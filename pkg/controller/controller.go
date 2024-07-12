package controller

// import (
// 	"chiko/pkg/entity"
// 	"chiko/pkg/ui"
// 	"context"
// 	"fmt"

// 	"github.com/rivo/tview"
// )

// type Controller struct {
// 	ctx       context.Context
// 	ui        ui.UI
// 	conn      *entity.Session
// 	bookmarks *[]entity.Bookmark
// 	theme     entity.Theme
// }

// func NewController() Controller {
// 	ui := ui.NewUI()
// 	conn := entity.Session{
// 		// Default server URL
// 		ServerURL: "localhost:20010",
// 	}

// 	bookmarks := []entity.Bookmark{}

// 	c := Controller{
// 		context.Background(),
// 		ui,
// 		&conn,
// 		&bookmarks,
// 		entity.TerminalTheme,
// 	}

// 	// Initialize bookmark tree view
// 	ui.BookmarkList.SetSelectedFunc(func(node *tview.TreeNode) {
// 		reference := node.GetReference()
// 		if reference == nil {
// 			return // Selecting the root node does nothing.
// 		}
// 		children := node.GetChildren()
// 		if len(children) == 0 {
// 			switch node.GetReference().(type) {
// 			case entity.Session:
// 				getSession := node.GetReference().(entity.Session)
// 				c.ShowBookmarkOptionsModal(c.ui.MenuList, getSession)
// 			}

// 		} else {
// 			// Collapse if visible, expand if collapsed.
// 			node.SetExpanded(!node.IsExpanded())
// 		}
// 	})

// 	return c
// }

// func (c Controller) Run() error {
// 	c.InitMenu()

// 	c.initSys()

// 	c.ui.App.EnableMouse(true)
// 	return c.ui.App.Run()
// }

// func (c Controller) initSys() {
// 	c.PrintLog(fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION), LOG_INFO)

// 	// Load bookmarks
// 	c.loadBookmarks()
// }
