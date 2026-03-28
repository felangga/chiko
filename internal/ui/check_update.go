package ui

import (
	"fmt"

	"github.com/felangga/chiko/internal/controller"
	"github.com/felangga/chiko/internal/entity"
)

func (u *UI) checkForUpdates() {
	latest, isNewer, err := controller.CheckLatestVersion(entity.APP_VERSION)
	if err != nil {
		return
	}

	if isNewer {
		u.PrintLog(entity.Log{
			Content: fmt.Sprintf(
				"🆕 New version available: %s (current: %s) — https://github.com/felangga/chiko/releases",
				latest, entity.APP_VERSION,
			),
			Type: entity.LOG_INFO,
		})
	}
}
