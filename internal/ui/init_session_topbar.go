package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// initSessionTopBar creates the Postman-style top bar:
// [ Server URL input ] [ Method label ] [ ▶ SEND button ]
func (u *UI) initSessionTopBar(sw *SessionWindow) *tview.Flex {
	focusHook := func() {
		u.GRPC = sw.GRPC
		u.ActiveSession = sw
	}

	// ── URL Field ─────────────────────────────────────────────
	urlField := tview.NewInputField()
	urlField.SetLabel(" 🔗 ")
	urlField.SetLabelColor(tcell.ColorAqua)
	urlField.SetText(sw.GRPC.Conn.ServerURL)
	urlField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	urlField.SetPlaceholder("  localhost:50051")
	urlField.SetPlaceholderStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	urlField.SetFocusFunc(focusHook)
	urlField.SetChangedFunc(func(text string) {
		sw.GRPC.Conn.ServerURL = text
		// Update window title dynamically
		title := text
		if title == "" {
			title = "New Request"
		}
		sw.WinBase.SetTitle(fmt.Sprintf(" %s ", title))
	})
	var lastConnectedURL string
	if sw.GRPC.Conn != nil {
		lastConnectedURL = sw.GRPC.Conn.ServerURL
	}

	urlField.SetBlurFunc(func() {
		currentURL := urlField.GetText()
		if currentURL != "" && currentURL != lastConnectedURL {
			lastConnectedURL = currentURL
			go u.connectSession(sw, urlField)
		}
		u.SaveWorkspace()
	})

	urlField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if u.History == nil || u.History.Entries == nil {
			return nil
		}

		// Collect unique URLs from history
		uniqueURLs := make(map[string]struct{})
		u.History.Mu.RLock()
		for _, entry := range *u.History.Entries {
			if entry.Session.ServerURL != "" {
				uniqueURLs[entry.Session.ServerURL] = struct{}{}
			}
		}
		u.History.Mu.RUnlock()

		if currentText == "" {
			for url := range uniqueURLs {
				entries = append(entries, url)
			}
			return entries
		}

		lowerCurrent := strings.ToLower(currentText)
		for url := range uniqueURLs {
			if strings.Contains(strings.ToLower(url), lowerCurrent) {
				entries = append(entries, url)
			}
		}
		return entries
	})

	urlField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			currentURL := urlField.GetText()
			lastConnectedURL = currentURL
			go u.connectSession(sw, urlField)
		}
		u.SaveWorkspace()
	})
	sw.URLField = urlField

	// ── Method label ──────────────────────────────────────────
	methodLabel := tview.NewButton("  [Select Method ▼]")
	methodLabel.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorDarkSlateBlue).
		Foreground(tcell.ColorWhite))
	// ── Method Field (Autocomplete) ───────────────────────────
	methodField := tview.NewInputField()
	methodField.SetLabel(" 🎯 ")
	methodField.SetLabelColor(tcell.ColorDarkGoldenrod)
	methodField.SetPlaceholder("  Select or type method...")
	methodField.SetPlaceholderStyle(tcell.StyleDefault.Foreground(tcell.ColorDarkGray))
	methodField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	methodField.SetFocusFunc(focusHook)

	methodField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if sw.GRPC == nil || sw.GRPC.Conn == nil || len(sw.GRPC.Conn.AvailableMethods) == 0 {
			return nil
		}
		if currentText == "" {
			return sw.GRPC.Conn.AvailableMethods
		}
		
		// Case-insensitive search 
		lowerCurrent := ""
		for _, char := range currentText {
			lowerCurrent += string(char)
			if lowerCurrent[len(lowerCurrent)-1] >= 'A' && lowerCurrent[len(lowerCurrent)-1] <= 'Z' {
				lowerCurrent = lowerCurrent[:len(lowerCurrent)-1] + string(lowerCurrent[len(lowerCurrent)-1]+32)
			}
		}
		
		for _, m := range sw.GRPC.Conn.AvailableMethods {
			lowerM := ""
			for _, char := range m {
				lowerM += string(char)
				if lowerM[len(lowerM)-1] >= 'A' && lowerM[len(lowerM)-1] <= 'Z' {
					lowerM = lowerM[:len(lowerM)-1] + string(lowerM[len(lowerM)-1]+32)
				}
			}

			// Simple substring match
			match := false
			for i := 0; i <= len(lowerM)-len(lowerCurrent); i++ {
				if lowerM[i:i+len(lowerCurrent)] == lowerCurrent {
					match = true
					break
				}
			}

			if match {
				entries = append(entries, m)
			}
		}
		return entries
	})

	methodField.SetChangedFunc(func(text string) {
		if text != "" {
			method := text
			sw.GRPC.Conn.SelectedMethod = &method
		} else {
			sw.GRPC.Conn.SelectedMethod = nil
		}
	})
	sw.MethodField = methodField

	// Update method field when method changes externally
	sw.RefreshTopBar = func() {
		u.App.QueueUpdateDraw(func() {
			if sw.GRPC.Conn.SelectedMethod != nil {
				methodField.SetText(*sw.GRPC.Conn.SelectedMethod)
			} else {
				methodField.SetText("")
			}
		})
	}

	// ── Send Button ────────────────────────────────────────────
	sendBtn := tview.NewButton("  ▶ SEND  ")
	sendBtn.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorDarkGreen).
		Foreground(tcell.ColorWhite))
	sendBtn.SetActivatedStyle(tcell.StyleDefault.
		Background(tcell.ColorGreen).
		Foreground(tcell.ColorWhite))
	sendBtn.SetFocusFunc(focusHook)
	sendBtn.SetSelectedFunc(func() {
		u.GRPC = sw.GRPC
		u.ActiveSession = sw
		u.InvokeRPC()
	})
	sw.SendBtn = sendBtn

	// ── Layout ─────────────────────────────────────────────────
	// Row 1: URL | SEND
	row1 := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(urlField, 0, 1, true).
		AddItem(sendBtn, 18, 0, false)

	// Row 2: Method
	bar := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(row1, 1, 0, true).
		AddItem(methodField, 1, 0, false)

	// Add 1 blank line of padding explicitly at the bottom to separate from Tabs
	bar.SetBorderPadding(0, 1, 1, 1)

	return bar
}
