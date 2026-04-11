package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

const (
	reqPanelBody     = "body"
	reqPanelMetadata = "metadata"
	reqPanelAuth     = "auth"
)

// initSessionRequestPanel creates a tabbed request section with Body/Metadata/Auth panels
func (u *UI) initSessionRequestPanel(sw *SessionWindow) (*tview.Pages, *tview.TextView) {
	pages := tview.NewPages()

	focusHook := func() {
		u.GRPC = sw.GRPC
		u.ActiveSession = sw
	}

	// ── Body Panel ────────────────────────────────────────────
	bodyArea := tview.NewTextArea()
	bodyArea.SetText(sw.GRPC.Conn.RequestPayload, false)
	bodyArea.SetBorder(false)
	bodyArea.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorLightGoldenrodYellow))
	bodyArea.SetFocusFunc(focusHook)
	bodyArea.SetMovedFunc(func() {
		// Sync text back to GRPC connection as user types
		sw.GRPC.Conn.RequestPayload = bodyArea.GetText()
		u.SaveWorkspace()
	})
	bodyArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlG {
			if sw.GRPC == nil || sw.GRPC.Conn == nil || sw.GRPC.Conn.SelectedMethod == nil {
				u.PrintOutput(entity.Output{
					SessionID:  sw.ID,
					Content:    "❗ please select rpc method first",
					WithHeader: false,
				})
				return nil
			}
			out, err := sw.GRPC.GenerateRPCPayloadSample()
			if err != nil {
				u.PrintOutput(entity.Output{
					SessionID:  sw.ID,
					Content:    fmt.Sprintf("❗ failed to generate sample: %s", err.Error()),
					WithHeader: false,
				})
				return nil
			}
			bodyArea.SetText(out, true)
			sw.GRPC.Conn.RequestPayload = out
			return nil
		}
		return event
	})
	sw.RequestBodyArea = bodyArea

	bodyWrapper := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(bodyArea, 0, 1, true)
	bodyWrapper.SetBorder(true).SetTitle(" Request Body (JSON) - [Ctrl+G] Generate Stub ")

	pages.AddPage(reqPanelBody, bodyWrapper, true, true)

	metaArea := tview.NewTextArea()
	metaText := ""
	if sw.GRPC.Conn != nil && len(sw.GRPC.Conn.Metadata) > 0 {
		for _, m := range sw.GRPC.Conn.Metadata {
			if m.Active {
				metaText += m.Key + ": " + m.Value + "\n"
			} else {
				metaText += "// " + m.Key + ": " + m.Value + "\n"
			}
		}
	} else {
		metaText = "// Example-Header: SomeValue\n"
	}

	metaArea.SetText(metaText, false)
	metaArea.SetBorder(false)
	metaArea.SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorLightCyan))
	metaArea.SetFocusFunc(focusHook)

	metaArea.SetMovedFunc(func() {
		// Sync metadata back instantly as the user types
		text := metaArea.GetText()
		lines := strings.Split(text, "\n")
		var parsed []*entity.Metadata

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			active := true
			if strings.HasPrefix(line, "//") {
				active = false
				line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				parsed = append(parsed, &entity.Metadata{
					Key:    strings.TrimSpace(parts[0]),
					Value:  strings.TrimSpace(parts[1]),
					Active: active,
				})
			}
		}
		sw.GRPC.Conn.Metadata = parsed
		u.SaveWorkspace()
	})

	metaWrapper := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(metaArea, 0, 1, true)
	metaWrapper.SetBorder(true).SetTitle(" Metadata (Key: Value) ")
	pages.AddPage(reqPanelMetadata, metaWrapper, true, false)

	authInput := tview.NewInputField()
	authInput.SetLabel(" Bearer Token: ")
	if sw.GRPC.Conn != nil && sw.GRPC.Conn.Authorization != nil && sw.GRPC.Conn.Authorization.BearerToken != nil {
		authInput.SetText(sw.GRPC.Conn.Authorization.BearerToken.Token)
	}
	authInput.SetPlaceholder("  Enter token here (without 'Bearer ' prefix)")
	authInput.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
	authInput.SetFocusFunc(focusHook)
	authInput.SetChangedFunc(func(text string) {
		if text == "" {
			sw.GRPC.Conn.Authorization = nil
		} else {
			sw.GRPC.Conn.Authorization = &entity.Auth{
				AuthType: entity.AuthTypeBearer,
				BearerToken: &entity.AuthValueBearerToken{
					Token: text,
				},
			}
		}
		u.SaveWorkspace()
	})

	authWrapper := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("\n  Configure Authorization\n").SetTextColor(tcell.ColorDarkGray), 3, 0, false).
		AddItem(authInput, 1, 0, true).
		AddItem(tview.NewBox(), 0, 1, false)

	authWrapper.SetBorder(true).SetTitle(" Authorization ")

	// Record references on SessionWindow so explicit shortcuts can target them
	sw.AuthInput = authInput
	sw.MetadataArea = metaArea

	pages.AddPage(reqPanelAuth, authWrapper, true, false)

	// ── Tab Bar ───────────────────────────────────────────────
	tabBar := tview.NewTextView()
	tabBar.SetDynamicColors(true)
	tabBar.SetRegions(true)
	tabBar.SetWrap(false)
	tabBar.SetBackgroundColor(tcell.ColorDarkSlateGray)
	tabBar.SetFocusFunc(focusHook)

	sw.RequestPages = pages
	sw.RequestTabBar = tabBar

	renderTabs := func(active string) {
		tabBar.Clear()
		tabs := []struct {
			id    string
			label string
		}{
			{reqPanelBody, " Body "},
			{reqPanelMetadata, " Metadata "},
			{reqPanelAuth, " Auth "},
		}
		for _, t := range tabs {
			if t.id == active {
				tabBar.Write([]byte(`["` + t.id + `"][white:darkblue::b]` + t.label + `[""] `))
			} else {
				tabBar.Write([]byte(`["` + t.id + `"][darkgray:black]` + t.label + `[""] `))
			}
		}
	}

	renderTabs(reqPanelBody)
	sw.RefreshRequestTabs = renderTabs

	tabBar.SetHighlightedFunc(func(added, removed, remaining []string) {
		if len(added) == 0 {
			return
		}
		tab := added[0]
		pages.SwitchToPage(tab)
		renderTabs(tab)
	})

	tabBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Key() == tcell.KeyRune && event.Rune() == 'b':
			tabBar.Highlight(reqPanelBody)
		case event.Key() == tcell.KeyRune && event.Rune() == 'm':
			tabBar.Highlight(reqPanelMetadata)
		case event.Key() == tcell.KeyRune && event.Rune() == 'a':
			tabBar.Highlight(reqPanelAuth)
		case event.Key() == tcell.KeyTAB:
			u.SetFocus(bodyArea)
		}
		return event
	})

	return pages, tabBar
}
