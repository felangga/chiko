package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) ShowMetadataModal() {
	table := tview.NewTable()
	table.SetBorders(false).
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical).
		SetBackgroundColor(u.Theme.Colors.WindowColor)

	form := tview.NewForm()
	form.SetButtonsAlign(tview.AlignRight)
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)
	form.SetButtonStyle(u.Theme.Style.ButtonStyle)

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 0, 1, 1)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(table, 0, 8, true)
	flex.AddItem(form, 0, 1, false)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 📖 Metadata ",
		rootView:      flex,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 100, 30},
		fallbackFocus: u.Layout.MenuList,
	})

	form.AddButton("Delete Metadata", func() {
		if u.GRPC.Conn.Metadata == nil {
			return
		}

		row, col := table.GetSelection()
		cell := table.GetCell(row, col)
		u.deleteMetadataModal(wnd, table, cell)
	})

	form.AddButton("Add Metadata", func() {
		u.showAddMetadataModal(wnd, table)
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if form.GetButton(form.GetButtonIndex("Add Metadata")).HasFocus() {
				u.SetFocus(table)
			}

		}
		return event
	})

	u.ShowMetadataModal_SetInputCapture(wnd, table, form)
	u.ShowMetadataModal_RefreshMetadataTable(table)
}

func (u *UI) ShowMetadataModal_SetInputCapture(wnd *winman.WindowBase, table *tview.Table, form *tview.Form) {
	wnd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
			return nil
		}

		return event
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// If no metadata, do nothing
			if u.GRPC.Conn.Metadata == nil {
				return event
			}

			row, col := table.GetSelection()
			cell := table.GetCell(row, col)

			u.showEditMetadataModal(wnd, table, cell)

		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
			return nil
		case tcell.KeyTab:
			u.SetFocus(form)
			return nil
		}

		return event
	})

	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftDoubleClick {
			// If no metadata, do nothing
			if u.GRPC.Conn.Metadata == nil {
				return action, event
			}
			row, col := table.GetSelection()
			cell := table.GetCell(row, col)

			u.showEditMetadataModal(wnd, table, cell)
		}
		return action, event
	})
}

func (u *UI) ShowMetadataModal_RefreshMetadataTable(table *tview.Table) {
	table.Clear()

	// Set headers
	table.SetCell(0, 0, tview.NewTableCell("Active").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetSelectable(false))

	table.SetCell(0, 1, tview.NewTableCell("Key").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetExpansion(3).
		SetSelectable(false))

	table.SetCell(0, 2, tview.NewTableCell("Value").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetExpansion(4).
		SetSelectable(false))

	// Populate the table with data
	for row, meta := range u.GRPC.Conn.Metadata {
		// Active indicator
		activeCell := "✖"
		if meta.Active {
			activeCell = "✔"
		}

		table.SetCell(row+1, 0, tview.NewTableCell(activeCell).
			SetAlign(tview.AlignCenter).
			SetReference(meta).
			SetSelectable(true))

		// Key
		table.SetCell(row+1, 1, tview.NewTableCell(meta.Key).
			SetAlign(tview.AlignLeft).
			SetReference(meta).
			SetSelectable(true))

		// Value
		table.SetCell(row+1, 2, tview.NewTableCell(meta.Value).
			SetAlign(tview.AlignLeft).
			SetReference(meta).
			SetSelectable(true))
	}

}

