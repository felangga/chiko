package entity

type Theme struct {
	Name   string
	Colors map[string]string
}

var (
	TerminalTheme = Theme{
		Name: "Terminal",
		Colors: map[string]string{
			"Error":       "red",
			"Warning":     "darkred",
			"Notice":      "silver",
			"WindowColor": "#444444",
			"ModalColor":  "#111111",
			"ButtonColor": "#5500FF",
			"FieldColor":  "#111111",
		},
	}

	SelectedTheme = TerminalTheme
)
