package main

import (
	"log"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	rc := widgets.NewRadarChart()
	rc.Title = "System Health"
	rc.SetRect(0, 0, 50, 25)
	rc.Data = [][]float64{
		{80, 40, 30, 90, 60}, // Current
		{50, 50, 50, 50, 50}, // Avg
	}
	rc.Labels = []string{"CPU", "Mem", "Disk", "Net", "Swap"}
	rc.LineColors = []ui.Color{ui.ColorGreen, ui.ColorRed}
	rc.MaxVal = 100

	ui.Render(rc)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
