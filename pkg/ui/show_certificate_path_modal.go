package ui

import (
	"os"

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
	txtCAPath.SetLabel("CA Path")
	txtCAPath.SetText(CAPath)

	txtCertPath = tview.NewInputField()
	txtCertPath.SetLabel("Cert Path")
	txtCertPath.SetText(CertPath)

	txtKeyPath = tview.NewInputField()
	txtKeyPath.SetLabel("Key Path")
	txtKeyPath.SetText(KeyPath)

	layout := tview.NewForm()
	layout.SetBorderPadding(1, 1, 1, 1)
	layout.SetButtonsAlign(tview.AlignRight)
	layout.SetBackgroundColor(u.Theme.Colors.WindowColor)
	layout.AddFormItem(txtCAPath)
	layout.AddFormItem(txtCertPath)
	layout.AddFormItem(txtKeyPath)

	wnd := u.CreateModalDialog(CreateModalDiaLog{
		title:         " 🔐 Certificate ",
		rootView:      layout,
		draggable:     true,
		size:          winSize{0, 0, 50, 11},
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

	u.ShowCertificatePathModal_SetInputCapture(wnd, parentWnd)
}

func (u *UI) ShowCertificatePathModal_SetInputCapture(wnd *winman.WindowBase, parentWnd *winman.WindowBase) {
	wnd.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.WinMan.RemoveWindow(wnd)
			u.SetFocus(parentWnd)
			os.Exit(0)
			return nil
		}

		return event
	})

}
