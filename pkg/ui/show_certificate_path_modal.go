package ui

import (
	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/pkg/entity"
)

var (
	txtCAPath   *tview.InputField
	txtCertPath *tview.InputField
	txtKeyPath  *tview.InputField
)

func (u *UI) ShowCertificatePathModal(parentWnd *winman.WindowBase) {
	var (
		CAPath   string
		CertPath string
		KeyPath  string
	)

	if u.GRPC.Conn.SSLCert != nil {
		if u.GRPC.Conn.SSLCert.CA_Path != nil {
			CAPath = *u.GRPC.Conn.SSLCert.CA_Path
		}
		if u.GRPC.Conn.SSLCert.ClientCert_Path != nil {
			CertPath = *u.GRPC.Conn.SSLCert.ClientCert_Path
		}
		if u.GRPC.Conn.SSLCert.ClientKey_Path != nil {
			KeyPath = *u.GRPC.Conn.SSLCert.ClientKey_Path
		}
	}

	txtCAPath = tview.NewInputField()
	txtCAPath.SetBackgroundColor(u.Theme.Colors.WindowColor)
	txtCAPath.SetFieldStyle(u.Theme.Style.FieldStyle)
	txtCAPath.SetText(CAPath)

	txtCertPath = tview.NewInputField()
	txtCertPath.SetLabel("Cert Path")
	txtCertPath.SetFieldStyle(u.Theme.Style.FieldStyle)
	txtCertPath.SetText(CertPath)

	txtKeyPath = tview.NewInputField()
	txtKeyPath.SetLabel("Key Path")
	txtKeyPath.SetText(KeyPath)

	layout := tview.NewForm()
	layout.SetBorderPadding(1, 1, 1, 1)
	layout.SetButtonsAlign(tview.AlignRight)
	layout.SetBackgroundColor(u.Theme.Colors.WindowColor)
	layout.SetButtonStyle(u.Theme.Style.ButtonStyle)
	layout.SetFieldStyle(u.Theme.Style.FieldStyle)

	layout.AddCheckbox("Enable TLS", u.GRPC.Conn.EnableTLS, func(checked bool) {
		u.GRPC.Conn.EnableTLS = checked
	})
	layout.AddCheckbox("Skip Verification", u.GRPC.Conn.InsecureSkipVerify, func(checked bool) {
		u.GRPC.Conn.InsecureSkipVerify = checked
	})
	layout.AddTextView("CA Path", "Provide CA Certificate if your server uses a self-signed certificate\n(Press Enter to Browse File)", 80, 2, true, false)
	layout.AddFormItem(txtCAPath)
	layout.AddTextView("mTLS Key", "Provide Client Certificate and Key if your server uses mutual TLS\n(Press Enter to Browse File)", 80, 2, true, false)
	layout.AddFormItem(txtCertPath)
	layout.AddFormItem(txtKeyPath)

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 🔐 TLS ",
		rootView:      layout,
		draggable:     true,
		resizeable:    false,
		size:          winSize{0, 0, 100, 21},
		fallbackFocus: parentWnd,
	})

	layout.AddButton("Cancel", func() {
		u.WinMan.RemoveWindow(wnd)
		u.SetFocus(parentWnd)
	})

	layout.AddButton("Apply", func() {
		u.GRPC.Conn.SSLCert = &entity.Cert{}

		// Get paths from input fields
		CAPath := txtCAPath.GetText()
		CertPath := txtCertPath.GetText()
		KeyPath := txtKeyPath.GetText()

		// Assign paths if not empty
		if CAPath != "" {
			u.GRPC.Conn.SSLCert.CA_Path = &CAPath
		}
		if CertPath != "" {
			u.GRPC.Conn.SSLCert.ClientCert_Path = &CertPath
		}
		if KeyPath != "" {
			u.GRPC.Conn.SSLCert.ClientKey_Path = &KeyPath
		}

		u.WinMan.RemoveWindow(wnd)
		u.SetFocus(parentWnd)
	})

	u.ShowCertificatePathModal_SetInputCapture(layout, wnd, parentWnd)
}

func (u *UI) ShowCertificatePathModal_SetInputCapture(form *tview.Form, wnd *winman.WindowBase, parentWnd *winman.WindowBase) {
	// Handle CA Path input box
	txtCAPath.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			u.ShowModalFilePicker(form, txtCAPath.GetText(), func(path string) {
				txtCAPath.SetText(path)
			})

			return nil
		}

		if event.Key() == tcell.KeyEscape {
			u.CloseModalDialog(wnd, parentWnd)
			return nil
		}

		return event
	})

	// Handle Cert Path input box
	txtCertPath.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			u.ShowModalFilePicker(form, txtCertPath.GetText(), func(path string) {
				txtCertPath.SetText(path)
			})
		}

		if event.Key() == tcell.KeyEscape {
			u.CloseModalDialog(wnd, parentWnd)
			return nil
		}

		return event
	})

	// Handle Key Path input box
	txtKeyPath.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			u.ShowModalFilePicker(form, txtKeyPath.GetText(), func(path string) {
				txtKeyPath.SetText(path)
			})
		}

		if event.Key() == tcell.KeyEscape {
			u.CloseModalDialog(wnd, parentWnd)
			return nil
		}

		return event
	})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, parentWnd)
		}
		return event
	})

}
