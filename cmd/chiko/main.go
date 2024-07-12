package main

import (
	"chiko/pkg/ui"
	"fmt"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	err := ui.NewUI().Run()
	if err != nil {
		panic(err)
	}
}
