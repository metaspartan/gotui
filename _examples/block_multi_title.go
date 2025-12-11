//go:build ignore
// +build ignore

package main

import (
	"log"

	ui "github.com/metaspartan/gotui"
	"github.com/metaspartan/gotui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// 1. Left + Right (Top)
	p1 := widgets.NewParagraph()
	p1.Title = "Main Title (Left)"
	p1.TitleRight = "Secondary (Right)"
	p1.Text = "Demonstrating Title + TitleRight"
	p1.SetRect(0, 0, 50, 10)

	// 2. Center + Right (Top) + Left + Right (Bottom)
	p2 := widgets.NewParagraph()
	p2.Title = "Centered Title"
	p2.TitleAlignment = ui.AlignCenter
	p2.TitleRight = "v1.0"
	p2.TitleBottom = "Page 1/5"
	p2.TitleBottomAlignment = ui.AlignLeft
	p2.TitleBottomRight = "Press 'q' to quit"
	p2.Text = "Demonstrating Center Top, Right Top, Left Bottom, Right Bottom"
	p2.SetRect(0, 11, 50, 21)

	// 3. Collision Test (Left + Right overlapping)
	p3 := widgets.NewParagraph()
	p3.Title = "This is a very long title that might overlap"
	p3.TitleRight = "with this also long right title"
	p3.Text = "Overlap Demo (User responsibility)"
	p3.SetRect(51, 0, 90, 10)

	ui.Render(p1, p2, p3)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
