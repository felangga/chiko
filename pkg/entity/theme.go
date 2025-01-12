package entity

import "github.com/gdamore/tcell/v2"

type ColorValue struct {
	Error            tcell.Color
	Warning          tcell.Color
	Notice           tcell.Color
	WindowColor      tcell.Color
	ModalColor       tcell.Color
	ButtonColor      tcell.Color
	FieldColor       tcell.Color
	PlaceholderColor tcell.Color
	CommandBarColor  tcell.Color
}

type Theme struct {
	Name   string
	Colors ColorValue
}

var (
	TerminalTheme = Theme{
		Name: "Terminal",
		Colors: ColorValue{
			Error:            tcell.GetColor("red"),
			Warning:          tcell.GetColor("darkred"),
			Notice:           tcell.GetColor("silver"),
			WindowColor:      tcell.GetColor("#444444"),
			ModalColor:       tcell.GetColor("#111111"),
			ButtonColor:      tcell.GetColor("#5500FF"),
			FieldColor:       tcell.GetColor("#111111"),
			PlaceholderColor: tcell.GetColor("#666666"),
			CommandBarColor:  tcell.GetColor("#333333"),
		},
	}
)
