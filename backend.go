package gotui

import (
	"github.com/gdamore/tcell/v2"
)

var Screen tcell.Screen

// Init initializes tcell and is required to render anything.
// After initialization, the library must be finalized with `Close`.
func Init() error {
	var err error
	Screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := Screen.Init(); err != nil {
		return err
	}
	Screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorDefault))
	Screen.EnableMouse()
	// Output mode is handled automatically by tcell usually (24-bit if supported)
	return nil
}

// Close closes tcell.
func Close() {
	if Screen != nil {
		Screen.Fini()
	}
}

func TerminalDimensions() (int, int) {
	if Screen == nil {
		return 0, 0
	}
	width, height := Screen.Size()
	return width, height
}

func Clear() {
	if Screen != nil {
		Screen.Clear()
	}
}

// ClearBackground sets the default background color and clears the screen.
func ClearBackground(c Color) {
	if Screen != nil {
		Screen.SetStyle(tcell.StyleDefault.Background(c))
		Screen.Clear()
	}
}
