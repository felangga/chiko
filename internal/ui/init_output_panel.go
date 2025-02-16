package ui

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

type InitOutputPanelComponents struct {
	Layout   *tview.Flex
	TextArea *tview.TextArea
	Buffer   string
}

var (
	commands map[string]Commands
)

// InitOutputPanel initializes the output panel on the main screen
func (u *UI) InitOutputPanel() InitOutputPanelComponents {
	output := tview.NewTextArea()
	output.SetWrap(false)
	output.SetMaxLength(1)
	output.SetTextStyle(tcell.StyleDefault.
		Foreground(tcell.ColorGreen))

	u.initOutputPanel_handleTextArea(output)

	layout := tview.NewFlex()
	layout.SetDirection(tview.FlexRow)
	layout.SetBorder(true)
	layout.SetTitle(" Output ")
	layout.AddItem(output, 0, 1, true)
	layout.AddItem(u.initOutputPanel_PanelBar(), 1, 1, false)

	return InitOutputPanelComponents{
		Layout:   layout,
		TextArea: output,
		Buffer:   "",
	}
}

func (u *UI) initOutputPanel_handleTextArea(textarea *tview.TextArea) {
	commands = map[string]Commands{
		"selectall": {
			KeyComb:    'a',
			CommandKey: "A",
			Text:       "Select All",
			OnExecute: func() {
				u.doSelectAll(textarea)
			},
		},
		"copy": {
			KeyComb:    'c',
			CommandKey: "C",
			Text:       "Copy",
			OnExecute: func() {
				u.doCopyText(textarea)
			},
		},
		"writefile": {
			KeyComb:    'w',
			CommandKey: "W",
			Text:       "Dump To File",
			OnExecute: func() {
				u.doWriteToFile(textarea)
			},
		},
	}

	textarea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Block backspace for editing the text area
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			return nil
		}

		// Change focus to menu list when pressed TAB
		if event.Key() == tcell.KeyTAB {
			u.SetFocus(u.Layout.LogList)
			return nil
		}

		for _, k := range commands {
			if event.Key() == tcell.KeyRune && tcell.Key(event.Rune()) == tcell.Key(k.KeyComb) {
				k.OnExecute()
				return nil
			}
		}
		return event
	})
}

// doCopyText copy text selected from the output window to clipboard
func (u *UI) doCopyText(textarea *tview.TextArea) {
	text, _, _ := textarea.GetSelection()
	if text == "" {
		u.PrintLog(entity.Log{
			Content: "❌ no text selected",
			Type:    entity.LOG_INFO,
		})
		return
	}

	if err := clipboard.WriteAll(text); err != nil {
		u.PrintLog(entity.Log{
			Content: "❌ failed to copied to clipboard",
			Type:    entity.LOG_INFO,
		})
		return
	}

	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("📋 %.2f kB copied to clipboard", float64(len(text))/1024),
		Type:    entity.LOG_INFO,
	})
}

func (u *UI) doSelectAll(textarea *tview.TextArea) {
	textarea.Select(0, textarea.GetTextLength())
}

func (u *UI) doWriteToFile(textarea *tview.TextArea) {
	// Show dialog box so user can input the file path
	dir, err := u.Storage.GetWorkingDirectory()
	if err != nil {
		u.PrintLog(entity.Log{
			Content: "❌ failed to get working directory",
			Type:    entity.LOG_ERROR,
		})
		return
	}

	txtPath := tview.NewInputField().SetText(dir + "/" + u.GRPC.Conn.ServerURL + ".txt")
	txtPath.SetFieldBackgroundColor(u.Theme.Colors.WindowColor)
	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:      " 💾 Export Path ",
		draggable:  true,
		resizeable: false,
		size: winSize{
			x:      0,
			y:      0,
			width:  80,
			height: 1,
		},
		rootView:      txtPath,
		fallbackFocus: u.Layout.MenuList,
	})

	txtPath.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		case tcell.KeyEnter:
			if err := u.Storage.SaveToFile(txtPath.GetText(), []byte(textarea.GetText())); err != nil {
				u.PrintLog(entity.Log{
					Content: "❌ failed to write to file: " + err.Error(),
					Type:    entity.LOG_ERROR,
				})
				return nil
			}
			// Remove the window and restore focus to menu list
			defer u.CloseModalDialog(wnd, u.Layout.MenuList)
			u.PrintLog(entity.Log{
				Content: "✅ file saved successfully: " + txtPath.GetText(),
				Type:    entity.LOG_INFO,
			})
		}
		return event
	})

}

type Commands struct {
	KeyComb    rune
	CommandKey string
	Text       string
	OnExecute  func()
}

func (u *UI) initOutputPanel_PanelBar() *tview.TextView {
	// Init Panel bar
	info := tview.NewTextView()
	info.SetDynamicColors(true)
	info.SetRegions(true)
	info.SetWrap(false)
	info.SetBackgroundColor(u.Theme.Colors.CommandBarColor)

	info.SetHighlightedFunc(func(added, removed, remaining []string) {
		if len(added) < 1 {
			return
		}

		getCommand := added[0]

		// Execute defined commands
		if c, ok := commands[getCommand]; ok {
			c.OnExecute()
		}

		// Set highlight to "blank" after 1 second
		go func() {
			time.Sleep(time.Second / 2)
			u.App.QueueUpdateDraw(func() {
				info.Highlight("blank")
			})
		}()
	})

	info.SetText(`["blank"][""]`)

	// Populate the panel bar
	for c, cmd := range commands {
		fmt.Fprintf(info, `[lightcyan]%s ["%s"][lightgrey]%s[""]  `, cmd.CommandKey, c, cmd.Text)
	}

	return info
}
