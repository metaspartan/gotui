package main

import (
	"image"
	"log"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	// Background content
	p := widgets.NewParagraph()
	p.Title = "Main Application"
	p.Text = "Press [ENTER] to open the modal popup.\nPress [q] to quit."
	// Modal
	modal := widgets.NewModal("Could not connect to database!\n\nDo you want to retry connection?")
	modal.Title = "Network Error"
	// Use new CenterIn helper
	termWidth, termHeight := ui.TerminalDimensions()
	p.SetRect(0, 0, termWidth, termHeight)

	modal.CenterIn(0, 0, termWidth, termHeight, 50, 12)

	// Styling
	modal.Border = true
	modal.BorderRounded = true
	modal.BorderStyle.Fg = ui.ColorPink
	modal.BorderStyle.Bg = ui.ColorBlack
	modal.TextStyle.Bg = ui.ColorBlack // Text background must match modal background
	modal.TextStyle.Fg = ui.ColorWhite // Text color

	// Buttons
	_ = modal.AddButton("Retry", nil)
	_ = modal.AddButton("Cancel", nil)

	// Default to first button active
	modal.ActiveButtonIndex = 0

	showModal := false
	if ui.ScreenshotMode {
		showModal = true
	}

	draw := func() {
		var items []ui.Drawable
		items = append(items, p)
		if showModal {
			items = append(items, modal)
		}
		ui.Render(items...)
	}

	uiEvents := ui.PollEvents()
	draw()

	for {
		e := <-uiEvents

		// Handle resizing
		if e.Type == ui.ResizeEvent {
			payload := e.Payload.(ui.Resize)
			p.SetRect(0, 0, payload.Width, payload.Height)
			modal.CenterIn(0, 0, payload.Width, payload.Height, 50, 12)
			draw()
			continue
		}

		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Enter>":
			if showModal {
				// Action based on active button
				if modal.ActiveButtonIndex == 0 { // Retry
					// Simulate action, then close
					showModal = false
					draw()
				} else if modal.ActiveButtonIndex == 1 { // Cancel
					showModal = false
					draw()
				}
			} else {
				showModal = true
				modal.ActiveButtonIndex = 0 // Reset focus
				draw()
			}
		case "<Tab>":
			if showModal {
				modal.ActiveButtonIndex++
				if modal.ActiveButtonIndex >= len(modal.Buttons) {
					modal.ActiveButtonIndex = 0
				}
				draw()
			}
		case "<S-Tab>", "<Up>", "<Left>": // Also arrow support
			if showModal {
				modal.ActiveButtonIndex--
				if modal.ActiveButtonIndex < 0 {
					modal.ActiveButtonIndex = len(modal.Buttons) - 1
				}
				draw()
			}
		case "<Down>", "<Right>":
			if showModal {
				modal.ActiveButtonIndex++
				if modal.ActiveButtonIndex >= len(modal.Buttons) {
					modal.ActiveButtonIndex = 0
				}
				draw()
			}
		case "<Escape>":
			showModal = false
			draw()
		case "<MouseLeft>":
			if showModal {
				payload := e.Payload.(ui.Mouse)
				pt := image.Pt(payload.X, payload.Y)

				// Check buttons click
				clicked := false
				for i, b := range modal.Buttons {
					if pt.In(b.GetRect()) {
						modal.ActiveButtonIndex = i
						clicked = true
						draw()
						break
					}
				}

				if !clicked {
					// Click outside modal to close
					if !pt.In(modal.GetRect()) {
						showModal = false
						draw()
					}
				}
			}
		case "<MouseRelease>":
			if showModal {
				payload := e.Payload.(ui.Mouse)
				pt := image.Pt(payload.X, payload.Y)

				// Confirm action on release inside the active button
				if modal.ActiveButtonIndex >= 0 && modal.ActiveButtonIndex < len(modal.Buttons) {
					activeBtn := modal.Buttons[modal.ActiveButtonIndex]
					if pt.In(activeBtn.GetRect()) {
						// Trigger action
						showModal = false
						draw()
					}
				}
			}
		}
	}
}
