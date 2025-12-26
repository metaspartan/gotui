package main

import (
	"log"
	"math/rand"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	h := widgets.NewHeatmap()
	h.Title = "Server Load Heatmap"

	// Generate random data
	// 10 rows, 10 cols
	data := make([][]float64, 15)
	for i := range data {
		row := make([]float64, 10)
		for j := range row {
			row[j] = rand.Float64() * 100
		}
		data[i] = row
	}
	h.Data = data
	h.CellWidth = 3
	h.YLabels = []string{"Node1", "Node2", "Node3", "Node4", "Node5", "Node6", "Node7", "Node8", "Node9", "Node10"}
	h.Colors = []ui.Color{
		ui.NewRGBColor(0, 255, 0),   // Green
		ui.NewRGBColor(255, 255, 0), // Yellow
		ui.NewRGBColor(255, 165, 0), // Orange
		ui.NewRGBColor(255, 0, 0),   // Red
	}

	termWidth, termHeight := ui.TerminalDimensions()
	h.SetRect(0, 0, termWidth, termHeight)

	ui.Render(h)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			h.SetRect(0, 0, payload.Width, payload.Height)
			ui.Render(h)
		}
	}
}
