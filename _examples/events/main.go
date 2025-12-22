package main

import (
	"log"
	"strconv"

	"github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

// InteractiveParagraph is a custom widget that handles clicks.
type InteractiveParagraph struct {
	*widgets.Paragraph
	clickCount int
}

func NewInteractiveParagraph() *InteractiveParagraph {
	p := &InteractiveParagraph{
		Paragraph: widgets.NewParagraph(),
	}
	p.Title = "Click Me!"
	p.Text = "No clicks yet."
	return p
}

func (p *InteractiveParagraph) HandleEvent(e gotui.Event) bool {
	if e.Type == gotui.MouseEvent {
		payload := e.Payload.(gotui.Mouse)
		if payload.X >= p.Inner.Min.X && payload.X < p.Inner.Max.X &&
			payload.Y >= p.Inner.Min.Y && payload.Y < p.Inner.Max.Y {
			// Handle Click inside widget
			if e.ID == "<MouseLeft>" {
				p.clickCount++
				p.Text = "Clicked " + strconv.Itoa(p.clickCount) + " times!"
				return true
			}
		}
	} else if e.Type == gotui.KeyboardEvent {
		// Handle specific keys if focused (simple version)
		if e.ID == "<Space>" || e.ID == " " {
			p.Text = "Space Pressed!"
			return true
		}
	}
	return false
}

func main() {
	if err := gotui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer gotui.Close()

	p1 := NewInteractiveParagraph()
	p1.SetRect(0, 0, 40, 5)

	p2 := widgets.NewParagraph()
	p2.Title = "Status"
	p2.Text = "Waiting for events..."
	p2.SetRect(0, 6, 40, 10)

	gotui.Render(p1, p2)

	uiEvents := gotui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}

		// Event Loop Pattern
		handled := false
		if p1.HandleEvent(e) {
			p2.Text = "Event Handled by Paragraph!"
			handled = true
		}

		if handled {
			gotui.Render(p1, p2)
		}
	}
}
