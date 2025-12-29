package main

import (
	"fmt"
	"log"
	"math"
	"time"

	ui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize gotui: %v", err)
	}
	defer ui.Close()

	w, h := ui.TerminalDimensions()

	// Helper to create charts
	createChart := func(title string) *widgets.StepChart {
		lc := widgets.NewStepChart()
		lc.Title = title
		lc.AxesColor = ui.ColorWhite
		lc.LineColors = []ui.Color{ui.ColorGreen, ui.ColorYellow}
		return lc
	}

	// 1. Top Left: No Axes (Simple)
	p0 := createChart("Simple (No Axes)")
	p0.ShowAxes = false

	// 2. Top Right: Standard Axes
	p1 := createChart("Standard Axes")
	p1.ShowAxes = true

	// 3. Middle Left: Float Labels
	p2 := createChart("Float Labels")
	p2.ShowAxes = false
	p2.ShowRightAxis = true

	// 4. Middle Right: Axes + Labels
	p3 := createChart("Axes + Labels")
	p3.ShowAxes = true
	p3.ShowRightAxis = true
	p3.DataLabels = []string{"CPU", "MEM"}

	// 5. Bottom: Full Width with Black Background
	p4 := createChart("Full Width (Black BG)")
	p4.ShowAxes = true
	p4.ShowRightAxis = true
	p4.DataLabels = []string{"Line 1", "Line 2"}
	p4.BorderStyle.Bg = ui.ColorBlack
	p4.TitleStyle.Bg = ui.ColorBlack
	p4.BackgroundColor = ui.ColorBlack
	p4.TitleStyle.Fg = ui.ColorWhite

	resize := func() {
		halfW := w / 2
		thirdH := h / 3
		// Top row (2 quadrants)
		p0.SetRect(0, 0, halfW, thirdH)
		p1.SetRect(halfW, 0, w, thirdH)
		// Middle row (2 quadrants)
		p2.SetRect(0, thirdH, halfW, thirdH*2)
		p3.SetRect(halfW, thirdH, w, thirdH*2)
		// Bottom row (full width)
		p4.SetRect(0, thirdH*2, w, h)
	}
	resize()

	// Helper to maintain data buffer size
	updateData := func(width int, baseData []float64) []float64 {
		needed := width * 2
		if len(baseData) < needed {
			diff := needed - len(baseData)
			lastVal := 0.0
			if len(baseData) > 0 {
				lastVal = baseData[len(baseData)-1]
			}
			for k := 0; k < diff; k++ {
				baseData = append(baseData, lastVal)
			}
		}
		return baseData
	}

	// Calculate available drawing width for a widget
	getDrawWidth := func(widget *widgets.StepChart) int {
		rect := widget.GetRect()
		width := rect.Dx() - 2 // Account for borders
		if widget.ShowAxes {
			width -= 5 // yAxisLabelsWidth + 1
		}
		if width < 1 {
			width = 1
		}
		return width
	}

	getVisibleData := func(d []float64, widget *widgets.StepChart) []float64 {
		width := getDrawWidth(widget)
		if len(d) > width {
			return d[len(d)-width:]
		}
		return d
	}

	render := func(data1, data2 []float64) {
		p0.Data = [][]float64{getVisibleData(data1, p0), getVisibleData(data2, p0)}
		p1.Data = [][]float64{getVisibleData(data1, p1), getVisibleData(data2, p1)}
		p2.Data = [][]float64{getVisibleData(data1, p2), getVisibleData(data2, p2)}
		p3.Data = [][]float64{getVisibleData(data1, p3), getVisibleData(data2, p3)}
		p4.Data = [][]float64{getVisibleData(data1, p4), getVisibleData(data2, p4)}
		ui.Render(p0, p1, p2, p3, p4)
	}

	// Initial data
	data := make([]float64, w) // Start with plenty
	data2 := make([]float64, w)
	for i := range data {
		data[i] = 15 + 10*math.Sin(float64(i)/5) + 5*math.Cos(float64(i)/3)
		data2[i] = 10 + 8*math.Sin(float64(i)/4) + 2*math.Cos(float64(i)/2)
	}

	render(data, data2)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(100 * time.Millisecond)
	x := float64(w)

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				w, h = payload.Width, payload.Height
				resize()

				// Ensure buffers are enough for the new max width
				data = updateData(w, data)
				data2 = updateData(w, data2)

				ui.Clear()
				render(data, data2)
			}
		case <-ticker.C:
			val := 15 + 10*math.Sin(x/5) + 5*math.Cos(x/3) + float64(time.Now().UnixNano()%10)*0.5
			val2 := 10 + 8*math.Sin(x/4) + 2*math.Cos(x/2) + float64(time.Now().UnixNano()%10)*0.5

			data = append(data, val)
			data2 = append(data2, val2)

			// Trim excess
			if len(data) > 4000 {
				data = data[1:]
				data2 = data2[1:]
			}

			// Update dynamic custom labels for p2
			p2.DataLabels = []string{
				fmt.Sprintf("%.1f%%", val),
				fmt.Sprintf("%.1f%%", val2),
			}

			render(data, data2)
			x += 0.5
		}
	}
}
