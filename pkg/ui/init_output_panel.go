package ui

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	"github.com/felangga/chiko/pkg/entity"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	output.SetWrap(true)
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
		"copy": {
			KeyComb:    'c',
			CommandKey: "C",
			Text:       "Copy",
			OnExecute: func() {
				textArea := u.Layout.OutputPanel.TextArea

				text, _, _ := textArea.GetSelection()
				if err := clipboard.WriteAll(text); err != nil {
					u.PrintLog(entity.Log{
						Content: "❌ failed to copied to clipboard",
						Type:    entity.LOG_INFO,
					})
					return
				}

				u.PrintLog(entity.Log{
					Content: "📋 Copied to clipboard",
					Type:    entity.LOG_INFO,
				})
			},
		},
		"writefile": {
			KeyComb:    'w',
			CommandKey: "W",
			Text:       "Write To File",
			OnExecute:  func() {},
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
