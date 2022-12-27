package ui

type Theme struct {
	Name   string
	Colors map[string]string
}

var (
	IceTheme = Theme{
		Name: "Ice",
		Colors: map[string]string{
			"Error":       "#aacdee",
			"Warning":     "#7997d6",
			"Notice":      "#5b73b4",
			"Message":     "gray",
			"Timestamp":   "#003d6b",
			"Service":     "#d5d8f0",
			"Shortcut":    "#00ffdf",
			"Menu":        "#85e7d2",
			"Description": "#003d6b",
		},
	}

	TerminalTheme = Theme{
		Name: "Terminal",
		Colors: map[string]string{
			"Error":       "red",
			"Warning":     "darkred",
			"Notice":      "silver",
			"Message":     "gray",
			"Timestamp":   "green",
			"Service":     "blue",
			"Shortcut":    "yellow",
			"Menu":        "teal",
			"Description": "green",
		},
	}

	theme = TerminalTheme
)
