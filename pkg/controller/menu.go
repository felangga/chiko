package controller

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/epiclabs-io/winman"
// 	"github.com/fullstorydev/grpcurl"
// 	"github.com/gdamore/tcell/v2"
// 	"github.com/jhump/protoreflect/desc"
// 	"github.com/lithammer/fuzzysearch/fuzzy"
// 	"github.com/rivo/tview"

// 	"chiko/pkg/entity"
// )

// func (c Controller) AddMenuSection(name string) {
// 	c.ui.MenuList.AddItem("[::d]"+name, "", 0, nil)
// 	c.ui.MenuList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
// }

// func (c Controller) InitMenu() {
// 	menuList := c.ui.MenuList
// 	menuList.AddItem("Server URL", "", 'u', c.SetServerURL)
// 	menuList.AddItem("Methods", "", 'm', c.SetRequestMethods)
// 	menuList.AddItem("Authorization", "", 'a', c.SetAuthorizationModal)
// 	menuList.AddItem("Metadata", "", 'd', nil)
// 	menuList.AddItem("Request Payload", "", 'p', c.SetRequestPayload)
// 	menuList.AddItem("Invoke", "", 'i', c.DoInvoke)
// 	menuList.AddItem("[::d]"+strings.Repeat(string(tcell.RuneHLine), 25), "", 0, nil)
// 	menuList.AddItem("Save to Bookmark", "", 'b', c.DoSaveBookmark)
// 	menuList.AddItem("Quit", "", 'q', c.DoQuit)
// }

// // DoQuit used to do shut down sequence and quit the program
// func (c Controller) DoQuit() {
// 	// Close active connection if we are going to close the application
// 	if (c.conn != nil) && (c.conn.ActiveConnection != nil) {
// 		err := c.conn.ActiveConnection.Close()
// 		if err != nil {
// 			c.PrintLog(err.Error(), LOG_ERROR)
// 		}
// 	}

// 	c.ui.App.Stop()
// }

// // SetServerURL used to show set server URL modal dialog
// func (c Controller) SetServerURL() {
// 	// Create Set Server URL From
// 	txtServerURL := tview.NewInputField().SetText(c.conn.ServerURL)
// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(txtServerURL).
// 		SetDraggable(true).
// 		SetTitle("🌏 Enter Server URL ")

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	txtServerURL.SetFieldBackgroundColor(wnd.GetBackgroundColor())

// 	txtServerURL.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		case tcell.KeyEnter:

// 			go c.CheckGRPC(txtServerURL.GetText())

// 			// Remove the window and restore focus to menu list
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		}
// 		return event
// 	})

// 	wnd.SetModal(true)
// 	wnd.SetRect(0, 0, 70, 1)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }

// // SetRequestMethods used to set the RPC request methods, it will show the available methods
// // if the server supports Server Reflection
// func (c Controller) SetRequestMethods() {
// 	if c.conn.ActiveConnection == nil {
// 		c.PrintLog("❗ no active connection", LOG_WARNING)
// 		return
// 	}

// 	// Create Set Server URL From
// 	// Set placeholder text and style
// 	style := tcell.StyleDefault.
// 		Background(c.theme.Colors.PlaceholderColor).
// 		Italic(true)

// 	listMethods := tview.NewList().
// 		ShowSecondaryText(false)
// 	listMethods.SetBorderPadding(1, 0, 0, 0)
// 	listMethods.SetBackgroundColor(c.theme.Colors.WindowColor)

// 	inpMethods := tview.NewInputField().
// 		SetText("").
// 		SetPlaceholder(" 🔍 Search methods...").
// 		SetPlaceholderStyle(style)

// 	inpMethods.SetFieldBackgroundColor(c.theme.Colors.PlaceholderColor)

// 	wndLayer := tview.NewFlex()
// 	wndLayer.SetDirection(tview.FlexRow)
// 	wndLayer.AddItem(inpMethods, 1, 1, true)
// 	wndLayer.AddItem(listMethods, 0, 1, false)

// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(wndLayer).
// 		SetDraggable(true).
// 		SetResizable(true).
// 		SetTitle(" Select RPC Methods ")

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	wnd.SetBorderPadding(1, 1, 1, 1)

// 	inpMethods.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		getText := inpMethods.GetText()

