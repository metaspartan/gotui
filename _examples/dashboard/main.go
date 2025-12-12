package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	ui "github.com/metaspartan/gotui/v4"
	"github.com/metaspartan/gotui/v4/widgets"
)

func main() {
	// 1. Header
	p := widgets.NewParagraph()
	p.Title = "gotui Dashboard"
	p.Text = "PRESS q TO QUIT | Responsive Grid Layout Demo | TrueColor Support"
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorLightCyan
	p.TitleStyle = ui.NewStyle(ui.ColorLightCyan, ui.ColorClear, ui.ModifierBold)
	p.TitleAlignment = ui.AlignCenter
	p.TitleRight = "v4.0.0"
	p.BorderRounded = false // Variety: Non-rounded border

	// 2. Sparklines (CPU Usage)
	slData := make([]float64, 200)
	sl := widgets.NewSparkline()
	sl.Data = slData

	// 4. Bar Chart (Network Traffic)
	bc := widgets.NewBarChart()
	bc.Title = "Network Traffic"
	bc.TitleBottom = "MB/s" // Bottom title demo
	bc.TitleBottomAlignment = ui.AlignRight
	bc.Labels = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	bc.BarColors = []ui.Color{ui.ColorBlue, ui.ColorLightCyan}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	bc.Data = []float64{3, 2, 5, 3, 9, 5}
	bc.TitleStyle.Fg = ui.ColorBlue
	bc.BorderStyle.Fg = ui.ColorBlue
	bc.BorderRounded = true
	bc.BarWidth = 0 // Auto-Width
	bc.BarGap = 1
	bc.MaxVal = 10 // Fix scaling to be consistent
	sl.LineColor = ui.ColorGreen
	sl.TitleStyle.Fg = ui.ColorWhite
	sl.MaxVal = 100

	sl2 := widgets.NewSparkline()
	sl2.Data = slData
	sl2.TitleStyle.Fg = ui.ColorWhite
	sl2.LineColor = ui.ColorMagenta
	sl2.MaxVal = 100

	slg := widgets.NewSparklineGroup(sl, sl2)
	slg.Title = "CPU Usage"
	slg.TitleStyle.Fg = ui.ColorGreen
	slg.BorderStyle.Fg = ui.ColorGreen
	slg.TitleRight = "Core 0 & 1" // Right title demo
	slg.BorderRounded = true

	// 3. Line Gauges (Memory & Load)
	lg1 := widgets.NewLineGauge()
	lg1.Title = "Memory"
	lg1.Percent = 45
	lg1.BarRune = '■'
	lg1.LineColor = ui.ColorYellow
	lg1.TitleStyle.Fg = ui.ColorYellow
	lg1.BorderRounded = true

	lg2 := widgets.NewLineGauge()
	lg2.Title = "Load"
	lg2.Percent = 60
	lg2.BarRune = '▰' // Future style
	lg2.BarRuneEmpty = '▱'
	lg2.LineColor = ui.ColorRed
	lg2.TitleStyle.Fg = ui.ColorRed
	lg2.BorderRounded = true

	// 5. Pie Chart (Disk Space)
	pc := widgets.NewPieChart()
	pc.Title = "Disk Usage"
	pc.Data = []float64{10, 20, 30, 40}
	pc.Colors = []ui.Color{ui.ColorRed, ui.ColorYellow, ui.ColorGreen, ui.ColorBlue}
	pc.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.0f%%", v)
	}
	pc.TitleStyle.Fg = ui.ColorMagenta
	pc.BorderStyle.Fg = ui.ColorMagenta
	pc.BorderRounded = true

	// 6. Plot (Response Time)
	plotData := make([][]float64, 2)
	plotData[0] = make([]float64, 50)
	plotData[1] = make([]float64, 50)

	plot := widgets.NewPlot()
	plot.Title = "Response Time"
	plot.TitleBottom = "(ms)"
	plot.Data = plotData
	plot.AxesColor = ui.ColorWhite
	plot.LineColors[0] = ui.ColorGreen  // Brighter color
	plot.LineColors[1] = ui.ColorYellow // Brighter color
	plot.Marker = widgets.MarkerDot     // Clearer than Braille
	plot.TitleStyle.Fg = ui.ColorLightCyan
	plot.BorderStyle.Fg = ui.ColorLightCyan
	plot.BorderRounded = true

	// 7. List (System Logs)
	l := widgets.NewList()
	l.Title = "System Logs"
	l.Rows = []string{
		"[INFO] System started",
		"[INFO] Service A initialized",
		"[WARN] Connection timeout (retrying)",
		"[INFO] Cache cleared",
		"[ERROR] User authentication failed",
		"[INFO] Worker pool scaled up",
		"[INFO] Backup completed",
		"[INFO] Health check passed",
		"[WARN] High memory usage detected",
		"[INFO] GC triggered",
		"[INFO] New client connected (192.168.1.10)",
		"[INFO] Request processed in 45ms",
		"[ERROR] Database connection lost",
		"[INFO] Reconnecting to DB...",
		"[INFO] DB Connected",
		"[INFO] Job #1234 completed",
		"[WARN] Rate limit exceeded",
		"[INFO] Retrying Job #1234",
		"[INFO] Cache warm-up started",
		"[INFO] Cache warm-up finished",
		"[INFO] Service B initialized",
		"[INFO] API Gateway ready",
		"[INFO] Listening on port 8080",
		"[WARN] Deprecated API usage /v1/users",
		"[INFO] User A logged out",
		"[INFO] User B logged in",
		"[INFO] Metrics flushed to InfluxDB",
		"[INFO] Rotating logs",
		"[INFO] Running cleanup task",
		"[INFO] Temp files deleted",
		"[ERROR] Disk space low (<10%)",
		"[INFO] Alert sent to admin",
	}
	l.TextStyle.Fg = ui.ColorYellow
	l.SelectedStyle = ui.NewStyle(ui.ColorBlack, ui.ColorYellow)
	l.TitleStyle.Fg = ui.ColorYellow
	l.BorderStyle.Fg = ui.ColorYellow
	l.TitleBottom = "Page 1/1"
	l.TitleBottomAlignment = ui.AlignRight
	l.BorderRounded = true

	// 8. Gauge (Download)
	g := widgets.NewGauge()
	g.Title = "Download"
	g.Percent = 50
	g.BarColor = ui.ColorGreen
	g.BorderStyle.Fg = ui.ColorGreen
	g.TitleStyle.Fg = ui.ColorGreen
	g.BorderRounded = true

	// 9. Grid Layout
	grid := ui.NewGrid()
	//termWidth, termHeight := ui.TerminalDimensions()
	//grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/10,
			ui.NewCol(1.0, p), // Header
		),
		ui.NewRow(2.0/10,
			ui.NewCol(1.0/2, slg), // Sparklines
			ui.NewCol(1.0/2, // Gauges Stacked
				ui.NewRow(1.0/2, lg1),
				ui.NewRow(1.0/2, lg2),
			),
		),
		ui.NewRow(3.5/10, // Increased Height
			ui.NewCol(1.0/3, bc),   // Bar Chart
			ui.NewCol(1.0/3, pc),   // Pie Chart
			ui.NewCol(1.0/3, plot), // Plot
		),
		ui.NewRow(3.5/10,
			ui.NewCol(2.0/3, l), // Logs
			ui.NewCol(1.0/3, g), // Gauge
		),
	)

	// Check for screenshot mode
	if len(os.Args) > 1 && os.Args[1] == "-screenshot" {
		termWidth, termHeight := 1024/7, 768/13 // Approximate cols/rows (e.g. 146x59)
		// Or standard terminal size
		termWidth, termHeight = 120, 40

		grid.SetRect(0, 0, termWidth, termHeight)

		if err := ui.SaveImage("screenshot.png", termWidth, termHeight, grid); err != nil {
			log.Fatalf("failed to take screenshot: %v", err)
		}
		fmt.Println("Screenshot saved to screenshot.png")
		return
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Render(grid)

	// Update Function
	tickerCount := 0
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(100 * time.Millisecond).C // 10 FPS updates

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<MouseWheelUp>":
				l.ScrollUp()
				ui.Render(grid)
			case "<MouseWheelDown>":
				l.ScrollDown()
				ui.Render(grid)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			tickerCount++

			// Update Sparklines
			sl.Data = append(sl.Data[1:], float64(rand.Intn(100)))
			sl2.Data = append(sl2.Data[1:], float64(rand.Intn(100)))

			// Update Gauges
			if tickerCount%5 == 0 {
				lg1.Percent = (lg1.Percent + rand.Intn(5)) % 100
				lg2.Percent = (lg2.Percent + rand.Intn(3)) % 100
			}

			// Update BarChart
			if tickerCount%10 == 0 {
				for i := range bc.Data {
					bc.Data[i] = float64(rand.Intn(10))
				}
			}

			// Update Gauge
			g.Percent = (g.Percent + 2) % 100

			// Update Plot
			plotData[0] = append(plotData[0][1:], 20+10*math.Sin(float64(tickerCount)/10.0)+float64(rand.Intn(5)))
			plotData[1] = append(plotData[1][1:], 40+20*math.Cos(float64(tickerCount)/15.0)+float64(rand.Intn(10)))
			plot.Data = plotData

			// Update Logs (Rotate selection)
			// if tickerCount%20 == 0 {
			// 	l.SelectedRow = (l.SelectedRow + 1) % len(l.Rows)
			// }

			ui.Render(grid)
		}
	}
}
