//go:build ignore
// +build ignore

package main

import (
	"log"

	ui "github.com/metaspartan/gotui"
	"github.com/metaspartan/gotui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewLogo()
	l.Title = "GOTUI"
	l.SetRect(0, 0, 50, 10)
	l.BorderStyle.Fg = ui.ColorCyan

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.Type == ui.KeyboardEvent {
			return
		}
	}
}
