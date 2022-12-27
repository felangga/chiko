package ui

import (
	"github.com/gdamore/tcell/v2"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type menuItem struct {
	Text     string
	Shortcut tcell.Key
	Selected func()
}

type Menu struct {
	*tview.Box

	Items []menuItem
}

func NewMenu(app *tview.Application) *Menu {
	m := &Menu{
		Box: tview.NewBox(),
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		for _, i := range m.Items {
			if event.Key() == i.Shortcut {
				if i.Selected != nil {
					i.Selected()
				}
				break
			}
		}

		return event
	})

	return m
}

func (m *Menu) Draw(screen tcell.Screen) {
	x, y, width, _ := m.GetInnerRect()
	x++

	for _, i := range m.Items {
		tview.Print(screen, tcell.KeyNames[i.Shortcut], x, y, width-2, tview.AlignLeft, tcell.GetColor(theme.Colors["Shortcut"]))
		x += 3

		tview.Print(screen, i.Text, x, y, width, tview.AlignLeft, tcell.GetColor(theme.Colors["Menu"]))
		x += runewidth.StringWidth(i.Text) + 2
	}
}

func (m *Menu) AddItem(text string, shortcut tcell.Key, selected func()) *Menu {
	m.Items = append(m.Items, menuItem{
		Text:     text,
		Shortcut: shortcut,
		Selected: selected,
	})

	return m
}

func (m *Menu) Clear() *Menu {
	m.Items = nil

	return m
}
