package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

var oneSession sync.Mutex // gotui uses a global ui.Screen -> serialize sessions for PoC

// ---- SSH -> tcell.Tty adapter ----

type sessionTTY struct {
	sess ssh.Session

	mu       sync.RWMutex
	w, h     int
	resizeCb func()

	winCh  <-chan ssh.Window
	closed chan struct{}
}

func newSessionTTY(sess ssh.Session) (*sessionTTY, error) {
	pty, winCh, ok := sess.Pty()
	if !ok {
		return nil, fmt.Errorf("no PTY requested (try: ssh -tt host -p 2222)")
	}

	t := &sessionTTY{
		sess:   sess,
		w:      pty.Window.Width,
		h:      pty.Window.Height,
		winCh:  winCh,
		closed: make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-t.closed:
				return
			case win, ok := <-t.winCh:
				if !ok {
					return
				}
				t.mu.Lock()
				t.w, t.h = win.Width, win.Height
				cb := t.resizeCb
				t.mu.Unlock()
				if cb != nil {
					cb()
				}
			}
		}
	}()

	return t, nil
}

func (t *sessionTTY) Size() (int, int) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.w, t.h
}

// tcell.Tty interface
func (t *sessionTTY) Start() error { return nil }
func (t *sessionTTY) Stop() error  { return nil }
func (t *sessionTTY) Drain() error { return nil }

func (t *sessionTTY) NotifyResize(cb func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.resizeCb = cb
}

func (t *sessionTTY) WindowSize() (tcell.WindowSize, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return tcell.WindowSize{Width: t.w, Height: t.h}, nil
}

func (t *sessionTTY) Read(p []byte) (int, error)  { return t.sess.Read(p) }
func (t *sessionTTY) Write(p []byte) (int, error) { return t.sess.Write(p) }

func (t *sessionTTY) Close() error {
	select {
	case <-t.closed:
	default:
		close(t.closed)
	}
	return t.sess.Close()
}

// ---- Dashboard construction ----

type dashboard struct {
	grid *ui.Grid

	// widgets we mutate over time / events
	sl, sl2   *widgets.Sparkline
	lg1, lg2  *widgets.LineGauge
	bc        *widgets.BarChart
	g         *widgets.Gauge
	plot      *widgets.Plot
	plotData  [][]float64
	logs      *widgets.List
	tickCount int
}

func newDashboard() *dashboard {
	d := &dashboard{}

	// Header
	p := widgets.NewParagraph()
	p.Title = "gotui Dashboard"
	p.Text = "PRESS q TO QUIT | Grid Layout Demo"
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan
	p.TitleStyle = ui.NewStyle(ui.ColorCyan, ui.ColorClear, ui.ModifierBold)
	p.TitleAlignment = ui.AlignCenter
	p.TitleRight = "ssh"
	p.BorderRounded = false

	// Sparklines
	slData := make([]float64, 200)
	d.sl = widgets.NewSparkline()
	d.sl.Data = slData
	d.sl.LineColor = ui.ColorGreen
	d.sl.TitleStyle.Fg = ui.ColorWhite
	d.sl.MaxVal = 100

	d.sl2 = widgets.NewSparkline()
	d.sl2.Data = slData
	d.sl2.LineColor = ui.ColorMagenta
	d.sl2.TitleStyle.Fg = ui.ColorWhite
	d.sl2.MaxVal = 100

	slg := widgets.NewSparklineGroup(d.sl, d.sl2)
	slg.Title = "CPU Usage"
	slg.TitleStyle.Fg = ui.ColorGreen
	slg.BorderStyle.Fg = ui.ColorGreen
	slg.TitleRight = "Core 0 & 1"
	slg.BorderRounded = true

	// Line gauges
	d.lg1 = widgets.NewLineGauge()
	d.lg1.Title = "Memory"
	d.lg1.Percent = 45
	d.lg1.BarRune = '■'
	d.lg1.LineColor = ui.ColorYellow
	d.lg1.TitleStyle.Fg = ui.ColorYellow
	d.lg1.BorderRounded = true

	d.lg2 = widgets.NewLineGauge()
	d.lg2.Title = "Load"
	d.lg2.Percent = 60
	d.lg2.BarRune = '▰'
	d.lg2.BarRuneEmpty = '▱'
	d.lg2.LineColor = ui.ColorRed
	d.lg2.TitleStyle.Fg = ui.ColorRed
	d.lg2.BorderRounded = true

	// Bar chart
	d.bc = widgets.NewBarChart()
	d.bc.Title = "Network Traffic"
	d.bc.TitleBottom = "MB/s"
	d.bc.TitleBottomAlignment = ui.AlignRight
	d.bc.Labels = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	d.bc.BarColors = []ui.Color{ui.ColorBlue, ui.ColorCyan}
	d.bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	d.bc.Data = []float64{3, 2, 5, 3, 9, 5}
	d.bc.TitleStyle.Fg = ui.ColorBlue
	d.bc.BorderStyle.Fg = ui.ColorBlue
	d.bc.BorderRounded = true
	d.bc.BarWidth = 0
	d.bc.BarGap = 1
	d.bc.MaxVal = 10

	// Pie chart
	pc := widgets.NewPieChart()
	pc.Title = "Disk Usage"
	pc.Data = []float64{10, 20, 30, 40}
	pc.Colors = []ui.Color{ui.ColorRed, ui.ColorYellow, ui.ColorGreen, ui.ColorBlue}
	pc.LabelFormatter = func(i int, v float64) string { return fmt.Sprintf("%.0f%%", v) }
	pc.TitleStyle.Fg = ui.ColorMagenta
	pc.BorderStyle.Fg = ui.ColorMagenta
	pc.BorderRounded = true

	// Plot
	d.plotData = make([][]float64, 2)
	d.plotData[0] = make([]float64, 50)
	d.plotData[1] = make([]float64, 50)

	d.plot = widgets.NewPlot()
	d.plot.Title = "Response Time"
	d.plot.TitleBottom = "(ms)"
	d.plot.Data = d.plotData
	d.plot.AxesColor = ui.ColorWhite
	d.plot.LineColors[0] = ui.ColorGreen
	d.plot.LineColors[1] = ui.ColorYellow
	d.plot.Marker = widgets.MarkerDot
	d.plot.TitleStyle.Fg = ui.ColorCyan
	d.plot.BorderStyle.Fg = ui.ColorCyan
	d.plot.BorderRounded = true

	// Logs list
	d.logs = widgets.NewList()
	d.logs.Title = "System Logs"
	d.logs.Rows = []string{
		"[INFO] System started",
		"[INFO] Service A initialized",
		"[WARN] Connection timeout (retrying)",
		"[ERROR] User authentication failed",
		"[INFO] Worker pool scaled up",
		"[WARN] High memory usage detected",
		"[INFO] Health check passed",
		"[ERROR] Disk space low (<10%)",
	}
	d.logs.TextStyle.Fg = ui.ColorYellow
	d.logs.SelectedStyle = ui.NewStyle(ui.ColorBlack, ui.ColorYellow)
	d.logs.TitleStyle.Fg = ui.ColorYellow
	d.logs.BorderStyle.Fg = ui.ColorYellow
	d.logs.TitleBottom = "Wheel to scroll"
	d.logs.TitleBottomAlignment = ui.AlignRight
	d.logs.BorderRounded = true

	// Gauge
	d.g = widgets.NewGauge()
	d.g.Title = "Download"
	d.g.Percent = 50
	d.g.BarColor = ui.ColorGreen
	d.g.BorderStyle.Fg = ui.ColorGreen
	d.g.TitleStyle.Fg = ui.ColorGreen
	d.g.BorderRounded = true

	// Grid
	d.grid = ui.NewGrid()
	d.grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0, p),
		),
		ui.NewRow(2.0/10,
			ui.NewCol(1.0/2, slg),
			ui.NewCol(1.0/2,
				ui.NewRow(1.0/2, d.lg1),
				ui.NewRow(1.0/2, d.lg2),
			),
		),
		ui.NewRow(3.5/10,
			ui.NewCol(1.0/3, d.bc),
			ui.NewCol(1.0/3, pc),
			ui.NewCol(1.0/3, d.plot),
		),
		ui.NewRow(3.5/10,
			ui.NewCol(2.0/3, d.logs),
			ui.NewCol(1.0/3, d.g),
		),
	)

	return d
}

