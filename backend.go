package gotui

import (
	"image"
	"io"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	Screen         tcell.Screen
	ScreenshotMode bool
)

// TTYHandle represents a custom terminal I/O source (e.g., SSH session).
// Implementations should provide Read/Write for terminal data.
type TTYHandle interface {
	io.ReadWriter
}

// InitConfig allows customizing initialization without exposing tcell types.
type InitConfig struct {
	// CustomTTY allows binding to a custom I/O source (e.g., SSH session).
	// The implementation should handle terminal protocol (ANSI/VT100).
	CustomTTY TTYHandle

	// Width and Height optionally specify dimensions for CustomTTY.
	// If zero, the library will attempt to detect them.
	Width, Height int

	// SimulationMode creates an in-memory screen for testing/screenshots.
	SimulationMode bool
	SimulationSize image.Point // e.g., image.Pt(120, 60)
}

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

// InitWithConfig initializes the library with custom configuration.
// This is useful for SSH servers, testing, or custom I/O scenarios.
// After initialization, the library must be finalized with `Close`.
func InitWithConfig(cfg *InitConfig) error {
	if cfg.SimulationMode {
		// Create in-memory screen for testing/screenshots
		Screen = tcell.NewSimulationScreen("UTF-8")
		if err := Screen.Init(); err != nil {
			return err
		}
		w, h := 120, 60
		if cfg.SimulationSize.X > 0 && cfg.SimulationSize.Y > 0 {
			w, h = cfg.SimulationSize.X, cfg.SimulationSize.Y
		}
		Screen.SetSize(w, h)
		return nil
	}

	if cfg.CustomTTY != nil {
		// Create screen from custom I/O (e.g., SSH session)
		// We need to wrap the TTYHandle into something tcell understands
		tty := &ttyAdapter{
			rw:     cfg.CustomTTY,
			width:  cfg.Width,
			height: cfg.Height,
		}

		var err error
		Screen, err = tcell.NewTerminfoScreenFromTty(tty)
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
		return nil
	}

	// Default: create standard screen
	return Init()
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

// ttyAdapter adapts a TTYHandle to tcell.Tty interface.
// This allows custom I/O sources (like SSH) without exposing tcell types.
type ttyAdapter struct {
	rw            io.ReadWriter
	width, height int
	resizeCB      func()
}

func (t *ttyAdapter) Read(p []byte) (int, error)  { return t.rw.Read(p) }
func (t *ttyAdapter) Write(p []byte) (int, error) { return t.rw.Write(p) }
func (t *ttyAdapter) Close() error {
	if c, ok := t.rw.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (t *ttyAdapter) Start() error { return nil }
func (t *ttyAdapter) Stop() error  { return nil }
func (t *ttyAdapter) Drain() error { return nil }

func (t *ttyAdapter) NotifyResize(cb func()) {
	t.resizeCB = cb
	// If the underlying TTYHandle supports resize notifications, hook them up
	type resizable interface {
		NotifyResize(func())
	}
	if r, ok := t.rw.(resizable); ok {
		r.NotifyResize(cb)
	}
}

func (t *ttyAdapter) WindowSize() (tcell.WindowSize, error) {
	// Try to detect dimensions from the underlying TTYHandle
	type windowSizer interface {
		WindowSize() (tcell.WindowSize, error)
	}
	if ws, ok := t.rw.(windowSizer); ok {
		return ws.WindowSize()
	}

	// Fall back to configured dimensions
	w, h := t.width, t.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	return tcell.WindowSize{Width: w, Height: h}, nil
}
