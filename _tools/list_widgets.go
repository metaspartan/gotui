package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	examplesDir := "_examples"
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		panic(err)
	}

	var folderNames []string
	for _, f := range files {
		if f.IsDir() {
			folderNames = append(folderNames, f.Name())
		}
	}
	sort.Strings(folderNames)

	fmt.Println("| Widget/Example | Screenshot | Code |")
	fmt.Println("| :--- | :---: | :--- |")

	for _, folderName := range folderNames {
		// Skip empty or non-example dirs if any check needed?
		// check if main.go exists
		mainPath := filepath.Join(examplesDir, folderName, "main.go")
		if _, err := os.Stat(mainPath); os.IsNotExist(err) {
			continue
		}

		prettyName := strings.Title(strings.ReplaceAll(folderName, "_", " "))

		// Markdown row
		// using HTML img tag for height control
		fmt.Printf("| **%s** | <img src=\"_examples/%s/screenshot.png\" height=\"80\" /> | [View Source](_examples/%s/main.go) |\n",
			prettyName, folderName, folderName)
	}
}
