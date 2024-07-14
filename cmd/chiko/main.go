package main

import "github.com/felangga/chiko/pkg/ui"

func main() {
	err := ui.NewUI().Run()
	if err != nil {
		panic(err)
	}
}
