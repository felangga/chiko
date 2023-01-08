package controller

import (
	"chiko/pkg/entity"
	"fmt"
	"strings"

	"github.com/epiclabs-io/winman"
	"github.com/fullstorydev/grpcurl"
	"github.com/gdamore/tcell/v2"
	"github.com/jhump/protoreflect/desc"
	"github.com/rivo/tview"
)

func (c Controller) initMenu() {
	c.ui.MenuList.AddItem("Server URL", "", 'u', c.setServerURL)
	c.ui.MenuList.AddItem("Method", "", 'm', c.setMethod)
	c.ui.MenuList.AddItem("Authorization", "", 'a', nil)
	c.ui.MenuList.AddItem("Metadata", "", 'd', nil)
	c.ui.MenuList.AddItem("Request Payload", "", 'p', c.setRequestPayload)
	c.ui.MenuList.AddItem("Invoke", "", 'i', c.doInvoke)
	c.ui.MenuList.AddItem("------------------", "", 0, nil)
	c.ui.MenuList.AddItem("Save to Bookmark", "", 'b', c.doSaveBookmark)
}

func (c Controller) setServerURL() {
	// Create Set Server URL From
	txtServerURL := tview.NewInputField().SetText(c.conn.ServerURL)
	wnd := winman.NewWindow().
		Show().
		SetRoot(txtServerURL).
		SetDraggable(true).
		SetTitle("🌏 Enter Server URL ")

	wnd.SetBackgroundColor(tcell.GetColor(entity.SelectedTheme.Colors["WindowColor"]))
	txtServerURL.SetFieldBackgroundColor(wnd.GetBackgroundColor())

	txtServerURL.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		case tcell.KeyEnter:

			go c.CheckGRPC(txtServerURL.GetText())

			// Remove the window and restore focus to menu list
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		}
		return event
	})

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 50, 1)
	wnd.AddButton(&winman.Button{
		Symbol: '❌',
		OnClick: func() {
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		},
	})

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

func (c Controller) setMethod() {
	if c.conn.ActiveConnection == nil {
		c.PrintLog("❗ no active connection", LOG_WARNING)
		return
	}

	// Create Set Server URL From
	listMethods := tview.NewList().
		ShowSecondaryText(false)

	wnd := winman.NewWindow().
		Show().
		SetRoot(listMethods).
		SetDraggable(true).
		SetResizable(true).
		SetTitle(" Select gRPC Method ")

	wnd.SetBackgroundColor(tcell.GetColor(entity.SelectedTheme.Colors["WindowColor"]))
	listMethods.SetBackgroundColor(wnd.GetBackgroundColor())

	for i, method := range c.conn.AvailableMethods {
		listMethods.AddItem(method, "", 0, nil)

		// Ignore if none was selected before
		if c.conn.SelectedMethod != nil && *c.conn.SelectedMethod == method {
			listMethods.SetCurrentItem(i)
		}
	}

	listMethods.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		case tcell.KeyEnter:
			idx := listMethods.GetCurrentItem()
			item, _ := listMethods.GetItemText(idx)
			c.conn.SelectedMethod = &item

			// Remove the window and restore focus to menu list
			c.PrintLog("👉 Method set to [blue]"+*c.conn.SelectedMethod, LOG_INFO)
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		}
		return event
	})

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 70, 7)
	wnd.AddButton(&winman.Button{
		Symbol: '❌',
		OnClick: func() {
			c.ui.WinMan.RemoveWindow(wnd)
			c.ui.SetFocus(c.ui.MenuList)
		},
	})

	var maxMinButton *winman.Button
	maxMinButton = &winman.Button{
		Symbol:    '▴',
		Alignment: winman.ButtonRight,
		OnClick: func() {
			if wnd.IsMaximized() {
				wnd.Restore()
				maxMinButton.Symbol = '▴'
			} else {
				wnd.Maximize()
				maxMinButton.Symbol = '▾'
			}
		},
	}
	wnd.AddButton(maxMinButton)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

func (c Controller) setRequestPayload() {
	if c.conn.ActiveConnection == nil {
		c.PrintLog("❗ no active connection", LOG_WARNING)
		return
	}

	requestPayload := c.conn.RequestPayload

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetTitle(" Request Payload ").
		SetDraggable(true)

	wnd.SetBackgroundColor(tcell.GetColor(entity.SelectedTheme.Colors["WindowColor"]))
	form.SetBackgroundColor(wnd.GetBackgroundColor())

	form.SetBorderPadding(1, 1, 0, 1)

	// Create text area for filling the payload
	txtPayload := tview.NewTextArea().SetText(requestPayload, true)
	txtPayload.SetSize(9, 100)
	txtPayload.SetTextStyle(tcell.StyleDefault.Background(tcell.Color100).Foreground(tcell.Color102))
	form.AddFormItem(txtPayload)

	form.AddButton("Generate Sample", func() {
		if c.conn.SelectedMethod == nil {
			c.PrintLog("please select grpc method first", LOG_ERROR)
			return
		}
		// Get service detail
		dsc, err := c.conn.DescriptorSource.FindSymbol(*c.conn.SelectedMethod)
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
			return
		}

		txt, err := grpcurl.GetDescriptorText(dsc, c.conn.DescriptorSource)
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
			return
		}

		// Parse the service to get request message name
		rr := c.parseRequestResponse(txt)
		if len(rr) < 2 {
			c.PrintLog(fmt.Sprintf("failed to parse service name: %s", txt), LOG_ERROR)
		}
		// Remove stream from request
		requestMessage := strings.ReplaceAll(rr[0][1], "stream", "")

		// Trim message
		requestMessage = strings.TrimSpace(requestMessage)
		if requestMessage[0:1] == "." {
			requestMessage = requestMessage[1:]
		}

		// Retrieve request message from descriptors
		dsc, err = c.conn.DescriptorSource.FindSymbol(requestMessage)
		if err != nil {
			c.PrintLog(err.Error(), LOG_ERROR)
			return
		}

		if dsc, ok := dsc.(*desc.MessageDescriptor); ok {
			// for messages, also show a template in JSON, to make it easier to
			// create a request to invoke an RPC
			tmpl := grpcurl.MakeTemplate(dsc)
			options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
			_, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), c.conn.DescriptorSource, nil, options)
			if err != nil {
				c.PrintLog(err.Error(), LOG_ERROR)
				return
			}
			str, err := formatter(tmpl)
			if err != nil {
				c.PrintLog(err.Error(), LOG_ERROR)
				return
			}
			txtPayload.SetText(str, true)
		}
	})

	form.AddButton("Set", func() {
		c.conn.RequestPayload = txtPayload.GetText()

		// Remove the window and restore focus to menu list
		c.PrintLog("\nRequest Payload:\n[yellow]"+c.conn.RequestPayload, LOG_INFO)
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})

	form.AddButton("Cancel", func() {
		// Remove the window and restore focus to menu list
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})
	form.SetButtonsAlign(tview.AlignRight)

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 70, 15)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}
