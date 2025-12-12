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

	// 1. Center Title, Center Text
	p1 := widgets.NewParagraph()
	p1.Title = "Centered Title"
	p1.TitleAlignment = ui.AlignCenter
	p1.Text = "This text is centered vertically and horizontally."
	p1.VerticalAlignment = ui.AlignMiddle
	p1.TextAlignment = ui.AlignCenter
	p1.SetRect(0, 0, 40, 10)

	// 2. Right Title, Bottom Text, Bottom Right Title
	p2 := widgets.NewParagraph()
	p2.Title = "Right Title"
	p2.TitleAlignment = ui.AlignRight
	p2.TitleBottom = "Bottom Right"
	p2.TitleBottomAlignment = ui.AlignRight
	p2.Text = "This text is at the bottom."
	p2.VerticalAlignment = ui.AlignBottom
	p2.SetRect(41, 0, 81, 10)

	// 3. Left Title, Top Text (Default)
	p3 := widgets.NewParagraph()
	p3.Title = "Left Title"
	p3.TitleAlignment = ui.AlignLeft
	p3.Text = "Default alignment (Top Left)."
	p3.SetRect(0, 11, 40, 21)

	termWidth, termHeight := ui.TerminalDimensions()
	gridX := termWidth / 2
	gridY := termHeight / 2 // Just split top/bottom roughly

	p1.SetRect(0, 0, gridX, gridY)
	p2.SetRect(gridX, 0, termWidth, gridY)
	p3.SetRect(0, gridY, gridX, termHeight)

	ui.Render(p1, p2, p3)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)

			termWidth = payload.Width
			termHeight = payload.Height
			gridX = termWidth / 2
			gridY = termHeight / 2

			p1.SetRect(0, 0, gridX, gridY)
			p2.SetRect(gridX, 0, termWidth, gridY)
			p3.SetRect(0, gridY, gridX, termHeight)

			ui.Clear()
			ui.Render(p1, p2, p3)
		}
	}
}
