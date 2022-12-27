package ui

import (
	"github.com/rivo/tview"
)

type View struct {
	App *tview.Application
}

func setupLayout() tview.Primitive {
	return tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Bookmarks"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(setupRequestPayload(), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Response"), 0, 3, false), 0, 2, false)
}

func setupRequestPayload() tview.Primitive {
	form := tview.NewForm().
		AddInputField("Server URL", "", 50, nil, nil).
		AddInputField("Last name", "", 20, nil, nil)
	form.SetBorder(true).SetTitle("Chiko v0.0.1").SetTitleAlign(tview.AlignCenter)
	return form
}

func NewView() View {
	app := tview.NewApplication()

	// list := tview.NewList().
	// 	ShowSecondaryText(false)
	// list.SetBorder(true).
	// 	SetTitleAlign(tview.AlignLeft)

	// main := tview.NewFlex().
	// 	AddItem(list, 0, 1, true)

	// pages := tview.NewPages().
	// 	AddPage("main", main, true, true)

	// frame := tview.NewFrame(pages)
	// frame.AddText("[::b][↓,j/↑,k][::-] Down/Up [::b][Enter,l/u,h][::-] Lower/Upper [::b][d[][::-] Download [::b][q[][::-] Quit", false, tview.AlignCenter, tcell.ColorWhite)
	layout := setupLayout()

	if err := app.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	v := View{
		app,
	}

	return v
}
