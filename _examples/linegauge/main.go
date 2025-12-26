package main

import (
	"log"
	"time"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
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

	g3 := widgets.NewLineGauge()
	g3.Title = "Line Gauge (Block)"
	g3.SetRect(0, 8, 50, 11)
	g3.Percent = 50
	g3.BarRune = '■'
	g3.LineColor = ui.ColorYellow

	g4 := widgets.NewLineGauge()
	g4.Title = "Line Gauge (Dots)"
	g4.SetRect(0, 12, 50, 15)
	g4.Percent = 60
	g4.BarRune = '⣿'
	g4.LineColor = ui.ColorMagenta
	g4.TitleStyle.Fg = ui.ColorWhite

	g5 := widgets.NewLineGauge()
	g5.Title = "Line Gauge (Future)"
	g5.SetRect(0, 16, 50, 19)
	g5.Percent = 45
	g5.BarRune = '▰'
	g5.BarRuneEmpty = '▱'
	g5.LineColor = ui.ColorLightCyan
	g5.LabelStyle = ui.NewStyle(ui.ColorLightCyan)

	ui.Render(g1, g2, g3, g4, g5)

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
			g3.Percent = (g3.Percent + 9) % 100
			g4.Percent = (g4.Percent + 4) % 100
			g5.Percent = (g5.Percent + 3) % 100
			ui.Render(g1, g2, g3, g4, g5)
		}
	}
}