// 		switch event.Key() {
// 		case tcell.KeyDown, tcell.KeyEnter, tcell.KeyTAB:
// 			c.ui.SetFocus(listMethods)
// 			return nil
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MainLayout)
// 			return nil
// 		case tcell.KeyRune:
// 			getText = getText + string(event.Rune())
// 		case tcell.KeyBackspace, tcell.KeyBackspace2:
// 			if len(getText) > 0 {
// 				getText = getText[:len(getText)-1]
// 			}
// 		}

// 		fuzzyFind := fuzzy.FindFold(getText, c.conn.AvailableMethods)
// 		c.refreshList(listMethods, fuzzyFind)

// 		return event
// 	})

// 	listMethods.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		case tcell.KeyEnter:
// 			idx := listMethods.GetCurrentItem()
// 			item, _ := listMethods.GetItemText(idx)
// 			c.conn.SelectedMethod = &item

// 			// Reset the search bar
// 			inpMethods.SetText("")

// 			// Remove the window and restore focus to menu list
// 			c.PrintLog("👉 Method set to [blue]"+*c.conn.SelectedMethod, LOG_INFO)
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		case tcell.KeyRune, tcell.KeyBackspace, tcell.KeyBackspace2:
// 			// Set focus to input field when user typing something
// 			c.ui.SetFocus(inpMethods)
// 		case tcell.KeyTAB:
// 			c.ui.SetFocus(inpMethods)
// 			return nil
// 		}
// 		return event
// 	})

// 	c.refreshList(listMethods, c.conn.AvailableMethods)

// 	wnd.SetModal(true)
// 	wnd.SetRect(0, 0, 100, 20)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	var maxMinButton *winman.Button
// 	maxMinButton = &winman.Button{
// 		Symbol:    '▴',
// 		Alignment: winman.ButtonRight,
// 		OnClick: func() {
// 			if wnd.IsMaximized() {
// 				wnd.Restore()
// 				maxMinButton.Symbol = '▴'
// 			} else {
// 				wnd.Maximize()
// 				maxMinButton.Symbol = '▾'
// 			}
// 		},
// 	}
// 	wnd.AddButton(maxMinButton)

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }

// func (c Controller) refreshList(listView *tview.List, items []string) {
// 	listView.Clear()
// 	listView.SetCurrentItem(0)

// 	for i, method := range items {
// 		listView.AddItem(method, "", 0, nil)

// 		// Ignore if none was selected before
// 		if c.conn.SelectedMethod != nil && *c.conn.SelectedMethod == method {
// 			listView.SetCurrentItem(i)
// 		}
// 	}

// }

// // SetAuthorizationModal used to show the authorization modal dialog
// func (c Controller) SetAuthorizationModal() {
// 	// Create set authorization modal
// 	var bearerToken string
// 	if c.conn.Authorization != nil {
// 		bearerToken = c.conn.Authorization.BearerToken.Token
// 	}

// 	txtAuthorization := tview.NewInputField().SetText(bearerToken)
// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(txtAuthorization).
// 		SetDraggable(true).
// 		SetTitle("🔑 Authorization ")

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	txtAuthorization.SetFieldBackgroundColor(wnd.GetBackgroundColor())
// 	txtAuthorization.SetBackgroundColor(wnd.GetBackgroundColor())

// 	// Set placeholder text and style
// 	style := tcell.StyleDefault.
// 		Background(wnd.GetBackgroundColor()).
// 		Italic(true)

// 	txtAuthorization.SetPlaceholder("Set token without the Bearer prefix")
// 	txtAuthorization.SetPlaceholderTextColor(c.theme.Colors.PlaceholderColor)
// 	txtAuthorization.SetPlaceholderStyle(style)

// 	txtAuthorization.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		case tcell.KeyEnter:
// 			// Check if empty
// 			if len(txtAuthorization.GetText()) < 1 {
// 				c.conn.Authorization = nil
// 				if len(bearerToken) > 0 {
// 					c.PrintLog("🔓 Authorization removed", LOG_INFO)
// 				}
// 			} else {
// 				auth := entity.Auth{
// 					AuthType: entity.AuthTypeBearer,
// 					BearerToken: &entity.AuthValueBearerToken{
// 						Token: txtAuthorization.GetText(),
// 					},
// 				}

// 				c.conn.Authorization = &auth
// 				c.PrintLog(fmt.Sprintf("🔒 Authorization set [blue]%s", txtAuthorization.GetText()), LOG_INFO)
// 			}

// 			// Remove the window and restore focus to menu list
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		}
// 		return event
// 	})

// 	wnd.SetModal(true)
// 	wnd.SetRect(0, 0, 100, 1)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }

