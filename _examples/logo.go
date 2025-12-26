//go:build ignore
// +build ignore

package main

import (
	"log"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewLogo()
	l.Title = "GOTUI"
	l.SetRect(0, 0, 50, 11)
	l.BorderStyle.Fg = ui.ColorLightCyan
	l.TitleBottom = "Go TUI Library By Carsen Klock"
	l.TitleBottomAlignment = ui.AlignCenter

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.Type == ui.KeyboardEvent {
			return
		}
	}
}
