//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()

	// Vertical Scrollbar
	vScroll := widgets.NewScrollbar()
	vScroll.SetRect(termWidth-3, 0, termWidth, termHeight-3)
	vScroll.Orientation = widgets.ScrollbarVertical
	vScroll.Max = 100
	vScroll.PageSize = 20
	vScroll.ThumbStyle = ui.NewStyle(ui.ColorYellow)
	vScroll.TrackStyle = ui.NewStyle(ui.ColorWhite)
	vScroll.Border = false

	// Horizontal Scrollbar
	hScroll := widgets.NewScrollbar()
	hScroll.SetRect(0, termHeight-3, termWidth-4, termHeight)
	hScroll.Orientation = widgets.ScrollbarHorizontal
	hScroll.Max = 100
	hScroll.PageSize = 20
	hScroll.ThumbStyle = ui.NewStyle(ui.ColorGreen)
	hScroll.TrackStyle = ui.NewStyle(ui.ColorYellow) // Match example Yellow track
	hScroll.Border = false
	// Ratatui example: track="-", thumb="▮", begin="<", end=">"
	hScroll.TrackRune = '-'
	hScroll.ThumbRune = '▮'
	hScroll.BeginRune = '<'
	hScroll.EndRune = '>'

	// Information
	p := widgets.NewParagraph()
	p.Title = "Controller"
	p.Text = "Use Arrow Keys to Scroll or MouseWheel.\nVertical: Up/Down\nHorizontal: Left/Right\nq: Quit"
	// make it bigger
	p.SetRect(termWidth/2-20, termHeight/2-3, termWidth/2+20, termHeight/2+5)

	render := func() {
		p.Text = fmt.Sprintf(
			"Use Arrow Keys to Scroll or MouseWheel.\nVertical: %d/%d\nHorizontal: %d/%d\nq: Quit",
			vScroll.Current, vScroll.Max,
			hScroll.Current, hScroll.Max,
		)
		ui.Render(vScroll, hScroll, p)
	}

	render()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if handleScrollbarEvents(e, vScroll, hScroll) {
			return
		}
		render()
	}
}

func handleScrollbarEvents(e ui.Event, vScroll, hScroll *widgets.Scrollbar) bool {
	switch e.ID {
	case "q", "<C-c>":
		return true
	case "<Up>":
		vScroll.Current--
		if vScroll.Current < 0 {
			vScroll.Current = 0
		}
	case "<Down>":
		vScroll.Current++
		if vScroll.Current > vScroll.Max-vScroll.PageSize {
			vScroll.Current = vScroll.Max - vScroll.PageSize
		}
	case "<Left>":
		hScroll.Current--
		if hScroll.Current < 0 {
			hScroll.Current = 0
		}
	case "<Right>":
		hScroll.Current++
		if hScroll.Current > hScroll.Max-hScroll.PageSize {
			hScroll.Current = hScroll.Max - hScroll.PageSize
		}
	case "PageUp":
		vScroll.Current -= vScroll.PageSize
	case "PageDown":
		vScroll.Current += vScroll.PageSize
	case "<MouseWheelUp>":
		vScroll.Current--
		if vScroll.Current < 0 {
			vScroll.Current = 0
		}
	case "<MouseWheelDown>":
		vScroll.Current++
		if vScroll.Current > vScroll.Max-vScroll.PageSize {
			vScroll.Current = vScroll.Max - vScroll.PageSize
		}
	}
	return false
}
