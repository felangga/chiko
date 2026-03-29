package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

// ShowHistoryModal opens the history panel as a modal dialog with a search box.
func (u *UI) ShowHistoryModal() {
	searchBox := tview.NewInputField().
		SetText("").
		SetPlaceholder(" 🔍 Search history...").
		SetPlaceholderStyle(u.Theme.Style.PlaceholderStyle)
	searchBox.SetFieldStyle(u.Theme.Style.FieldStyle)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBox, 1, 0, true).
		AddItem(u.Layout.HistoryPanel, 0, 1, false)
	layout.SetBackgroundColor(u.Theme.Colors.WindowColor)
	u.Layout.HistoryPanel.SetBackgroundColor(u.Theme.Colors.WindowColor)
	u.Layout.HistoryPanel.SetGraphicsColor(tcell.ColorWhite)

	if root := u.Layout.HistoryPanel.GetRoot(); root != nil {
		root.SetColor(tcell.ColorWhite)
	}

	wnd := u.CreateModalDialog(CreateModalDialogParam{
		title:         " 🕓 History ",
		rootView:      layout,
		draggable:     true,
		resizeable:    true,
		size:          winSize{0, 0, 80, 25},
		fallbackFocus: u.activeSessionFocus(),
	})
	wnd.SetBorderPadding(1, 0, 1, 1)

	u.filterHistoryPanel("")

	searchBox.SetChangedFunc(func(text string) {
		u.filterHistoryPanel(text)
	})

	searchBox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.activeSessionFocus())
			return nil
		case tcell.KeyTAB, tcell.KeyDown:
			u.App.SetFocus(u.Layout.HistoryPanel)
			return nil
		}
		return event
	})

	// ESC / TAB from the tree: go back to search box
	u.Layout.HistoryPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(wnd, u.activeSessionFocus())
			return nil
		case tcell.KeyTAB:
			// Move up to search box when at the top of the tree
			u.App.SetFocus(searchBox)
			return nil
		}
		return event
	})

	u.Layout.HistoryPanel.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		switch v := ref.(type) {
		case *entity.HistoryEntry:
			u.applyHistoryEntry(*v)
		default:
			node.SetExpanded(!node.IsExpanded())
		}
	})

	// Focus the search box first
	u.App.SetFocus(searchBox)
	_ = wnd
}

// filterHistoryPanel rebuilds the history tree filtered by query (case-insensitive).
// Matches on method name or server URL.
func (u *UI) filterHistoryPanel(query string) {
	u.History.Mu.RLock()
	entries := make([]entity.HistoryEntry, len(*u.History.Entries))
	copy(entries, *u.History.Entries)
	u.History.Mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))

	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	normalGroupStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(u.Theme.Colors.WindowColor)
	normalEntryStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(u.Theme.Colors.WindowColor)
	selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue).Bold(true)

	groupOrder := []string{}
	groupNodes := map[string]*tview.TreeNode{}

	for i := range entries {
		entry := &entries[i]

		method := "(no method)"
		if entry.Session.SelectedMethod != nil {
			method = *entry.Session.SelectedMethod
		}

		if q != "" {
			haystack := strings.ToLower(method + " " + entry.Session.ServerURL)
			if !strings.Contains(haystack, q) {
				continue
			}
		}

		local := entry.InvokedAt.Local()
		entryDay := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())

		var groupLabel string
		switch {
		case entryDay.Equal(today):
			groupLabel = "📅 Today"
		case entryDay.Equal(yesterday):
			groupLabel = "📅 Yesterday"
		default:
			groupLabel = "📅 " + entryDay.Format("02 January 2006")
		}

		if _, exists := groupNodes[groupLabel]; !exists {
			groupNode := tview.NewTreeNode(groupLabel).
				SetTextStyle(normalGroupStyle).
				SetSelectedTextStyle(selectedStyle)
			groupNode.SetExpanded(true)
			groupNode.SetReference(groupLabel)
			groupNodes[groupLabel] = groupNode
			groupOrder = append(groupOrder, groupLabel)
		}

		timeStr := local.Format("15:04:05")
		label := fmt.Sprintf("[white][%s] [yellow]%s [white]- %s", timeStr, entry.Session.ServerURL, method)
		entryNode := tview.NewTreeNode(label).
			SetTextStyle(normalEntryStyle).
			SetSelectedTextStyle(selectedStyle)
		entryNode.SetReference(entry)
		groupNodes[groupLabel].AddChild(entryNode)
	}

	go u.App.QueueUpdateDraw(func() {
		root := u.Layout.HistoryPanel.GetRoot()
		root.ClearChildren()
		for _, label := range groupOrder {
			root.AddChild(groupNodes[label])
		}
	})
}

// applyHistoryEntry loads a history entry into the active session and
// immediately initiates a new gRPC connection, mirroring ApplyBookmark.
func (u *UI) applyHistoryEntry(entry entity.HistoryEntry) {
	*u.GRPC.Conn = entry.Session

	method := "(no method)"
	if u.GRPC.Conn.SelectedMethod != nil {
		method = *u.GRPC.Conn.SelectedMethod
	}

	u.PrintLog(entity.Log{
		Content: fmt.Sprintf("🌏 connecting to [lightblue]%s [white]→ %s", u.GRPC.Conn.ServerURL, method),
		Type:    entity.LOG_INFO,
	})

	go func() {
		err := u.GRPC.Connect()
		if err != nil {
			u.PrintLog(entity.Log{
				Content: "❌ failed to connect to [blue]" + u.GRPC.Conn.ServerURL + " [red]" + err.Error(),
				Type:    entity.LOG_ERROR,
			})
			return
		}
	}()
}
