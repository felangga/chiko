package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/controller/grpc"
	"github.com/felangga/chiko/internal/entity"
)

// CreateSessionWindow spawns a robust MDI window holding request state
func (u *UI) CreateSessionWindow(existingSession *entity.Session) {
	var session *entity.Session
	
	if existingSession != nil {
		session = existingSession
		if session.ID == uuid.Nil {
			session.ID = uuid.New()
		}
	} else {
		session = &entity.Session{
			ID:                 uuid.New(),
			AllowUnknownFields: true,
			RequestPayload:     "{\n}",
		}
	}

	grpcCtrl := grpc.NewGRPC(u.Logger, session)
	
	// Create output panel
	outputPanel := u.InitOutputPanel()
	
	sw := &SessionWindow{
		ID:          session.ID,
		GRPC:        &grpcCtrl,
		OutputPanel: outputPanel,
	}

	// Build the Postman-style session layout — top bar calls focusHook internally
	layout := u.setupSessionLayout(sw)

	serverUrl := session.ServerURL
	if serverUrl == "" {
		serverUrl = "New Request"
	}

	w := u.WinMan.NewWindow().
		Show().
		SetRoot(layout).
		SetBorder(true).
		SetDraggable(true).
		SetResizable(true).
		SetTitle(fmt.Sprintf(" %s ", serverUrl))

	// Start at a comfortable default size (not maximized).
	// Position slightly offset so multiple windows stack visually.
	offset := len(u.Sessions) * 2
	w.SetRect(5+offset, 2+offset, 140, 38)

	// ── Title Bar Buttons ────────────────────────────────────
	// Maximize / Restore toggle button (right-aligned)
	w.AddButton(&winman.Button{
		Symbol:    '▣',
		Alignment: winman.ButtonRight,
		OnClick: func() {
			if w.IsMaximized() {
				w.Restore()
			} else {
				w.Maximize()
			}
		},
	})

	sw.WinBase = w

	u.Sessions = append(u.Sessions, sw)
	u.ActiveSession = sw
	u.GRPC = sw.GRPC

	// Crucial: expressly send the Dashboard back to the absolute bottom so older windows aren't hidden
	if u.HomeWindow != nil {
		u.WinMan.SetZ(u.HomeWindow, winman.WindowZBottom)
	}

	// explicitly transfer focus to the new window so it pops to the front and accepts input
	u.SetFocus(u.activeSessionFocus())
}

// CloseSession cleanly removes a session window
func (u *UI) CloseSession(sw *SessionWindow) {
	u.WinMan.RemoveWindow(sw.WinBase)
	for i, s := range u.Sessions {
		if s.ID == sw.ID {
			u.Sessions = append(u.Sessions[:i], u.Sessions[i+1:]...)
			break
		}
	}
	
	// Update active pointer safely
	if len(u.Sessions) > 0 {
		u.ActiveSession = u.Sessions[len(u.Sessions)-1]
		u.GRPC = u.ActiveSession.GRPC
	} else {
		u.ActiveSession = nil
		// Return to home menu automatically
		// Since WinMan is root, closing all windows reveals background Home Dashboard!
	}
}

func (u *UI) setupSessionLayout(sw *SessionWindow) *tview.Flex {
	topBar := u.initSessionTopBar(sw)
	requestPanel, tabBar := u.initSessionRequestPanel(sw)

	// Output panel also sets focus hook
	sw.OutputPanel.TextArea.SetFocusFunc(func() {
		u.GRPC = sw.GRPC
		u.ActiveSession = sw
	})

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topBar, 3, 0, true).                  // Row 1 & 2: URL bar + Send, Method Field + Margin
		AddItem(tabBar, 1, 0, false).                 // Row 4: Body / Metadata / Auth tabs
		AddItem(requestPanel, 0, 2, false).           // Row 3: swappable request panel
		AddItem(sw.OutputPanel.Layout, 0, 3, false)   // Row 4: response/output

	return layout
}

func (u *UI) FocusNextSession() {
    if len(u.Sessions) == 0 { return }
    idx := -1
    for i, s := range u.Sessions {
        if s == u.ActiveSession { idx = i; break }
    }
    if idx == -1 { idx = 0 }
    
    nextIdx := (idx + 1) % len(u.Sessions)
    u.ActiveSession = u.Sessions[nextIdx]
    u.GRPC = u.ActiveSession.GRPC

    if u.HomeWindow != nil {
		u.WinMan.SetZ(u.HomeWindow, winman.WindowZBottom)
	}
    u.ActiveSession.WinBase.Show()
    u.SetFocus(u.activeSessionFocus())
}

// CycleFocus shifts keyboard focus cleanly around the SessionWindow components
func (sw *SessionWindow) CycleFocus(u *UI, step int) {
	cycle := []tview.Primitive{
		sw.URLField,
		sw.MethodField,
		sw.SendBtn,
		sw.RequestTabBar,
		sw.RequestBodyArea,
		sw.OutputPanel.TextArea,
	}

	currentFocus := u.App.GetFocus()
	focusIdx := -1
	for i, p := range cycle {
		if currentFocus == p {
			focusIdx = i
			break
		}
	}

	if focusIdx == -1 {
		// Default to URL field if not tracked
		u.SetFocus(cycle[0])
		return
	}

	nextIdx := (focusIdx + step + len(cycle)) % len(cycle)
	u.SetFocus(cycle[nextIdx])
}

func (u *UI) FocusPrevSession() {
    if len(u.Sessions) == 0 { return }
    idx := -1
    for i, s := range u.Sessions {
        if s == u.ActiveSession { idx = i; break }
    }
    
    prevIdx := idx - 1
    if prevIdx < 0 { prevIdx = len(u.Sessions) - 1 }
    
    u.ActiveSession = u.Sessions[prevIdx]
    u.GRPC = u.ActiveSession.GRPC
    u.ActiveSession.WinBase.Show()
}
