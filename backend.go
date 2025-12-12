package gotui

import (
	"io"
	"os"

	"github.com/gdamore/tcell/v2"
)

func NewBackend(cfg *InitConfig) (*Backend, error) {
	b := &Backend{}
	if err := b.InitWithConfig(cfg); err != nil {
		return nil, err
	}
	return b, nil
}

func Init() error {
	if err := DefaultBackend.Init(); err != nil {
		return err
	}
	Screen = DefaultBackend.Screen
	ScreenshotMode = DefaultBackend.ScreenshotMode
	return nil
}

func InitWithConfig(cfg *InitConfig) error {
	if err := DefaultBackend.InitWithConfig(cfg); err != nil {
		return err
	}
	Screen = DefaultBackend.Screen
	ScreenshotMode = DefaultBackend.ScreenshotMode
	return nil
}

func Close() {
	DefaultBackend.Close()
}

func TerminalDimensions() (int, int) {
	return DefaultBackend.TerminalDimensions()
}

func Clear() {
	DefaultBackend.Clear()
}

func ClearBackground(c Color) {
	DefaultBackend.ClearBackground(c)
}

func (b *Backend) Init() error {
	for i, arg := range os.Args {
		if arg == "-screenshot" {
			b.ScreenshotMode = true
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

func (b *Backend) InitWithConfig(cfg *InitConfig) error {
	if cfg == nil {
		return b.Init()
	}
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

func (b *Backend) Close() {
	if b.Screen != nil {
		b.Screen.Fini()
	}
}

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

func (b *Backend) Clear() {
	if b.Screen != nil {
		b.Screen.Clear()
	}
}

func (b *Backend) ClearBackground(c Color) {
	if b.Screen != nil {
		b.Screen.SetStyle(tcell.StyleDefault.Background(c))
		b.Screen.Clear()
	}
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
