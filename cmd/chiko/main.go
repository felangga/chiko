package main

import (
	"log"
	"os"

	"github.com/felangga/chiko/internal/controller"
	"github.com/felangga/chiko/internal/entity"
	"github.com/felangga/chiko/internal/ui"
)

// Version is injected at build time via -ldflags "-X main.Version=<tag>"
var Version string

func main() {
	if Version != "" {
		entity.APP_VERSION = Version
	}

	flags, err := controller.ParseFlags()
	if err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}

	err = ui.NewUI(flags).Run()
	if err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
