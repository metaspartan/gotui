//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/metaspartan/gotui"
	"github.com/metaspartan/gotui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	pc := widgets.NewPieChart()
	pc.Title = "Donut Chart"
	pc.SetRect(0, 0, 50, 25)
	pc.Data = []float64{10, 20, 30, 40}
	pc.AngleOffset = -90 // Standard start
	pc.InnerRadius = 0.5 // Donut mode!

	pc.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.0f", v)
	}

	ui.Render(pc)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
