//go:build ignore
// +build ignore

package main

import (
	"log"
	"time"

	ui "github.com/metaspartan/gotui"
	"github.com/metaspartan/gotui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	wm := widgets.NewWorldMap()
	wm.Title = "World Map (Server Locations)"
	wm.SetRect(0, 0, 100, 30)
	wm.Color = ui.ColorCyan

	// Sample data: Servers in US, Europe, Asia
	wm.Data = []widgets.MapPoint{
		{Lat: 40.7128, Lon: -74.0060, Color: ui.ColorRed},     // New York
		{Lat: 37.7749, Lon: -122.4194, Color: ui.ColorRed},    // San Francisco
		{Lat: 51.5074, Lon: -0.1278, Color: ui.ColorGreen},    // London
		{Lat: 48.8566, Lon: 2.3522, Color: ui.ColorGreen},     // Paris
		{Lat: 35.6895, Lon: 139.6917, Color: ui.ColorYellow},  // Tokyo
		{Lat: -33.8688, Lon: 151.2093, Color: ui.ColorYellow}, // Sydney
	}

	ui.Render(wm)

	// Simple animation loop adding random points
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker.C:
			// Blink effect or move points?
			// Just re-render for now
			ui.Render(wm)
		}
	}
}
