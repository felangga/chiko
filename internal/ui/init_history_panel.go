package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/felangga/chiko/internal/entity"
)

// InitHistoryPanel initializes the persistent history panel (right of the output area)
func (u *UI) InitHistoryPanel() *tview.TreeView {
	selStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	treeRoot := tview.NewTreeNode("🕓 History").
		SetTextStyle(tcell.StyleDefault.Background(u.Theme.Colors.WindowColor).Foreground(tcell.ColorWhite)).
		SetSelectedTextStyle(selStyle)

	historyPanel := tview.NewTreeView().
		SetRoot(treeRoot).
		SetCurrentNode(treeRoot)

	historyPanel.SetBorderPadding(1, 1, 1, 1)

	u.InitHistoryPanel_SetInputCapture(historyPanel)
	u.InitHistoryPanel_SetSelection(historyPanel)

	return historyPanel
}

func (u *UI) InitHistoryPanel_SetInputCapture(historyPanel *tview.TreeView) {
	historyPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			u.SetFocus(u.Layout.MenuList)
			return nil
		}
		return event
	})
}

func (u *UI) InitHistoryPanel_SetSelection(historyPanel *tview.TreeView) {
	historyPanel.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		switch v := ref.(type) {
		case *entity.HistoryEntry:
			u.applyHistoryEntry(*v)
		default:
			// Date group node — toggle expand/collapse
			node.SetExpanded(!node.IsExpanded())
		}
	})
}

// RefreshHistoryPanel rebuilds the history tree, grouping entries by local date.
func (u *UI) RefreshHistoryPanel() {
	u.History.Mu.RLock()
	entries := make([]entity.HistoryEntry, len(*u.History.Entries))
	copy(entries, *u.History.Entries)
	u.History.Mu.RUnlock()

	go u.App.QueueUpdateDraw(func() {
		root := u.Layout.HistoryPanel.GetRoot()
		root.ClearChildren()

		if len(entries) == 0 {
			return
		}

		now := time.Now().Local()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		yesterday := today.AddDate(0, 0, -1)

		groupOrder := []string{}
		groupNodes := map[string]*tview.TreeNode{}

		for i := range entries {
			entry := &entries[i]
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

			normalGroupStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(u.Theme.Colors.WindowColor)
			normalEntryStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(u.Theme.Colors.WindowColor)
			selectedStyle := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite).Bold(true)

			if _, exists := groupNodes[groupLabel]; !exists {
				groupNode := tview.NewTreeNode(groupLabel).
					SetTextStyle(normalGroupStyle).
					SetSelectedTextStyle(selectedStyle)
				groupNode.SetExpanded(true)
				groupNode.SetReference(groupLabel)
				groupNodes[groupLabel] = groupNode
				groupOrder = append(groupOrder, groupLabel)
			}

			method := "(no method)"
			if entry.Session.SelectedMethod != nil {
				method = *entry.Session.SelectedMethod
			}

			timeStr := local.Format("15:04:05")
			label := fmt.Sprintf("[%s] [yellow]%s [white]- %s", timeStr, entry.Session.ServerURL, method)

			entryNode := tview.NewTreeNode(label).
				SetTextStyle(normalEntryStyle).
				SetSelectedTextStyle(selectedStyle)
			entryNode.SetReference(entry)
			groupNodes[groupLabel].AddChild(entryNode)
		}

		for _, label := range groupOrder {
			root.AddChild(groupNodes[label])
		}
	})
}
