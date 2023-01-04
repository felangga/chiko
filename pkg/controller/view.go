package controller

import (
	"chiko/pkg/entity"
	"fmt"
	"strings"

	"github.com/epiclabs-io/winman"
	"github.com/fullstorydev/grpcurl"
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
}

func (c Controller) initSys() {
	c.PrintLog(fmt.Sprintf("✨ Welcome to Chiko v%s", entity.APP_VERSION), LOG_INFO)
}

func (c Controller) setServerURL() {
	tmpURL := c.conn.ServerURL

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetDraggable(true)

	form.AddInputField("Server URL", c.conn.ServerURL, 0, nil, func(txt string) {
		tmpURL = txt
	})

	form.AddButton("Set", func() {
		c.conn.ServerURL = tmpURL
		// Remove the window and restore focus to menu list
		c.PrintLog("Server URL set to [blue]"+c.conn.ServerURL, LOG_INFO)
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)

		c.CheckGRPC()
	})

	form.AddButton("Cancel", func() {
		// Remove the window and restore focus to menu list
		c.ui.WinMan.RemoveWindow(wnd)
		c.ui.SetFocus(c.ui.MenuList)
	})
	form.SetButtonsAlign(tview.AlignRight)

	wnd.SetModal(true)
	wnd.SetRect(0, 0, 50, 7)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

func (c Controller) setMethod() {
	selectedMethod := c.conn.SelectedMethod

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetDraggable(true)

	form.AddDropDown("Select Method", c.conn.AvailableMethods, 0, func(option string, optionIndex int) {
		selectedMethod = &option
	}).SetBorderPadding(1, 1, 1, 1)

	form.AddButton("Set", func() {
		c.conn.SelectedMethod = selectedMethod

		// Remove the window and restore focus to menu list
		c.PrintLog(" 👉 Method set to [blue]"+*c.conn.SelectedMethod, LOG_INFO)
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
	wnd.SetRect(0, 0, 70, 7)

	c.ui.WinMan.AddWindow(wnd)
	c.ui.WinMan.Center(wnd)
	c.ui.SetFocus(wnd)
}

func (c Controller) setRequestPayload() {
	requestPayload := c.conn.RequestPayload

	// Create Set Server URL From
	form := tview.NewForm()
	wnd := winman.NewWindow().
		Show().
		SetRoot(form).
		SetTitle("Request Payload").
		SetDraggable(true)

	form.SetBorderPadding(1, 1, 0, 1)
	txtPayload := tview.NewTextArea().SetText(requestPayload, true)
	txtPayload.SetSize(9, 100)

	// form.AddTextArea("", requestPayload, 0, 9, 0, nil)
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
