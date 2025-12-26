package main

import (
	"image"
	"log"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

type App struct {
	p         *widgets.Paragraph
	modal     *widgets.Modal
	showModal bool
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	app := NewApp()
	app.Draw()

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if !app.HandleEvent(e) {
			return
		}
	}
}

func NewApp() *App {
	p := widgets.NewParagraph()
	p.Title = "Main Application"
	p.Text = "Press [ENTER] to open the modal popup.\nPress [q] to quit."

	modal := widgets.NewModal("Could not connect to database!\n\nDo you want to retry connection?")
	modal.Title = "Network Error"
	modal.Border = true
	modal.BorderRounded = true
	modal.BorderStyle.Fg = ui.ColorPink
	modal.BorderStyle.Bg = ui.ColorBlack
	modal.TextStyle.Bg = ui.ColorBlack
	modal.TextStyle.Fg = ui.ColorWhite
	modal.ActiveButtonIndex = 0

	app := &App{
		p:     p,
		modal: modal,
	}

	if ui.ScreenshotMode {
		app.showModal = true
	}

	termWidth, termHeight := ui.TerminalDimensions()
	app.resize(termWidth, termHeight)

	// Buttons
	_ = modal.AddButton("Retry", func() {
		app.showModal = false
	})
	_ = modal.AddButton("Cancel", func() {
		app.showModal = false
	})

	return app
}

func (a *App) resize(w, h int) {
	a.p.SetRect(0, 0, w, h)
	a.modal.CenterIn(0, 0, w, h, 50, 12)
}

func (a *App) Draw() {
	var items []ui.Drawable
	items = append(items, a.p)
	if a.showModal {
		items = append(items, a.modal)
	}
	ui.Render(items...)
}

func (a *App) HandleEvent(e ui.Event) bool {
	if e.Type == ui.ResizeEvent {
		payload := e.Payload.(ui.Resize)
		a.resize(payload.Width, payload.Height)
		a.Draw()
		return true
	}

	switch e.ID {
	case "q", "<C-c>":
		return false
	case "<Escape>":
		a.showModal = false
		a.Draw()
	case "<Enter>":
		a.handleEnter()
	case "<Tab>":
		a.cycleButton(1)
	case "<S-Tab>", "<Up>", "<Left>":
		a.cycleButton(-1)
	case "<Down>", "<Right>":
		a.cycleButton(1)
	case "<MouseLeft>":
		a.handleMouseLeft(e)
	case "<MouseRelease>":
		a.handleMouseRelease(e)
	}
	return true
}

func (a *App) handleEnter() {
	if a.showModal {
		if a.modal.ActiveButtonIndex >= 0 && a.modal.ActiveButtonIndex < len(a.modal.Buttons) {
			btn := a.modal.Buttons[a.modal.ActiveButtonIndex]
			if btn.OnClick != nil {
				btn.OnClick()
			}
			a.Draw()
		}
	} else {
		a.showModal = true
		a.modal.ActiveButtonIndex = 0
		a.Draw()
	}
}

func (a *App) cycleButton(dir int) {
	if !a.showModal {
		return
	}
	a.modal.ActiveButtonIndex += dir
	if a.modal.ActiveButtonIndex >= len(a.modal.Buttons) {
		a.modal.ActiveButtonIndex = 0
	} else if a.modal.ActiveButtonIndex < 0 {
		a.modal.ActiveButtonIndex = len(a.modal.Buttons) - 1
	}
	a.Draw()
}

func (a *App) handleMouseLeft(e ui.Event) {
	if !a.showModal {
		return
	}
	payload := e.Payload.(ui.Mouse)
	pt := image.Pt(payload.X, payload.Y)

	clicked := false
	for i, b := range a.modal.Buttons {
		if pt.In(b.GetRect()) {
			a.modal.ActiveButtonIndex = i
			clicked = true
			a.Draw()
			break
		}
	}

	if !clicked {
		if !pt.In(a.modal.GetRect()) {
			a.showModal = false
			a.Draw()
		}
	}
}

func (a *App) handleMouseRelease(e ui.Event) {
	if !a.showModal {
		return
	}
	payload := e.Payload.(ui.Mouse)
	pt := image.Pt(payload.X, payload.Y)

	if a.modal.ActiveButtonIndex >= 0 && a.modal.ActiveButtonIndex < len(a.modal.Buttons) {
		activeBtn := a.modal.Buttons[a.modal.ActiveButtonIndex]
		if pt.In(activeBtn.GetRect()) {
			if activeBtn.OnClick != nil {
				activeBtn.OnClick()
			}
			a.Draw()
		}
	}
}
