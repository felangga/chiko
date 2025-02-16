package main

import (
	"log"
	"os"

	"github.com/felangga/chiko/internal/controller"
	"github.com/felangga/chiko/internal/ui"
)

func main() {
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
