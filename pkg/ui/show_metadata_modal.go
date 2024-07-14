package ui

func (u *UI) ShowMetadataModal() {
	u.ShowMessageBox(ShowMessageBoxParam{
		title:   " ✨ Coming Soon ",
		message: "This feature is not yet implemented. Please stay tuned for the next update.",
	})
}
