package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"

	"github.com/felangga/chiko/internal/entity"
)

// ShowDeleteCategoryConfirmation shows a confirmation dialog to delete a category directly,
// triggered by the Delete key on a focused category node.
func (u *UI) ShowDeleteCategoryConfirmation(category entity.Category) {
	u.ShowMessageBox(ShowMessageBoxParam{
		title:   " 🗑️ Delete Category ",
		message: fmt.Sprintf("Delete category \"%s\" and all its bookmarks?", category.Name),
		buttons: []Button{
			{
				Name: "Yes",
				OnClick: func(wnd *winman.WindowBase) {
					err := u.DeleteCategory(category)
					if err != nil {
						u.PrintLog(entity.Log{
							Content: fmt.Sprintf("❌ failed to delete category, err: %v", err),
							Type:    entity.LOG_ERROR,
						})
						return
					}

					u.PrintLog(entity.Log{
						Content: fmt.Sprintf("✅ category [blue]%s [white]deleted", category.Name),
						Type:    entity.LOG_INFO,
					})

					u.CloseModalDialog(wnd, u.Layout.BookmarkList)
				},
			},
			{
				Name: "No",
				OnClick: func(wnd *winman.WindowBase) {
					u.CloseModalDialog(wnd, u.Layout.BookmarkList)
				},
			},
		},
	})
}

// DeleteCategory removes a category (and all its sessions) from the bookmark list and saves it
func (u *UI) DeleteCategory(category entity.Category) error {
	categories := *u.Bookmark.Categories
	for i, cat := range categories {
		if cat.Name == category.Name {
			*u.Bookmark.Categories = append(categories[:i], categories[i+1:]...)

			err := u.Bookmark.SaveBookmark()
			if err != nil {
				return err
			}

			// Refresh the bookmark list
			u.loadBookmarks()
			return nil
		}
	}

	return nil
}
