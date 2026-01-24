package main

import (
	"fmt"
	"os"
	"os/exec"
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
		if !f.IsDir() {
			continue
		}
		folderName := f.Name()
		targetDir := filepath.Join(examplesDir, folderName)
		mainGo := filepath.Join(targetDir, "main.go")
		if _, err := os.Stat(mainGo); os.IsNotExist(err) {
			continue
		}
		fmt.Printf("Processing %s...\n", folderName)
		if err := stripBuildTags(mainGo); err != nil {
			fmt.Printf("⚠️ Failed to strip tags %s: %v\n", folderName, err)
		}
		args := []string{"run", ".", "-screenshot"}
		cmd := exec.Command("go", args...)
		cmd.Dir = targetDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Failed to screenshot %s: %v\n", folderName, err)
		} else {
			fmt.Printf("✅ Screenshot generated for %s\n", folderName)
		}
		readmePath := filepath.Join(targetDir, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			readmeFile, err := os.Create(readmePath)
			if err != nil {
				fmt.Printf("Failed to create README %s: %v\n", readmePath, err)
				continue
			}
			data := struct {
				Name   string
				Folder string
			}{
				Name:   strings.Title(strings.ReplaceAll(folderName, "_", " ")),
				Folder: folderName,
			}
			if err := tmpl.Execute(readmeFile, data); err != nil {
				fmt.Printf("Failed to write README %s: %v\n", readmePath, err)
			}
			readmeFile.Close()
			fmt.Printf("📝 README created for %s\n", folderName)
		}
	}
	fmt.Println("🎉 Generation complete.")
}
func stripBuildTags(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	var newLines []string
	changed := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:build ignore") ||
			strings.HasPrefix(trimmed, "// +build ignore") {
			changed = true
			continue
		}
		newLines = append(newLines, line)
	}
	if changed {
		return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
	}
	return nil
}
