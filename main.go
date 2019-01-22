package main

import (
	"fmt"
	"os"

	"github.com/liamg/xu/app"
)

func main() {

	file, err := os.Open("./xu")
	if err != nil {
		fmt.Printf("Error: Cannot open file: %s\n", err)
		os.Exit(1)
	}

	editor, err := app.NewEditor(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	editor.Init()
	defer editor.Close()

	if err := editor.Run(); err != nil {
		editor.Close()
		fmt.Printf("xu crashed: %s\n", err)
		os.Exit(1)
	}
}
