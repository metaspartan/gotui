package widgets

import (
	"image"
	"math"

	ui "github.com/metaspartan/gotui"
)

type RadarChart struct {
	ui.Block
	Data       [][]float64
	DataLabels []string // Names for each data slice
	Labels     []string // Names for axes
	MaxVal     float64
	LineColors []ui.Color
	LabelStyle ui.Style
	DotStyle   ui.Style
}

func NewRadarChart() *RadarChart {
	return &RadarChart{
		Block:      *ui.NewBlock(),
		LineColors: ui.Theme.Plot.Lines,
		LabelStyle: ui.NewStyle(ui.Theme.Plot.Axes), // Fix: Wrap color in Style
		DotStyle:   ui.NewStyle(ui.ColorWhite),
		Data:       [][]float64{},
	}
}

func (rc *RadarChart) Draw(buf *ui.Buffer) {
	rc.Block.Draw(buf)

	if len(rc.Data) == 0 || len(rc.Data[0]) == 0 {
		return
	}

	// Use inner area
	canvas := ui.NewCanvas()
	canvas.Rectangle = rc.Inner

	// Braille canvas is 2x width, 4x height of cell coordinates.
	w := rc.Inner.Dx() * 2
	h := rc.Inner.Dy() * 4

	bcx := float64(w) / 2.0
	bcy := float64(h) / 2.0
	radius := math.Min(bcx, bcy) - 20.0 // Increased padding for labels

	numAxes := len(rc.Data[0])
	angleStep := (2 * math.Pi) / float64(numAxes)

	// Draw Axes
	for i := 0; i < numAxes; i++ {
		angle := float64(i)*angleStep - (math.Pi / 2) // Start at top
		ex := bcx + math.Cos(angle)*radius
		ey := bcy + math.Sin(angle)*radius

		canvas.SetLine(
			image.Pt(int(bcx), int(bcy)),
			image.Pt(int(ex), int(ey)),
			ui.ColorWhite,
		)

		// Labels (standard cell coordinates)
		if i < len(rc.Labels) {
			lx := rc.Inner.Min.X + int(ex/2)
			ly := rc.Inner.Min.Y + int(ey/4)

			// Simple adjustments to keep labels from overlapping chart too much
			if conversionX := math.Cos(angle); conversionX > 0.5 {
				lx++
			} else if conversionX < -0.5 {
				lx -= len(rc.Labels[i])
			}
			if conversionY := math.Sin(angle); conversionY < -0.5 {
				ly--
			} else if conversionY > 0.5 {
				ly++
			}

			buf.SetString(
				rc.Labels[i],
				rc.LabelStyle,
				image.Pt(lx, ly),
			)
		}
	}

	maxVal := rc.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64From2dSlice(rc.Data)
	}

	// Draw Data Polygons
	for i, dataSet := range rc.Data {
		color := ui.SelectColor(rc.LineColors, i)
		var firstPoint image.Point
		var lastPoint image.Point

		for j, val := range dataSet {
			angle := float64(j)*angleStep - (math.Pi / 2)
			valRadius := (val / maxVal) * radius

			px := bcx + math.Cos(angle)*valRadius
			py := bcy + math.Sin(angle)*valRadius
			p := image.Pt(int(px), int(py))

			if j == 0 {
				firstPoint = p
			} else {
				canvas.SetLine(lastPoint, p, color)
			}
			lastPoint = p

			// Optional: draw dot at vertex
			// canvas.SetPoint(p, color)
		}
		// Close the loop
		canvas.SetLine(lastPoint, firstPoint, color)
	}

	// Draw to buffer
	canvas.Draw(buf)
}