// // SetRequestPayload used to set the request payload and also user can generate the sample payload if the
// // Server Reflection is supported
// func (c Controller) SetRequestPayload() {
// 	if c.conn.ActiveConnection == nil {
// 		c.PrintLog("❗ no active connection", LOG_WARNING)
// 		return
// 	}

// 	if c.conn.SelectedMethod == nil {
// 		c.PrintLog("❗ please select rpc method first", LOG_ERROR)
// 		return
// 	}

// 	requestPayload := c.conn.RequestPayload

// 	// Create Set Server URL From
// 	form := tview.NewForm()
// 	wnd := winman.NewWindow().
// 		Show().
// 		SetRoot(form).
// 		SetTitle(" Request Payload ").
// 		SetDraggable(true)

// 	wnd.SetBackgroundColor(c.theme.Colors.WindowColor)

// 	form.SetBackgroundColor(c.theme.Colors.WindowColor)
// 	form.SetBorderPadding(1, 1, 0, 1)

// 	// Create text area for filling the payload
// 	txtPayload := tview.NewTextArea().SetText(requestPayload, true)
// 	txtPayload.SetSize(9, 100)
// 	form.SetFieldBackgroundColor(c.theme.Colors.FieldColor)
// 	form.AddFormItem(txtPayload)
// 	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
// 		switch event.Key() {
// 		case tcell.KeyEscape:
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 			return nil
// 		}
// 		return event
// 	})

// 	// Populate Buttons
// 	form.AddButton("Generate Sample", func() {

// 		// Get service detail
// 		dsc, err := c.conn.DescriptorSource.FindSymbol(*c.conn.SelectedMethod)
// 		if err != nil {
// 			c.PrintLog(err.Error(), LOG_ERROR)
// 			return
// 		}

// 		txt, err := grpcurl.GetDescriptorText(dsc, c.conn.DescriptorSource)
// 		if err != nil {
// 			c.PrintLog(err.Error(), LOG_ERROR)
// 			return
// 		}

// 		// Parse the service to get request message name
// 		rr := c.parseRequestResponse(txt)
// 		if len(rr) < 2 {
// 			c.PrintLog(fmt.Sprintf("failed to parse service name: %s", txt), LOG_ERROR)
// 		}
// 		// Remove stream from request
// 		requestMessage := strings.ReplaceAll(rr[0][1], "stream", "")

// 		// Trim message
// 		requestMessage = strings.TrimSpace(requestMessage)
// 		if requestMessage[0:1] == "." {
// 			requestMessage = requestMessage[1:]
// 		}

// 		// Retrieve request message from descriptors
// 		dsc, err = c.conn.DescriptorSource.FindSymbol(requestMessage)
// 		if err != nil {
// 			c.PrintLog(err.Error(), LOG_ERROR)
// 			return
// 		}

// 		if dsc, ok := dsc.(*desc.MessageDescriptor); ok {
// 			// for messages, also show a template in JSON, to make it easier to
// 			// create a request to invoke an RPC
// 			tmpl := grpcurl.MakeTemplate(dsc)
// 			options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
// 			_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), c.conn.DescriptorSource, nil, options)
// 			if err != nil {
// 				c.PrintLog(err.Error(), LOG_ERROR)
// 				return
// 			}
// 			str, err := formatter(tmpl)
// 			if err != nil {
// 				c.PrintLog(err.Error(), LOG_ERROR)
// 				return
// 			}
// 			txtPayload.SetText(str, true)
// 		}
// 	})

// 	form.AddButton("Apply", func() {
// 		c.conn.RequestPayload = txtPayload.GetText()

// 		// Remove the window and restore focus to menu list
// 		c.PrintLog("\nRequest Payload:\n[yellow]"+c.conn.RequestPayload, LOG_INFO)
// 		c.ui.WinMan.RemoveWindow(wnd)
// 		c.ui.SetFocus(c.ui.MenuList)
// 	})

// 	form.SetButtonsAlign(tview.AlignRight)

// 	wnd.SetModal(true)
// 	wnd.SetRect(0, 0, 70, 15)
// 	wnd.AddButton(&winman.Button{
// 		Symbol: 'X',
// 		OnClick: func() {
// 			c.ui.WinMan.RemoveWindow(wnd)
// 			c.ui.SetFocus(c.ui.MenuList)
// 		},
// 	})

// 	c.ui.WinMan.AddWindow(wnd)
// 	c.ui.WinMan.Center(wnd)
// 	c.ui.SetFocus(wnd)
// }
