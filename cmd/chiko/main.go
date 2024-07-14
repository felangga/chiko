package main

import (
	"chiko/pkg/ui"
)

func main() {
	err := ui.NewUI().Run()
	if err != nil {
		panic(err)
	}
}
