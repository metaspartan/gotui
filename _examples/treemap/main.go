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

	tm := widgets.NewTreeMap()
	tm.Title = "Disk Usage"
	tm.SetRect(0, 0, 80, 24)
	tm.TextColor = ui.ColorWhite

	// Sample Hierarchical Data
	root := &widgets.TreeMapNode{
		Children: []*widgets.TreeMapNode{
			{
				Label: "/var",
				Value: 40,
				Style: ui.NewStyle(ui.ColorWhite, ui.ColorRed),
			},
			{
				Label: "/home",
				Value: 30,
				Style: ui.NewStyle(ui.ColorWhite, ui.ColorBlue),
			},
			{
				Label: "/usr",
				Value: 20,
				Style: ui.NewStyle(ui.ColorWhite, ui.ColorGreen),
			},
			{
				Label: "/etc",
				Value: 10,
				Style: ui.NewStyle(ui.ColorWhite, ui.ColorYellow),
			},
		},
	}
	tm.Root = root

	ui.Render(tm)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
