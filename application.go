package gotui

import (
	"sync"
)

// Application represents the application.
type Application struct {
	root    Widget
	focus   Widget
	running bool
	stop    chan struct{}
	sync.Mutex
}

// NewApp returns a new Application.
func NewApp() *Application {
	return &Application{}
}

// SetRoot sets the root widget of the application.
// If focus is true, the root widget is also focused.
func (a *Application) SetRoot(root Widget, focus bool) {
	a.Lock()
	defer a.Unlock()
	a.root = root
	if focus {
		a.focus = root
	}
}

// SetFocus sets the focus to the given widget.
func (a *Application) SetFocus(p Widget) {
	a.Lock()
	defer a.Unlock()
	a.focus = p
}

// Stop stops the application.
func (a *Application) Stop() {
	a.Lock()
	defer a.Unlock()
	if a.running && a.stop != nil {
		close(a.stop)
		a.running = false
	}
}

// getRoot returns the root widget under lock.
func (a *Application) getRoot() Widget {
	a.Lock()
	defer a.Unlock()
	return a.root
}

// Run runs the application.
func (a *Application) Run() error {
	if err := Init(); err != nil {
		return err
	}
	defer Close()

	a.Lock()
	a.running = true
	a.stop = make(chan struct{}) // Recreate for each Run
	a.Unlock()

	// Size the root widget to terminal size after init
	root := a.getRoot()
	if root != nil {
		w, h := TerminalDimensions()
		root.SetRect(0, 0, w, h)
		Render(root)
	}

	uiEvents := PollEvents()
	for {
		select {
		case <-a.stop:
			return nil
		case e := <-uiEvents:
			if a.handleEvent(e) {
				return nil
			}
		}
	}
}

// handleEvent processes a single event. Returns true if the application should stop.
func (a *Application) handleEvent(e Event) bool {
	if e.Type == ResizeEvent {
		a.handleResize(e)
		return false
	}

	handled := a.dispatchKeyOrMouse(e)

	// Default handlers (like Quit)
	if !handled {
		if e.ID == "<C-c>" || e.ID == "q" {
			return true
		}
	}

	// Re-render
	root := a.getRoot()
	if root != nil {
		Render(root)
	}
	return false
}

func (a *Application) handleResize(e Event) {
	payload := e.Payload.(Resize)
	root := a.getRoot()
	if root != nil {
		root.SetRect(0, 0, payload.Width, payload.Height)
		Render(root)
	}
}

func (a *Application) dispatchKeyOrMouse(e Event) bool {
	handled := false
	a.Lock()
	focus := a.focus
	root := a.root
	a.Unlock()

	// 1. Dispatch to Focus (Keyboard)
	if e.Type == KeyboardEvent && focus != nil {
		if focus.HandleEvent(e) {
			handled = true
		}
	}

	// 2. Dispatch to Root (Mouse? Bubble up?)
	if !handled && root != nil {
		if root.HandleEvent(e) {
			handled = true
		}
	}
	return handled
}
