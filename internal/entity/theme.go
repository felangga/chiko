package entity

import "github.com/gdamore/tcell/v2"

type ColorValue struct {
	Error           tcell.Color
	Warning         tcell.Color
	Notice          tcell.Color
	WindowColor     tcell.Color
	ModalColor      tcell.Color
	CommandBarColor tcell.Color
}

type ComponentStyle struct {
	ButtonStyle       tcell.Style
	PlaceholderStyle  tcell.Style
	FieldStyle        tcell.Style
	ListMainTextStyle tcell.Style
	ListBorderStyle   tcell.Style
	TextAreaStyle     tcell.Style
}

type Theme struct {
	Name   string
	Colors ColorValue
	Style  ComponentStyle
}

var (
	TerminalTheme = Theme{
		Name: "Terminal",
		Colors: ColorValue{
			Error:           tcell.GetColor("red"),
			Warning:         tcell.GetColor("darkred"),
			Notice:          tcell.GetColor("silver"),
			WindowColor:     tcell.GetColor("#444444"),
			ModalColor:      tcell.GetColor("#111111"),
			CommandBarColor: tcell.GetColor("#333333"),
		},
		Style: ComponentStyle{
			ButtonStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#5500FF")),
			PlaceholderStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#666666")).Italic(true),
			FieldStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#666666")),
			ListMainTextStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#444444")),
			ListBorderStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#444444")),
			TextAreaStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#444444")),
		},
	}
)
