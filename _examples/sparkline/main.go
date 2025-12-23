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

	data := []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}

	sl0 := widgets.NewSparkline()
	sl0.Data = data[3:]
	sl0.LineColor = ui.ColorGreen

	// single
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = "Sparkline 0"
	slg0.SetRect(0, 0, 20, 10)

	sl1 := widgets.NewSparkline()
	sl1.Title = "Sparkline 1"
	sl1.Data = data
	sl1.LineColor = ui.ColorRed

	sl2 := widgets.NewSparkline()
	sl2.Title = "Sparkline 2"
	sl2.Data = data[5:]
	sl2.LineColor = ui.ColorMagenta

	slg1 := widgets.NewSparklineGroup(sl0, sl1, sl2)
	slg1.Title = "Group Sparklines"
	slg1.SetRect(0, 10, 25, 25)

	sl3 := widgets.NewSparkline()
	sl3.Title = "Enlarged Sparkline"
	sl3.Data = data
	sl3.LineColor = ui.ColorYellow

	slg2 := widgets.NewSparklineGroup(sl3)
	slg2.Title = "Tweeked Sparkline"
	slg2.SetRect(20, 0, 50, 10)
	slg2.BorderStyle.Fg = ui.ColorLightCyan

	// Sparkline with explicit black background
	sl4 := widgets.NewSparkline()
	sl4.Title = "Black Background Sparkline"
	sl4.Data = data
	sl4.LineColor = ui.ColorCyan
	sl4.BackgroundColor = ui.ColorBlack

	slg3 := widgets.NewSparklineGroup(sl4)
	slg3.Title = "With Background Color"
	slg3.SetRect(25, 10, 55, 20)
	slg3.BackgroundColor = ui.ColorBlack

	ui.Render(slg0, slg1, slg2, slg3)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
