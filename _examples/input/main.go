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

	i1 := widgets.NewInput()
	i1.Title = "Username"
	i1.Placeholder = "Enter username"
	i1.SetRect(10, 5, 50, 8)

	i2 := widgets.NewInput()
	i2.Title = "Password"
	i2.Placeholder = "Enter password"
	i2.EchoMode = widgets.EchoPassword
	i2.SetRect(10, 10, 50, 13)

	current := 0
	inputs := []*widgets.Input{i1, i2}

	ui.Render(i1, i2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			return
		case "<Tab>":
			current = (current + 1) % len(inputs)
			ui.Render(i1, i2)
		case "<Enter>":
			// Submit action?
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			// Resize logic... for this simple example we keep fixed or adjust relative?
			// Let's just center them
			midX := payload.Width / 2
			midY := payload.Height / 2
			i1.SetRect(midX-20, midY-5, midX+20, midY-2)
			i2.SetRect(midX-20, midY, midX+20, midY+3)
			ui.Clear()
			ui.Render(i1, i2)
		default:
			// Pass event to active input
			active := inputs[current]

			// Handle Input events
			switch e.ID {
			case "<Backspace>":
				active.Backspace()
			case "<Left>":
				active.MoveCursorLeft()
			case "<Right>":
				active.MoveCursorRight()
			case "<Space>":
				active.InsertRune(' ')
			default:
				if len(e.ID) == 1 {
					active.InsertRune([]rune(e.ID)[0])
				}
			}
			ui.Render(i1, i2)
		}
	}
}
