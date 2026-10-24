package ui

// App shortcut runes. When --vim is on, keys that conflict with vim
// motions/modes are remapped so features stay reachable.
//
// Vim-reserved on text widgets: h j k l w b e v i a
func (u *UI) keyHistory() rune {
	if u.VimEnabled {
		return 'y' // historY
	}
	return 'h'
}

func (u *UI) keyAuthorization() rune {
	if u.VimEnabled {
		return 't' // Token
	}
	return 'a'
}

func (u *UI) keyInvoke() rune {
	if u.VimEnabled {
		return 'r' // Run
	}
	return 'i'
}

func (u *UI) keyBookmark() rune {
	if u.VimEnabled {
		return 's' // Save bookmark
	}
	return 'b'
}

func (u *UI) keySelectAll() rune {
	if u.VimEnabled {
		return 'A' // Shift-a; plain 'a' enters Edit
	}
	return 'a'
}

func (u *UI) keyDumpFile() rune {
	if u.VimEnabled {
		return 'f' // dump to File; plain 'w' is word-forward
	}
	return 'w'
}
