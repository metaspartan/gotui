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

	ta := widgets.NewTextArea()
	ta.Title = "TextArea (Press 'q' or 'Esc' to quit)"
	ta.Text = "Welcome to Gotui TextArea!\n\nType something here...\nSupports multiline and navigation."
	ta.SetRect(0, 0, 50, 20)

	termWidth, termHeight := ui.TerminalDimensions()
	ta.SetRect(0, 0, termWidth, termHeight)

	ui.Render(ta)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if handleTextAreaEvents(ta, e) {
			return
		}
	}
}

func handleTextAreaEvents(ta *widgets.TextArea, e ui.Event) bool {
	switch e.ID {
	case "q", "<C-c>", "<Escape>":
		return true
	case "<Resize>":
		payload := e.Payload.(ui.Resize)
		ta.SetRect(0, 0, payload.Width, payload.Height)
		ui.Clear()
		ui.Render(ta)
	case "<MouseWheelUp>":
		ta.MoveCursor(0, -1)
		ui.Render(ta)
	case "<MouseWheelDown>":
		ta.MoveCursor(0, 1)
		ui.Render(ta)
	case "<Up>":
		ta.MoveCursor(0, -1)
		ui.Render(ta)
	case "<Down>":
		ta.MoveCursor(0, 1)
		ui.Render(ta)
	case "<Left>":
		ta.MoveCursor(-1, 0)
		ui.Render(ta)
	case "<Right>":
		ta.MoveCursor(1, 0)
		ui.Render(ta)
	case "<Enter>":
		ta.InsertNewline()
		ui.Render(ta)
	case "<Backspace>", "<Delete>":
		// Basic Backspace support
		ta.DeleteRune()
		ui.Render(ta)
	case "<Space>":
		ta.InsertRune(' ')
		ui.Render(ta)
	case "<Tab>":
		ta.InsertRune('\t')
		ui.Render(ta)
	default:
		// Handle simple character input
		if len(e.ID) == 1 {
			r := []rune(e.ID)[0]
			// TODO: Better check for printable characters
			ta.InsertRune(r)
			ui.Render(ta)
		}
	}
	return false
}
