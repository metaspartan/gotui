package gotui

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	Screen         tcell.Screen
	ScreenshotMode bool
)

// Init initializes tcell and is required to render anything.
// After initialization, the library must be finalized with `Close`.
func Init() error {
	// Check for -screenshot flag automatically
	for i, arg := range os.Args {
		if arg == "-screenshot" {
			ScreenshotMode = true
			// Remove flag so app logic doesn't see it
			os.Args = append(os.Args[:i], os.Args[i+1:]...)

			// Initialize a simulation screen so tcell's color palette works correctly
			Screen = tcell.NewSimulationScreen("UTF-8")
			if err := Screen.Init(); err != nil {
				return err
			}
			Screen.SetSize(120, 60)
			return nil
		}
	}

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
	if ScreenshotMode {
		return 120, 60
	}
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
