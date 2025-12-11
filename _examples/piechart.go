//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"math/rand"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	pc := widgets.NewPieChart()
	pc.Title = "Pie Chart Example"
	pc.SetRect(0, 0, 50, 22)
	pc.Data = []float64{.25, .25, .25, .25}
	pc.AngleOffset = -.5 * 3.14159 // Start from top
	pc.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	termWidth, termHeight := ui.TerminalDimensions()
	pc.SetRect(0, 0, termWidth, termHeight)

	ui.Render(pc)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			pc.SetRect(0, 0, payload.Width, payload.Height)
			ui.Render(pc)
		case "<Space>":
			// Randomize data
			data := make([]float64, 4)
			for i := range data {
				data[i] = rand.Float64()
			}
			pc.Data = data
			ui.Render(pc)
		}
	}
}
