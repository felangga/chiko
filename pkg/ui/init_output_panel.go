package ui

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitOutputPanel initializes the output panel on the main screen
func (u *UI) InitOutputPanel() *tview.Table {

	output := tview.NewTable()
	//	output.SetBorders(true)
	output.SetTitle(" ⚙️ Output ").SetBorder(true)
	output.SetSelectable(true, false)
	output.SetSeparator(tview.Borders.Vertical)

	// Set headers
	output.SetCell(0, 0, tview.NewTableCell("Created At").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false))

	output.SetCell(0, 1, tview.NewTableCell("Log").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetExpansion(3).
		SetSelectable(false))

	output.SetCell(1, 0, tview.NewTableCell(time.Now().Format(time.RFC3339)).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(true))

	output.SetCell(1, 1, tview.NewTableCell(`the quick brown fox jumps over the lazy dog. 
	the quick brown fox jumps over the lazy dog. the quick brown fox jumps over the lazy dog.`).
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignLeft).
		SetExpansion(3).
		SetSelectable(true).
		SetMaxWidth(5))

	output.SetCell(2, 0, tview.NewTableCell(time.Now().Format(time.RFC3339)).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetExpansion(1).
		SetSelectable(true))

	output.SetCell(2, 1, tview.NewTableCell("the quick brown fox jumps over the lazy dog. the quick brown fox jumps over the lazy dog. the quick brown fox jumps over the lazy dog.").
		SetTextColor(tcell.ColorWhite).
		SetAlign(tview.AlignLeft).
		SetExpansion(3).
		SetSelectable(true))

	return output
}

func (u *UI) InitOutputPanel_SetInputCapture(outPanel *tview.TextArea) {
	outPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.App.SetFocus(u.Layout.MenuList)
		}
		return event
	})
}
