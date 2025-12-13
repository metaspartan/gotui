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

	// 1. Logo with Gradient
	l := widgets.NewLogo()
	l.Title = "Gradient Logo"
	l.Gradient.Enabled = true
	l.Gradient.Start = ui.NewRGBColor(0, 255, 255) // Cyan
	l.Gradient.End = ui.NewRGBColor(255, 0, 255)   // Magenta

	termWidth, termHeight := ui.TerminalDimensions()
	l.SetRect(termWidth/2-25, termHeight/2-20, termWidth/2+26, termHeight/2-9)

	// 2. Gradient Paragraph
	p := widgets.NewParagraph()
	p.Title = "Gradient Paragraph"
	p.Text = "This paragraph has gradient text enabled!\nIt transitions from Green to Yellow.\nThe gradient is applied to the raw text, even with wrapping.\n\nEnjoy the colors!"
	p.Gradient.Enabled = true
	p.Gradient.Start = ui.NewRGBColor(0, 255, 0) // Green
	p.Gradient.End = ui.NewRGBColor(255, 255, 0) // Yellow
	p.SetRect(termWidth/2-25, termHeight/2-8, termWidth/2+26, termHeight/2-1)

	// 3. Gradient Gauge
	g := widgets.NewGauge()
	g.Title = "Gradient Gauge"
	g.Percent = 75
	g.Gradient.Enabled = true
	g.Gradient.Start = ui.NewRGBColor(255, 0, 0) // Red
	g.Gradient.End = ui.NewRGBColor(0, 0, 255)   // Blue
	g.SetRect(termWidth/2-25, termHeight/2+0, termWidth/2+26, termHeight/2+3)

	// 4. Gradient List (Horizontal Border)
	list := widgets.NewList()
	list.Title = "Gradient Selection List"
	list.Rows = []string{
		"[1] Option One",
		"[2] Option Two",
		"[3] Option Three",
		"[4] Option Four",
		"[5] Option Five",
	}
	list.Gradient.Enabled = true
	list.Gradient.Start = ui.NewRGBColor(255, 100, 0) // Orange
	list.Gradient.End = ui.NewRGBColor(255, 0, 100)   // Pink
	list.BorderGradient.Enabled = true
	list.BorderGradient.Direction = ui.GradientHorizontal
	list.BorderGradient.Start = ui.NewRGBColor(0, 255, 255) // Cyan
	list.BorderGradient.End = ui.NewRGBColor(0, 0, 255)     // Blue
	list.SelectedRow = 2
	list.SetRect(termWidth/2-25, termHeight/2+4, termWidth/2+26, termHeight/2+12)

	// 5. Vertical Gradient Border Block
	vb := ui.NewBlock()
	vb.Title = "Vertical Gradient Border"
	vb.BorderGradient.Enabled = true
	vb.BorderGradient.Direction = ui.GradientVertical
	vb.BorderGradient.Start = ui.NewRGBColor(255, 255, 0) // Yellow
	vb.BorderGradient.End = ui.NewRGBColor(255, 0, 0)     // Red
	vb.SetRect(termWidth/2-25, termHeight/2+13, termWidth/2+26, termHeight/2+18)

	ui.Render(l, p, g, list, vb)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
