package main

import (
	"time"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// Create spinners for each style
	spinners := []*widgets.Spinner{}

	descriptions := []struct {
		Title  string
		Frames []string
	}{
		{"Line", widgets.SpinnerLine},
		{"Dots", widgets.SpinnerDots},
		{"MiniDots", widgets.SpinnerMiniDots},
		{"Pulse", widgets.SpinnerPulse},
		{"Points", widgets.SpinnerPoints},
		// {"Globe", widgets.SpinnerGlobe},
		// {"Moon", widgets.SpinnerMoon},
		// {"Clock", widgets.SpinnerClock},
		// {"Monkey", widgets.SpinnerMonkey},
		// {"Star", widgets.SpinnerStar},
		// {"Hamburger", widgets.SpinnerHamburger},
		{"Grow Vert", widgets.SpinnerGrowVertical},
		{"Grow Horiz", widgets.SpinnerGrowHorizontal},
		{"Arrow", widgets.SpinnerArrow},
		{"Triangle", widgets.SpinnerTriangle},
		{"Halves", widgets.SpinnerCircleHalves},
		{"Ball", widgets.SpinnerBouncingBall},
	}

	termWidth, termHeight := ui.TerminalDimensions()
	gridWidth := 20
	gridHeight := 4
	cols := termWidth / gridWidth

	for i, desc := range descriptions {
		s := widgets.NewSpinner()
		s.Title = desc.Title
		s.Frames = desc.Frames
		s.Label = "loading..."
		s.FormatString = "%s %s" // icon space label

		// Grid layout
		col := i % cols
		row := i / cols
		x1 := col * gridWidth
		y1 := row * gridHeight
		x2 := x1 + gridWidth
		y2 := y1 + gridHeight

		s.SetRect(x1, y1, x2, y2)
		s.BorderStyle.Fg = ui.ColorGreen

		spinners = append(spinners, s)
	}

	// Instructions
	p := widgets.NewParagraph()
	p.Text = "Press [q] to quit"
	p.SetRect(0, termHeight-3, termWidth, termHeight)
	p.Border = false

	render := func() {
		var items []ui.Drawable
		for _, s := range spinners {
			items = append(items, s)
		}
		items = append(items, p)
		ui.Render(items...)
	}

	render()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-ticker.C:
			for _, s := range spinners {
				s.Advance()
			}
			render()
		case e := <-uiEvents:
			if e.Type == ui.ResizeEvent {
				payload := e.Payload.(ui.Resize)
				termWidth, termHeight = payload.Width, payload.Height
				cols = termWidth / gridWidth

				for i, s := range spinners {
					col := i % cols
					row := i / cols
					x1 := col * gridWidth
					y1 := row * gridHeight
					x2 := x1 + gridWidth
					y2 := y1 + gridHeight
					s.SetRect(x1, y1, x2, y2)
				}
				p.SetRect(0, termHeight-3, termWidth, termHeight)
				render()
			}
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	}
}
