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

	c := widgets.NewCalendar()
	c.Title = "Calendar"
	termWidth, termHeight := ui.TerminalDimensions()
	c.SetRect(0, 0, termWidth, termHeight)

	ui.Render(c)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Right>":
			c.Month++
			if c.Month > 12 {
				c.Month = 1
				c.Year++
			}
			ui.Render(c)
		case "<Left>":
			c.Month--
			if c.Month < 1 {
				c.Month = 12
				c.Year--
			}
			ui.Render(c)
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			c.SetRect(0, 0, payload.Width, payload.Height)
			ui.Render(c)
		}
	}
}
