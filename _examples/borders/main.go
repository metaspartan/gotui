package main

import (
	"log"

	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := gotui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer gotui.Close()

	p1 := widgets.NewParagraph()
	p1.Text = "Default Border"
	p1.SetRect(0, 0, 25, 5)

	p2 := widgets.NewParagraph()
	p2.Text = "Double Border"
	double := gotui.BorderSetDouble()
	p2.Block.BorderSet = &double
	p2.SetRect(30, 0, 55, 5)

	p3 := widgets.NewParagraph()
	p3.Text = "Thick Border"
	thick := gotui.BorderSetThick()
	p3.Block.BorderSet = &thick
	p3.SetRect(0, 10, 25, 15)

	p4 := widgets.NewParagraph()
	p4.Text = "Round Border"
	round := gotui.BorderSetRound()
	p4.Block.BorderSet = &round
	p4.SetRect(30, 10, 55, 15)

	p5 := widgets.NewParagraph()
	p5.Text = "Custom Border"
	custom := gotui.BorderSet{
		Top:         '-',
		Bottom:      '-',
		Left:        '|',
		Right:       '|',
		TopLeft:     '*',
		TopRight:    '*',
		BottomLeft:  '*',
		BottomRight: '*',
		TopT:        'v',
		BottomT:     '^',
		LeftT:       '>',
		RightT:      '<',
	}
	p5.Block.BorderSet = &custom
	p5.SetRect(0, 20, 25, 25)

	p6 := widgets.NewParagraph()
	p6.Text = "Hidden Border"
	hidden := gotui.BorderSetHidden()
	p6.Block.BorderSet = &hidden
	p6.SetRect(30, 20, 55, 25)

	gotui.Render(p1, p2, p3, p4, p5, p6)

	uiEvents := gotui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
