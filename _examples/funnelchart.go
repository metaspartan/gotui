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

	fc := widgets.NewFunnelChart()
	fc.Title = "Sales Funnel"
	fc.SetRect(0, 0, 50, 30)
	fc.Data = []float64{1000, 750, 500, 250, 50}
	fc.Labels = []string{"Leads", "Interested", "Trial", "Negotiation", "Won"}
	fc.Colors = []ui.Color{
		ui.ColorBlue,
		ui.ColorCyan,
		ui.ColorGreen,
		ui.ColorYellow,
		ui.ColorRed,
	}

	ui.Render(fc)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
