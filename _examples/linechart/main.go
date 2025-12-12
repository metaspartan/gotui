package main

import (
	"log"
	"math"
	"time"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	// 1. Basic Sine Wave (Braille)
	p1 := widgets.NewPlot()
	p1.Title = "Braille Line Chart (Sine Wave)"
	p1.Data = make([][]float64, 2)
	p1.Data[0] = make([]float64, 100)
	p1.Data[1] = make([]float64, 100)
	p1.AxesColor = ui.ColorWhite
	p1.LineColors[0] = ui.ColorCyan
	p1.LineColors[1] = ui.ColorYellow
	p1.Marker = widgets.MarkerBraille
	p1.Fill = true // Enable filled area mode // Default, gives high resolution lines

	// 2. Dot Mode Comparison
	p2 := widgets.NewPlot()
	p2.Title = "Dot Mode (Same Data)"
	p2.Data = make([][]float64, 2)
	p2.AxesColor = ui.ColorWhite
	p2.LineColors[0] = ui.ColorCyan
	p2.LineColors[1] = ui.ColorYellow
	p2.Marker = widgets.MarkerDot

	// Grid layout
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0,
			ui.NewCol(0.5, p1),
			ui.NewCol(0.5, p2),
		),
	)

	update := func(tick int) {
		for i := 0; i < 100; i++ {
			p1.Data[0][i] = math.Sin(float64(i+tick) / 10)
			p1.Data[1][i] = math.Cos(float64(i+tick) / 10)
		}
		p2.Data = p1.Data
	}

	update(0)
	ui.Render(grid)

	ticker := time.NewTicker(50 * time.Millisecond).C
	uiEvents := ui.PollEvents()

	count := 0
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			count++
			update(count)
			ui.Render(grid)
		}
	}
}
