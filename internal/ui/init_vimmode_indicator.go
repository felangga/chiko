package ui

import (
	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/events"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InitVimModeComponents struct {
	Layout *tview.Flex
}

func (u *UI) InitVimModeIndicator() *InitVimModeComponents {
	v := new(InitVimModeComponents)
	v.Layout = tview.NewFlex().SetDirection(tview.FlexRowCSS)
	vimIndc := tview.NewTextView().
		SetText("Normal").
		SetTextAlign(tview.AlignCenter)
	v.Layout.AddItem(vimIndc, 0, 1, false)
	v.Layout.SetBorder(true)

	// register events
	events.GetEventManager().Listen(events.EventTypeVimModeChange, func(i events.IEvent) {
		u.App.QueueUpdateDraw(func() {
			vmc := i.(*events.VimModeChanged)
			vimIndc.SetText(string(vmc.NewMode))
			switch vmc.NewMode {
			case entity.VimModeNormal:
				vimIndc.SetBackgroundColor(tcell.ColorGreen)
				vimIndc.SetTextColor(tcell.ColorBlack)
			case entity.VimModeEdit:
				vimIndc.SetBackgroundColor(tcell.ColorOrange)
				vimIndc.SetTextColor(tcell.ColorBlack)
			case entity.VimModeVisual:
				vimIndc.SetBackgroundColor(tcell.ColorPurple)
				vimIndc.SetTextColor(tcell.ColorBlack)
			}
		})
	})
	return v
}
