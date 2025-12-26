package main

import (
	"image"
	"log"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Title = "Interaction Demo"
	p.Text = "Click the buttons or checkboxes! Press q to quit."
	p.SetRect(5, 5, 40, 10)
	p.BorderStyle.Fg = ui.ColorLightCyan

	b1 := widgets.NewButton("Submit")
	b1.SetRect(5, 12, 20, 15)
	b1.Border = true
	b1.BorderStyle.Fg = ui.ColorGreen

	b2 := widgets.NewButton("Cancel")
	b2.SetRect(25, 12, 40, 15)
	b2.Border = true
	b2.ActiveStyle = ui.NewStyle(ui.ColorWhite, ui.ColorRed) // Different active color
	b2.BorderStyle.Fg = ui.ColorRed

	c1 := widgets.NewCheckbox("Enable Logging")
	c1.SetRect(5, 17, 30, 20)
	c1.Border = false // Frameless look

	c2 := widgets.NewCheckbox("Verbose Mode")
	c2.SetRect(5, 20, 30, 23)
	c2.Border = false
	c2.Checked = true
	c2.CheckedRune = '✓'

	ui.Render(p, b1, b2, c1, c2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<MouseLeft>":
			payload := e.Payload.(ui.Mouse)
			pt := image.Pt(payload.X, payload.Y)

			// Simple hit testing
			if pt.In(b1.GetRect()) {
				b1.Activate()
				ui.Render(b1)
			} else {
				if b1.IsActive {
					b1.Deactivate()
					ui.Render(b1)
				}
			}

			if pt.In(b2.GetRect()) {
				b2.Activate()
				ui.Render(b2)
			} else {
				if b2.IsActive {
					b2.Deactivate()
					ui.Render(b2)
				}
			}

			if pt.In(c1.GetRect()) {
				c1.Toggle()
				ui.Render(c1)
			}
			if pt.In(c2.GetRect()) {
				c2.Toggle()
				ui.Render(c2)
			}
		case "<MouseRelease>":
			// Deactivate buttons on release
			b1.Deactivate()
			b2.Deactivate()
			ui.Render(b1, b2)
		}
	}
}
