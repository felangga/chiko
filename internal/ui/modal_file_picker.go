package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	modalWnd   *winman.WindowBase
	fileList   *tview.TreeView
	currentDir string = "/"
	treeRoot   *tview.TreeNode
)

// sortDirectoryContents sorts directory contents with directories first, then files
// and sorts alphabetically within each group
func sortDirectoryContents(contents []string, baseDir string) []string {
	sort.Slice(contents, func(i, j int) bool {
		iPath := filepath.Join(baseDir, contents[i])
		jPath := filepath.Join(baseDir, contents[j])

		iIsDir, errI := isDirectory(iPath)
		jIsDir, errJ := isDirectory(jPath)

		// Prioritize directories
		if errI == nil && errJ == nil && iIsDir != jIsDir {
			return iIsDir
		}

		// If both are the same type (both dir or both file), sort alphabetically
		return strings.ToLower(contents[i]) < strings.ToLower(contents[j])
	})
	return contents
}

func (u UI) ShowModalFilePicker(parentWnd tview.Primitive, curDir string, onFileSelected func(string)) {
	// Initialize file list view
	fileList = u.createFileListView()

	// Create modal window for file picker
	modalWnd = u.createFilePickerModal(fileList, parentWnd)

	// Determine starting directory
	currentDir = u.determineStartingDirectory(curDir)

	// Set up input capture and load file tree
	u.ModalFilePicker_SetInputCapture(onFileSelected, parentWnd)
	u.loadFileTree()
}

// createFileListView sets up the tree view for file selection
func (u UI) createFileListView() *tview.TreeView {
	treeRoot = tview.NewTreeNode("🏠 /")
	fileList = tview.NewTreeView()
	fileList.SetBackgroundColor(u.Theme.Colors.WindowColor)
	fileList.SetRoot(treeRoot)
	fileList.SetCurrentNode(treeRoot)

	treeRoot.SetTextStyle(tcell.StyleDefault.
		Background(u.Theme.Colors.WindowColor))

	return fileList
}

// createFilePickerModal creates the modal dialog for file selection
func (u UI) createFilePickerModal(fileList *tview.TreeView, parentWnd tview.Primitive) *winman.WindowBase {
	return u.CreateModalDialog(CreateModalDialogParam{
		title:         " Select File ",
		rootView:      fileList,
		draggable:     true,
		resizeable:    true,
		size:          winSize{0, 0, 50, 20},
		fallbackFocus: parentWnd,
	})
}

// determineStartingDirectory resolves the initial directory for file selection
func (u UI) determineStartingDirectory(curDir string) string {
	// Get user home directory with robust fallback
	homeDir, err := os.Getwd()
	if err != nil {
		log.Printf("Could not get home directory: %v. Falling back to root.", err)
		homeDir = "/"
	}

	// If curDir is empty, return home directory
	if curDir == "" {
		return homeDir
	}

	// Check if the path is a file or directory
	info, err := os.Stat(curDir)
	if err != nil {
		// If path is invalid, fall back to home directory
		log.Printf("Could not stat path %s: %v. Falling back to home directory.", curDir, err)
		return homeDir
	}

	// If it's a file, extract its directory
	if !info.IsDir() {
		return filepath.Dir(curDir)
	}

	// If it's already a directory, return it
	return curDir
}

func (u UI) loadFileTree() {
	fileList.GetRoot().ClearChildren()

	// Update root node with current directory, truncate if too long
	displayDir := currentDir
	if len(displayDir) > 30 {
		displayDir = "..." + displayDir[len(displayDir)-27:]
	}
	treeRoot.SetText(fmt.Sprintf("🐶 %s", displayDir))

	// Load and sort directory contents
	dirContents := loadDirectoryContent(currentDir)
	sortedContents := sortDirectoryContents(dirContents, currentDir)

	// Create nodes for each file/directory
	for _, item := range sortedContents {
		var icon string
		path := filepath.Join(currentDir, item)

		isDir, err := isDirectory(path)
		if err != nil {
			log.Printf("Error checking path %s: %v", path, err)
			continue
		}

		if isDir {
			icon = "📁"
		} else {
			icon = "📄"
		}

		child := tview.NewTreeNode(fmt.Sprintf("%s %s", icon, item))
		child.SetReference(path)
		child.SetTextStyle(tcell.StyleDefault.
			Background(u.Theme.Colors.WindowColor))

		fileList.GetRoot().AddChild(child)
	}
}

// Helper function to check if a path is a directory with improved error handling
func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("path does not exist: %s", path)
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func loadDirectoryContent(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Error reading directory %s: %v", dir, err)
		return []string{}
	}

	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		// Optionally, you can add filters here (e.g., ignore hidden files)
		if !strings.HasPrefix(entry.Name(), ".") {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames
}

func (u *UI) ModalFilePicker_SetInputCapture(onFileSelected func(string), parentWnd tview.Primitive) {

	fileList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			u.CloseModalDialog(modalWnd, parentWnd)
			return nil
		}
		return event
	})

	fileList.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()

		// Root node or no reference - go to parent directory
		if reference == nil {
			parentDir := filepath.Dir(currentDir)

			// Prevent going above root
			if parentDir == currentDir {
				return
			}

			currentDir = parentDir
			u.loadFileTree()
			return
		}

		path := reference.(string)

		// Check if it's a directory
		isDir, err := isDirectory(path)
		if err != nil {
			log.Printf("Error checking directory: %v", err)
			return
		}

		if isDir {
			// If no children, load directory contents
			if len(node.GetChildren()) == 0 {
				currentDir = path
				u.loadFileTree()
			}
		} else {
			// Close modal and return focus to parent window
			u.CloseModalDialog(modalWnd, parentWnd)
			onFileSelected(path)
		}
	})
}
