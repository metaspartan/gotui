package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const readmeTemplate = `# {{.Name}} Example

This example demonstrates the **{{.Name}}** widget/feature.

## 🚀 Run

` + "```bash" + `
go run _examples/{{.Folder}}/main.go
` + "```" + `

## 📸 Screenshot

![{{.Name}} Screenshot](screenshot.png)

## 📝 Code

See [main.go](main.go) for the implementation.
`

func main() {
	examplesDir := "_examples"
	files, err := os.ReadDir(examplesDir)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			continue
		}

		filename := f.Name()
		name := strings.TrimSuffix(filename, ".go")

		// Skip special files if any
		if name == "doc" || name == "go.mod" {
			continue
		}

		fmt.Printf("Migrating %s...\n", filename)

		// Create directory
		targetDir := filepath.Join(examplesDir, name)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Printf("Failed to create dir %s: %v\n", targetDir, err)
			continue
		}

		// Move file to targetDir/main.go
		srcPath := filepath.Join(examplesDir, filename)
		dstPath := filepath.Join(targetDir, "main.go")

		if err := moveFile(srcPath, dstPath); err != nil {
			fmt.Printf("Failed to move file %s: %v\n", srcPath, err)
			continue
		}

		// Create README.md
		readmePath := filepath.Join(targetDir, "README.md")
		readmeFile, err := os.Create(readmePath)
		if err != nil {
			fmt.Printf("Failed to create README %s: %v\n", readmePath, err)
			continue
		}

		data := struct {
			Name   string
			Folder string
		}{
			Name:   strings.Title(strings.ReplaceAll(name, "_", " ")),
			Folder: name,
		}

		if err := tmpl.Execute(readmeFile, data); err != nil {
			fmt.Printf("Failed to write README %s: %v\n", readmePath, err)
		}
		readmeFile.Close()
	}
	fmt.Println("Migration complete.")
}

func moveFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, input, 0644); err != nil {
		return err
	}
	return os.Remove(src)
}
