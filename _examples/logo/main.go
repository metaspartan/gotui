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

	l := widgets.NewLogo()
	l.Title = "GOTUI"
	l.SetRect(0, 0, 50, 11)
	l.BorderStyle.Fg = ui.ColorLightCyan
	l.TitleBottom = "Go TUI Library By Carsen Klock"
	l.TitleBottomAlignment = ui.AlignCenter
	l.Gradient.Enabled = true
	l.Gradient.Stops = []ui.Color{
		ui.NewRGBColor(57, 255, 20), // Neon Green
		ui.NewRGBColor(255, 215, 0), // Gold
		ui.NewRGBColor(148, 0, 211), // Neon Purple
	}

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.Type == ui.KeyboardEvent {
			return
		}
	}
}