func (d *dashboard) onResize(w, h int) {
	d.grid.SetRect(0, 0, w, h)
	ui.Clear()
	ui.Render(d.grid)
}

func (d *dashboard) onTick() (dirty bool) {
	d.tickCount++

	// Sparklines
	d.sl.Data = append(d.sl.Data[1:], float64(rand.Intn(100)))
	d.sl2.Data = append(d.sl2.Data[1:], float64(rand.Intn(100)))
	dirty = true

	// Gauges (slower)
	if d.tickCount%5 == 0 {
		d.lg1.Percent = (d.lg1.Percent + rand.Intn(5)) % 100
		d.lg2.Percent = (d.lg2.Percent + rand.Intn(3)) % 100
	}

	// Bar chart (slower)
	if d.tickCount%10 == 0 {
		for i := range d.bc.Data {
			d.bc.Data[i] = float64(rand.Intn(10))
		}
	}

	// Download gauge
	d.g.Percent = (d.g.Percent + 2) % 100

	// Plot
	d.plotData[0] = append(d.plotData[0][1:], 20+10*math.Sin(float64(d.tickCount)/10.0)+float64(rand.Intn(5)))
	d.plotData[1] = append(d.plotData[1][1:], 40+20*math.Cos(float64(d.tickCount)/15.0)+float64(rand.Intn(10)))
	d.plot.Data = d.plotData

	return true
}

// ---- SSH session runner ----

func runDashboardOverSSH(sess ssh.Session) {
	oneSession.Lock()
	defer oneSession.Unlock()

	tty, err := newSessionTTY(sess)
	if err != nil {
		fmt.Fprintln(sess.Stderr(), err)
		return
	}
	defer tty.Close()

	// Use new InitWithConfig API - no tcell exposure needed!
	err = ui.InitWithConfig(&ui.InitConfig{
		CustomTTY: tty, // sessionTTY implements io.ReadWriter
	})
	if err != nil {
		fmt.Fprintln(sess.Stderr(), "init:", err)
		return
	}
	defer ui.Close()

	d := newDashboard()
	w, h := tty.Size()
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	d.onResize(w, h)

	events := ui.PollEventsWithContext(sess.Context())
	ticker := time.NewTicker(150 * time.Millisecond) // ~6-7 FPS feels nicer over SSH
	defer ticker.Stop()

	for {
		select {
		case <-sess.Context().Done():
			return

		case e, ok := <-events:
			if !ok {
				return // channel closed
			}
			switch e.ID {
			case "q", "<C-c>", "<Escape>":
				return

			case "<Resize>":
				r := e.Payload.(ui.Resize)
				d.onResize(r.Width, r.Height)

			case "<MouseWheelUp>":
				d.logs.ScrollUp()
				ui.Render(d.grid)

			case "<MouseWheelDown>":
				d.logs.ScrollDown()
				ui.Render(d.grid)
			}

		case <-ticker.C:
			if d.onTick() {
				ui.Render(d.grid)
			}
		}
	}
}

func main() {
	ssh.Handle(runDashboardOverSSH)
	log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("hostkey")))
}
