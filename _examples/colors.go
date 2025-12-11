//go:build ignore
// +build ignore

package main

import (
	"log"
	"sort"

	"github.com/gdamore/tcell/v2"
	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	// Get all color names from tcell
	var names []string
	for name := range tcell.ColorNames {
		names = append(names, name)
	}
	sort.Strings(names)

	l := widgets.NewList()
	l.Title = "All Tcell Colors"
	l.Rows = names
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.WrapText = false
	coloredRows := make([]string, len(names))
	for i, name := range names {
		coloredRows[i] = "[" + name + "](fg:black,bg:" + name + ")"
	}
	l.Rows = coloredRows

	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(0, 0, termWidth, termHeight)

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
			ui.Render(l)
		case "k", "<Up>":
			l.ScrollUp()
			ui.Render(l)
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			l.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(l)
		case "<MouseWheelUp>":
			l.ScrollUp()
			ui.Render(l)
		case "<MouseWheelDown>":
			l.ScrollDown()
			ui.Render(l)
		}
	}
}
