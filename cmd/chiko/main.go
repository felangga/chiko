package main

import (
	"log"
	"os"

	"github.com/felangga/chiko/internal/ui"
)

func main() {
	err := ui.NewUI().Run()
	if err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}
