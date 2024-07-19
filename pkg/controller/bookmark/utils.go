package bookmark

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetOSConfigDir is used to get the path of the config directory based on the operating system.
func GetOSConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA")
	case "darwin":
		return filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			return xdgConfigHome
		}
		return filepath.Join(os.Getenv("HOME"), ".config")
	default:
		return ""
	}
}
