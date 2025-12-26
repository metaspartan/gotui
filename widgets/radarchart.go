package widgets

import (
	"image"
	"math"

	ui "github.com/metaspartan/gotui/v5"
)

// RadarChart represents a widget that displays a radar chart.
type RadarChart struct {
	ui.Block
	Data       [][]float64
	DataLabels []string
	Labels     []string
	MaxVal     float64
	LineColors []ui.Color
	LabelStyle ui.Style
	DotStyle   ui.Style
}

// NewRadarChart returns a new RadarChart.
func NewRadarChart() *RadarChart {
	return &RadarChart{
		Block:      *ui.NewBlock(),
		LineColors: ui.Theme.Plot.Lines,
		LabelStyle: ui.NewStyle(ui.Theme.Plot.Axes),
		DotStyle:   ui.NewStyle(ui.ColorWhite),
		Data:       [][]float64{},
	}
}

// Draw draws the radar chart to the buffer.
func (rc *RadarChart) Draw(buf *ui.Buffer) {
	rc.Block.Draw(buf)
	if len(rc.Data) == 0 || len(rc.Data[0]) == 0 {
		return
	}
	canvas := ui.NewCanvas()
	canvas.SetRect(rc.Inner.Min.X, rc.Inner.Min.Y, rc.Inner.Max.X, rc.Inner.Max.Y)
	canvas.Border = false
	w := rc.Inner.Dx() * 2
	h := rc.Inner.Dy() * 4
	bcx := float64(w) / 2.0
	bcy := float64(h) / 2.0
	radius := math.Min(bcx, bcy) - 20.0
	numAxes := len(rc.Data[0])
	angleStep := (2 * math.Pi) / float64(numAxes)
	for i := 0; i < numAxes; i++ {
		angle := float64(i)*angleStep - (math.Pi / 2)
		ex := bcx + math.Cos(angle)*radius
		ey := bcy + math.Sin(angle)*radius
		canvas.SetLine(
			image.Pt(int(bcx), int(bcy)),
			image.Pt(int(ex), int(ey)),
			ui.ColorWhite,
		)
		if i < len(rc.Labels) {
			lx := rc.Inner.Min.X + int(ex/2)
			ly := rc.Inner.Min.Y + int(ey/4)
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
	for i, dataSet := range rc.Data {
		color := ui.SelectColor(rc.LineColors, i)
		if color == ui.ColorClear || color == 0 {
			color = ui.ColorWhite
		}
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
		}
		canvas.SetLine(lastPoint, firstPoint, color)
	}
	canvas.Draw(buf)
}