func (u *UI) showAddMetadataModal(parentWnd *winman.WindowBase, table *tview.Table) {
	form := tview.NewForm()
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)
	form.SetButtonStyle(u.Theme.Style.ButtonStyle)
	form.SetFieldStyle(u.Theme.Style.FieldStyle)
	form.SetButtonsAlign(tview.AlignRight)

	inpKey := tview.NewInputField()
	inpKey.SetLabel("Key")
	form.AddFormItem(inpKey)

	inpValue := tview.NewInputField()
	inpValue.SetLabel("Value")
	form.AddFormItem(inpValue)

	chkActive := tview.NewCheckbox()
	chkActive.SetChecked(true)
	chkActive.SetLabel("Active")
	form.AddFormItem(chkActive)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " Add Metadata ",
		rootView:      form,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 50, 11},
		fallbackFocus: parentWnd,
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, parentWnd)
		}
		return event
	})

	form.AddButton("Cancel", func() {
		u.CloseModalDialog(wnd, parentWnd)
	})

	form.AddButton("Save", func() {
		// Create new metadata
		metadata := &entity.Metadata{
			Active: chkActive.IsChecked(),
			Key:    inpKey.GetText(),
			Value:  inpValue.GetText(),
		}

		u.GRPC.Conn.Metadata = append(u.GRPC.Conn.Metadata, metadata)
		u.ShowMetadataModal_RefreshMetadataTable(table)
		u.CloseModalDialog(wnd, parentWnd)
	})
}

func (u *UI) showEditMetadataModal(parentWnd *winman.WindowBase, table *tview.Table, cell *tview.TableCell) {
	metadata := cell.GetReference().(*entity.Metadata)

	form := tview.NewForm()
	form.SetBackgroundColor(u.Theme.Colors.WindowColor)
	form.SetButtonStyle(u.Theme.Style.ButtonStyle)
	form.SetFieldStyle(u.Theme.Style.FieldStyle)
	form.SetButtonsAlign(tview.AlignRight)

	inpKey := tview.NewInputField()
	inpKey.SetLabel("Key")
	inpKey.SetText(metadata.Key)
	form.AddFormItem(inpKey)

	inpValue := tview.NewInputField()
	inpValue.SetLabel("Value")
	inpValue.SetText(metadata.Value)
	form.AddFormItem(inpValue)

	chkActive := tview.NewCheckbox()
	chkActive.SetChecked(metadata.Active)
	chkActive.SetLabel("Active")
	form.AddFormItem(chkActive)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " Edit Metadata ",
		rootView:      form,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 50, 11},
		fallbackFocus: parentWnd,
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, parentWnd)
		}
		return event
	})

	form.AddButton("Cancel", func() {
		u.CloseModalDialog(wnd, parentWnd)
	})

	form.AddButton("Save", func() {
		// Create new metadata
		metadata.Active = chkActive.IsChecked()
		metadata.Key = inpKey.GetText()
		metadata.Value = inpValue.GetText()

		u.ShowMetadataModal_RefreshMetadataTable(table)

		u.CloseModalDialog(wnd, parentWnd)
	})
}

func (u *UI) deleteMetadataModal(parentWnd *winman.WindowBase, table *tview.Table, cell *tview.TableCell) {
	// Return if the metadata table is empty
	if len(u.GRPC.Conn.Metadata) < 1 {
		return
	}
	u.ShowMessageBox(ShowMessageBoxParam{
		title:   "Delete Metadata",
		message: "Are you sure you want to delete this metadata?",
		buttons: []Button{
			{
				Name: "Yes",
				OnClick: func(wnd *winman.WindowBase) {

					metadata := cell.GetReference().(*entity.Metadata)

					// Delete the metadata from the table
					// Find the index of the metadata in the slice
					var indexToRemove int
					for i, meta := range u.GRPC.Conn.Metadata {
						if meta == metadata {
							indexToRemove = i
							break
						}
					}
					// Remove the metadata from the slice
					u.GRPC.Conn.Metadata = append(u.GRPC.Conn.Metadata[:indexToRemove], u.GRPC.Conn.Metadata[indexToRemove+1:]...)

					u.ShowMetadataModal_RefreshMetadataTable(table)
					u.CloseModalDialog(wnd, parentWnd)
				},
			},
			{
				Name: "No",
				OnClick: func(wnd *winman.WindowBase) {
					u.CloseModalDialog(wnd, parentWnd)
				},
			},
		},
	})
}
