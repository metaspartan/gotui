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

	// Helper to create a merged block
	createBlock := func(title, text string, x1, y1, x2, y2 int) *widgets.Paragraph {
		p := widgets.NewParagraph()
		p.Title = title
		p.Text = text
		p.SetRect(x1, y1, x2, y2)
		p.BorderCollapse = true // Enable merging
		return p
	}

	// Layout with overlapping borders
	// Top: 0-10
	// Bottom: 10-20 (Overlap at 10? No, 10 is exclusive in Rect implies 0..9.
	// If Rect is (0,0, 10,10), pixels are 0..9.
	// If next is (0,9, 10,19), overlap at y=9.

	// Let's create a 2x2 grid overlapping.
	// TopLeft: 0,0 - 20,10
	// TopRight: 19,0 - 40,10 (Overlap x=19)
	// BotLeft: 0,9 - 20,20 (Overlap y=9)
	// BotRight: 19,9 - 40,20 (Overlap x=19, y=9)

	p1 := createBlock("Top Left", "This block shares borders.\nRight and Bottom are merged.", 0, 0, 20, 10)
	p2 := createBlock("Top Right", "Merged Left and Bottom.\nCorner is a T-junction.", 19, 0, 40, 10)
	p3 := createBlock("Bottom Left", "Merged Top and Right.", 0, 9, 20, 20)
	p4 := createBlock("Bottom Right", "Merged Top and Left.\nCenter is a Cross.", 19, 9, 40, 20)

	// A center block overlapping all 4?
	// Let's stick to 2x2 for clean T and Cross testing.

	// Rendering order matters.
	// p1 drawn first.
	// p2 draws and merges left (19,y).
	// p3 draws and merges top (x,9).
	// p4 draws and merges top (x,9) and left (19,y). AND the corner (19,9) should resolve to CROSS.

	uiEvents := ui.PollEvents()
	for {
		// Redraw on every loop
		ui.Render(p1, p2, p3, p4)

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			_ = payload
		}
	}
}
