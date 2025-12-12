package gotui

import (
	"image"
	"io"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	// DefaultBackend is the default backend instance for global function compatibility.
	DefaultBackend = &Backend{}
)

// Helper to expose Screen for legacy code accessing it directly.
// Deprecated: usage of Screen variable is discouraged. Use Backend instances.
var Screen tcell.Screen

// Helper to expose ScreenshotMode for legacy code.
var ScreenshotMode bool

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

// Backend encapsulates the tcell screen and state, allowing multiple instances.
type Backend struct {
	Screen         tcell.Screen
	ScreenshotMode bool
}

// NewBackend creates a new Backend with the provided config.
func NewBackend(cfg *InitConfig) (*Backend, error) {
	b := &Backend{}
	if err := b.InitWithConfig(cfg); err != nil {
		return nil, err
	}
	return b, nil
}

// Init initializes the default backend.
func Init() error {
	if err := DefaultBackend.Init(); err != nil {
		return err
	}
	Screen = DefaultBackend.Screen
	ScreenshotMode = DefaultBackend.ScreenshotMode
	return nil
}

// InitWithConfig initializes the default backend with custom config.
func InitWithConfig(cfg *InitConfig) error {
	if err := DefaultBackend.InitWithConfig(cfg); err != nil {
		return err
	}
	Screen = DefaultBackend.Screen
	ScreenshotMode = DefaultBackend.ScreenshotMode
	return nil
}

// Close closes the default backend.
func Close() {
	DefaultBackend.Close()
}

// TerminalDimensions returns the dimensions of the default backend.
func TerminalDimensions() (int, int) {
	return DefaultBackend.TerminalDimensions()
}

// Clear clears the default backend.
func Clear() {
	DefaultBackend.Clear()
}

// ClearBackground sets background on default backend.
func ClearBackground(c Color) {
	DefaultBackend.ClearBackground(c)
}

// Init initializes the backend's tcell screen.
func (b *Backend) Init() error {
	// Check for -screenshot flag automatically (only for default backend usually, but check anyway)
	for i, arg := range os.Args {
		if arg == "-screenshot" {
			b.ScreenshotMode = true
			// Remove flag so app logic doesn't see it
			os.Args = append(os.Args[:i], os.Args[i+1:]...)

			b.Screen = tcell.NewSimulationScreen("UTF-8")
			if err := b.Screen.Init(); err != nil {
				return err
			}
			b.Screen.SetSize(120, 60)
			return nil
		}
	}

	var err error
	b.Screen, err = tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := b.Screen.Init(); err != nil {
		return err
	}
	b.Screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorDefault))
	b.Screen.EnableMouse()
	return nil
}

// InitWithConfig initializes the backend with custom configuration.
func (b *Backend) InitWithConfig(cfg *InitConfig) error {
	if cfg.SimulationMode {
		b.Screen = tcell.NewSimulationScreen("UTF-8")
		if err := b.Screen.Init(); err != nil {
			return err
		}
		w, h := 120, 60
		if cfg.SimulationSize.X > 0 && cfg.SimulationSize.Y > 0 {
			w, h = cfg.SimulationSize.X, cfg.SimulationSize.Y
		}
		b.Screen.SetSize(w, h)
		return nil
	}

	if cfg.CustomTTY != nil {
		tty := &ttyAdapter{
			rw:     cfg.CustomTTY,
			width:  cfg.Width,
			height: cfg.Height,
		}

		var err error
		b.Screen, err = tcell.NewTerminfoScreenFromTty(tty)
		if err != nil {
			return err
		}
		if err := b.Screen.Init(); err != nil {
			return err
		}
		b.Screen.SetStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDefault))
		b.Screen.EnableMouse()
		return nil
	}

	return b.Init()
}

// Close closes the backend.
func (b *Backend) Close() {
	if b.Screen != nil {
		b.Screen.Fini()
	}
}

// TerminalDimensions returns the dimensions of the screen.
func (b *Backend) TerminalDimensions() (int, int) {
	if b.ScreenshotMode {
		return 120, 60
	}
	if b.Screen == nil {
		return 0, 0
	}
	width, height := b.Screen.Size()
	return width, height
}

// Clear clears the screen.
func (b *Backend) Clear() {
	if b.Screen != nil {
		b.Screen.Clear()
	}
}

// ClearBackground sets the default background color and clears.
func (b *Backend) ClearBackground(c Color) {
	if b.Screen != nil {
		b.Screen.SetStyle(tcell.StyleDefault.Background(c))
		b.Screen.Clear()
	}
}

// ttyAdapter adapts a TTYHandle to tcell.Tty interface.
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
	type resizable interface {
		NotifyResize(func())
	}
	if r, ok := t.rw.(resizable); ok {
		r.NotifyResize(cb)
	}
}

func (t *ttyAdapter) WindowSize() (tcell.WindowSize, error) {
	type tcellWindowSizer interface {
		WindowSize() (tcell.WindowSize, error)
	}
	type simpleWindowSizer interface {
		WindowSize() (int, int, error)
	}

	if ws, ok := t.rw.(tcellWindowSizer); ok {
		return ws.WindowSize()
	}
	if ws, ok := t.rw.(simpleWindowSizer); ok {
		w, h, err := ws.WindowSize()
		return tcell.WindowSize{Width: w, Height: h}, err
	}

	w, h := t.width, t.height
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	return tcell.WindowSize{Width: w, Height: h}, nil
}
