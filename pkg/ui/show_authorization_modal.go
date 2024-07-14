package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

// ShowAuthorizationModal used to show the authorization modal dialog
func (u *UI) ShowAuthorizationModal() {

	var currentToken string
	if u.GRPC.Conn.Authorization != nil {
		currentToken = u.GRPC.Conn.Authorization.BearerToken.Token
	}

	txtAuthorization := tview.NewInputField()
	txtAuthorization.SetText(currentToken)
	txtAuthorization.SetFieldBackgroundColor(u.Theme.Colors.WindowColor)
	txtAuthorization.SetBackgroundColor(u.Theme.Colors.WindowColor)

	// Set placeholder text and style
	style := tcell.StyleDefault.
		Background(u.Theme.Colors.WindowColor).
		Italic(true)

	txtAuthorization.SetPlaceholder("Set token without the Bearer prefix")
	txtAuthorization.SetPlaceholderTextColor(u.Theme.Colors.PlaceholderColor)
	txtAuthorization.SetPlaceholderStyle(style)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
		title:         " 🔑 Authorization ",
		rootView:      txtAuthorization,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 100, 1},
		fallbackFocus: u.Layout.MenuList,
	})

	u.ShowAuthorizationModal_SetInputCapture(wnd, txtAuthorization, currentToken)
}

func (u *UI) ShowAuthorizationModal_SetInputCapture(wnd *winman.WindowBase, txtAuthorization *tview.InputField, currentToken string) {
	txtAuthorization.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		case tcell.KeyEnter:
			// Check if empty
			if len(txtAuthorization.GetText()) < 1 {
				u.GRPC.Conn.Authorization = nil
				if len(currentToken) > 0 {
					u.PrintLog(entity.Log{
						Content: "🔒 Authorization removed",
						Type:    entity.LOG_INFO,
					})
				}
			} else {
				auth := entity.Auth{
					AuthType: entity.AuthTypeBearer,
					BearerToken: &entity.AuthValueBearerToken{
						Token: txtAuthorization.GetText(),
					},
				}

				u.GRPC.Conn.Authorization = &auth
				u.PrintLog(entity.Log{
					Content: fmt.Sprintf("🔒 Authorization set [blue]%s", txtAuthorization.GetText()),
					Type:    entity.LOG_INFO,
				})
			}

			// Remove the window and restore focus to menu list
			u.CloseModalDialog(wnd, u.Layout.MenuList)
		}
		return event
	})
}
