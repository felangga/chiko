package entity

type VimMode string

const (
	VimModeNormal VimMode = "Normal"
	VimModeEdit   VimMode = "Edit"
	VimModeVisual VimMode = "Visual"
)
