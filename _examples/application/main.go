package main

import (
	"log"
	"strconv"

	"github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

// MainWidget acts as a container or main logic handler.
type MainWidget struct {
	*widgets.Paragraph
	clickCount int
}

func NewMainWidget() *MainWidget {
	p := &MainWidget{
		Paragraph: widgets.NewParagraph(),
	}
	p.Title = "Application Example"
	p.Text = "Press <Space> to change text.\nClick to increment counter.\nPress <C-c> to quit."
	return p
}

func (p *MainWidget) HandleEvent(e gotui.Event) bool {
	if e.Type == gotui.KeyboardEvent {
		switch e.ID {
		case "<Space>", " ":
			p.Text = "Space Pressed! (Event Handled)"
			return true
		}
	} else if e.Type == gotui.MouseEvent {
		if e.ID == "<MouseLeft>" {
			p.clickCount++
			p.Text = "Clicked " + strconv.Itoa(p.clickCount) + " times!"
			return true
		}
	}
	return false
}

func main() {
	app := gotui.NewApp()

	w := NewMainWidget()
	// Size will be set by Run() after initialization

	app.SetRoot(w, true)

	if err := app.Run(); err != nil {
		log.Fatalf("App run failed: %v", err)
	}
}
