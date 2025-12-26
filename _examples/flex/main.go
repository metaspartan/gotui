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

	// 1. Header (Fixed height 3)
	p1 := widgets.NewParagraph()
	p1.Title = "Header"
	p1.Text = "Fixed Height: 3 rows\nStandard Layout Demo"
	p1.Border = true

	// 2. Sidebar (Fixed Width 20)
	p2 := widgets.NewParagraph()
	p2.Title = "Sidebar"
	p2.Text = "Fixed Width: 20\n\n- Item 1\n- Item 2\n- Item 3"
	p2.Border = true

	// 3. Main Content (Proportion 1 aka 100% of rest)
	p3 := widgets.NewParagraph()
	p3.Title = "Main Content"
	p3.Text = "This block takes up all remaining space.\nResize the window to see it adapt!"
	p3.Border = true

	// 4. Footer (Fixed height 1)
	p4 := widgets.NewParagraph()
	p4.Text = "Footer: Fixed Height 1"
	p4.Border = false
	p4.TextStyle.Bg = ui.ColorBlue

	// Compose logic:
	// Root is Vertical: Header, Middle, Footer
	// Middle is Horizontal: Sidebar, Content

	middleFlex := widgets.NewFlex()
	middleFlex.Direction = widgets.FlexRow // Horizontal layout
	middleFlex.AddItem(p2, 20, 0, false)   // Fixed 20 width
	middleFlex.AddItem(p3, 0, 1, false)    // 100% remaining

	rootFlex := widgets.NewFlex()
	rootFlex.Direction = widgets.FlexColumn   // Vertical layout
	rootFlex.AddItem(p1, 3, 0, false)         // Fixed 3 height
	rootFlex.AddItem(middleFlex, 0, 1, false) // 100% remaining height
	rootFlex.AddItem(p4, 1, 0, false)         // Fixed 1 height

	// Set root size to screen
	w, h := ui.TerminalDimensions()
	rootFlex.SetRect(0, 0, w, h)

	ui.Render(rootFlex)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			rootFlex.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(rootFlex)
		}
	}
}
