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

	// 1. Basic Block
	b1 := widgets.NewParagraph()
	b1.Text = "Basic Block\n(Default Style)"
	b1.SetRect(0, 0, 25, 5)

	// 2. Title Styling
	b2 := widgets.NewParagraph()
	b2.Text = "Styled Block\n(Red Border, Yellow Title)"
	b2.SetRect(26, 0, 51, 5)
	b2.Title = "Custom Style"
	b2.TitleStyle = ui.NewStyle(ui.ColorYellow)
	b2.BorderStyle = ui.NewStyle(ui.ColorRed)

	// 3. Title Alignments
	b3 := widgets.NewParagraph()
	b3.Text = "Left Aligned Title"
	b3.SetRect(0, 6, 25, 11)
	b3.Title = "Left"
	b3.TitleAlignment = ui.AlignLeft

	b4 := widgets.NewParagraph()
	b4.Text = "Center Aligned Title"
	b4.SetRect(26, 6, 51, 11)
	b4.Title = "Center"
	b4.TitleAlignment = ui.AlignCenter

	b5 := widgets.NewParagraph()
	b5.Text = "Right Aligned Title"
	b5.SetRect(52, 6, 77, 11)
	b5.Title = "Right"
	b5.TitleAlignment = ui.AlignRight

	// 4. Bottom Titles
	b6 := widgets.NewParagraph()
	b6.Text = "Bottom Title (Center)"
	b6.SetRect(0, 12, 25, 17)
	b6.TitleBottom = "Bottom Info"
	b6.TitleBottomAlignment = ui.AlignCenter

	b7 := widgets.NewParagraph()
	b7.Text = "Top & Bottom Titles"
	b7.SetRect(26, 12, 51, 17)
	b7.Title = "Top"
	b7.TitleAlignment = ui.AlignLeft
	b7.TitleBottom = "Bottom"
	b7.TitleBottomAlignment = ui.AlignRight

	// 5. Padding
	b8 := widgets.NewParagraph()
	b8.Text = "Rounded Block\n(2px Padding)"
	b8.SetRect(52, 12, 77, 17)
	b8.BorderRounded = true
	b8.PaddingLeft = 2
	b8.PaddingTop = 1
	b8.PaddingRight = 2
	b8.PaddingBottom = 1

	// Instructions
	info := widgets.NewParagraph()
	info.Text = "Press q to quit"
	info.SetRect(0, 18, 50, 21)
	info.Border = false

	ui.Render(b1, b2, b3, b4, b5, b6, b7, b8, info)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
