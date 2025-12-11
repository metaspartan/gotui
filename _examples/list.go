//go:build ignore
// +build ignore

package main

import (
	"log"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/metaspartan/gotui/v4",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 25, 8)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if handleListEvents(l, e, previousKey) {
			return
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}
}

func handleListEvents(l *widgets.List, e ui.Event, previousKey string) bool {
	switch e.ID {
	case "q", "<C-c>":
		return true
	case "j", "<Down>":
		l.ScrollDown()
	case "k", "<Up>":
		l.ScrollUp()
	case "<C-d>":
		l.ScrollHalfPageDown()
	case "<C-u>":
		l.ScrollHalfPageUp()
	case "<C-f>":
		l.ScrollPageDown()
		ui.Render(l) // Render calls in switch are redundant if main loop renders, but explicit calls might be intended for immediate update?
		// Actually main loop renders at end. Safe to remove explicit renders inside switch if logic allows.
		// Original code had ui.Render(l) in some cases.
	case "<MouseWheelUp>":
		l.ScrollUp()
		ui.Render(l)
	case "<MouseWheelDown>":
		l.ScrollDown()
		ui.Render(l)
	case "<C-b>":
		l.ScrollPageUp()
	case "g":
		if previousKey == "g" {
			l.ScrollTop()
		}
	case "<Home>":
		l.ScrollTop()
	case "G", "<End>":
		l.ScrollBottom()
	}
	return false
}
