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

	// Set global background to Blue
	ui.ClearBackground(ui.ColorBlue)

	p := widgets.NewParagraph()
	p.Text = "The entire terminal background should be BLUE!"
	p.SetRect(5, 5, 50, 10)
	p.Border = true
	p.Title = "Global Background"

	// Style: Blue Background, Black Text, Black Border
	p.BackgroundColor = ui.ColorBlue
	p.TextStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)
	p.BorderStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)
	p.TitleStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)

	// Instructions
	info := widgets.NewParagraph()
	info.Text = "Press q to quit"
	info.SetRect(5, 12, 50, 15)
	info.BackgroundColor = ui.ColorBlue
	info.TextStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)
	info.BorderStyle = ui.NewStyle(ui.ColorBlack, ui.ColorBlue)

	ui.Render(p, info)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
