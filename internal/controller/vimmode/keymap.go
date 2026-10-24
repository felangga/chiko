package vimmode

import "github.com/gdamore/tcell/v2"

var (
	NormalKeyMap KeyMap = make(KeyMap)
	EditKeyMap   KeyMap = make(KeyMap)
	VisualKeyMap KeyMap = make(KeyMap)
)

// keyID identifies an EventKey by value so Bind/GetBinding work across
// separately allocated *tcell.EventKey pointers.
type keyID struct {
	key tcell.Key
	ch  rune
	mod tcell.ModMask
}

func eventKeyID(ev *tcell.EventKey) keyID {
	k := ev.Key()
	if k == tcell.KeyRune {
		return keyID{key: k, ch: ev.Rune(), mod: ev.Modifiers()}
	}
	// Non-rune keys: ignore rune payload (terminals disagree on ch).
	return keyID{key: k, ch: 0, mod: ev.Modifiers()}
}

type KeyMap map[keyID]*KeyBind

func (km KeyMap) Bind(from *tcell.EventKey, to *tcell.EventKey) {
	km[eventKeyID(from)] = &KeyBind{
		From: from,
		To:   to,
	}
}

func (km KeyMap) GetBinding(from *tcell.EventKey) *KeyBind {
	if v, ok := km[eventKeyID(from)]; ok {
		return v
	}
	// KeyCtrlH/J/K/L often arrive with ModCtrl and/or ModShift set even
	// though bindings are registered with ModNone.
	switch from.Key() {
	case tcell.KeyCtrlH, tcell.KeyCtrlJ, tcell.KeyCtrlK, tcell.KeyCtrlL:
		if v, ok := km[keyID{key: from.Key(), mod: tcell.ModNone}]; ok {
			return v
		}
	}
	return &KeyBind{
		From: from,
		To:   from,
	}
}

type KeyBind struct {
	From *tcell.EventKey
	To   *tcell.EventKey
}
