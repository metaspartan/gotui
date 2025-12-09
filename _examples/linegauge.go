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

	g1 := widgets.NewLineGauge()
	g1.Title = "Line Gauge (Standard)"
	g1.SetRect(0, 0, 50, 3)
	g1.Percent = 50
	g1.LineColor = ui.ColorRed

	g2 := widgets.NewLineGauge()
	g2.Title = "Line Gauge (No Label)"
	g2.SetRect(0, 4, 50, 7)
	g2.Percent = 75
	g2.Label = " " // Hide label
	g2.LineColor = ui.ColorGreen

	ui.Render(g1, g2)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C

	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent {
				return
			}
		case <-ticker:
			g1.Percent = (g1.Percent + 5) % 100
			g2.Percent = (g2.Percent + 2) % 100
			ui.Render(g1, g2)
		}
	}
}
