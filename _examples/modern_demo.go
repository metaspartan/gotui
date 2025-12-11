package main

import (
	"fmt"
	"log"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	// 1. Block with Bottom Title
	p1 := widgets.NewParagraph()
	p1.Title = "Top Title"
	p1.TitleBottom = "Bottom Title (New Feature!)"
	p1.Text = "This block demonstrates the new bottom title feature.\nPress q to quit."
	p1.SetRect(5, 5, 45, 15)
	p1.BorderStyle.Fg = ui.ColorYellow

	// 2. Block with RGB Color (TrueColor)
	p2 := widgets.NewParagraph()
	p2.Title = "TrueColor Support"
	p2.Text = "This background should be a custom RGB color (Tele-ish).\nIf you see this, RGB is working!"
	p2.SetRect(50, 5, 90, 15)

	// Cobalt Blue / Teal custom color
	// R: 0x30, G: 0xD5, B: 0xC8 -> 48, 213, 200
	customColor := ui.NewRGBColor(48, 213, 200)
	p2.TextStyle.Bg = customColor
	p2.BorderStyle.Fg = customColor
	p2.TitleStyle.Fg = customColor

	ui.Render(p1, p2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			p1.Text = fmt.Sprintf("Resized to: %dx%d", payload.Width, payload.Height)
			ui.Clear()
			ui.Render(p1, p2)
		}
	}
}
