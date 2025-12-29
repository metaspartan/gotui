package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui/v5"
)

// StepChart represents a stepped line chart
type StepChart struct {
	*Plot
	ShowRightAxis bool
	DataLabels    []string
}

// NewStepChart returns a new StepChart initialized with default settings
// for a solid-line step graph.
func NewStepChart() *StepChart {
	p := NewPlot()
	p.AxesColor = ui.ColorWhite
	p.LineColors = []ui.Color{ui.ColorGreen}
	p.ShowAxes = true

	return &StepChart{
		Plot: p,
	}
}

// Draw draws the StepChart to the buffer using box drawing characters for a stepped line.
func (sc *StepChart) Draw(buf *ui.Buffer) {
	sc.Block.Draw(buf)

	maxVal := sc.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64From2dSlice(sc.Data)
	}
	if maxVal == 0 {
		maxVal = 1
	}

	if sc.ShowAxes {
		sc.Plot.plotAxes(buf, maxVal)
	}

	drawArea := sc.getDrawArea()
	height := drawArea.Dy()

	scale := sc.HorizontalScale
	if scale <= 0 {
		scale = 1
	}

	for i, lineData := range sc.Data {
		lineColor := ui.SelectColor(sc.LineColors, i)
		style := ui.NewStyle(lineColor, sc.Block.BorderStyle.Bg)
		sc.drawLine(buf, lineData, drawArea, maxVal, float64(height), scale, style)

		if sc.ShowRightAxis {
			sc.drawLabel(buf, lineData, i, drawArea, maxVal, float64(height), style)
		}
	}
}

func (sc *StepChart) drawLabel(buf *ui.Buffer, lineData []float64, lineIdx int, drawArea image.Rectangle, maxVal, height float64, style ui.Style) {
	if len(lineData) == 0 {
		return
	}

	val := lineData[len(lineData)-1]
	y := drawArea.Max.Y - 1 - int((val/maxVal)*(height-1))
	y = clamp(y, drawArea.Min.Y, drawArea.Max.Y-1)

	label := ""
	if lineIdx < len(sc.DataLabels) {
		label = sc.DataLabels[lineIdx]
	} else {
		label = fmt.Sprintf("%.2f", val)
	}

	x := max(drawArea.Max.X-len(label), drawArea.Min.X)

	buf.SetString(label, style, image.Pt(x, y))
}

func (sc *StepChart) getDrawArea() image.Rectangle {
	drawArea := sc.Inner
	if sc.ShowAxes {
		drawArea = image.Rect(
			sc.Inner.Min.X+yAxisLabelsWidth+1, sc.Inner.Min.Y,
			sc.Inner.Max.X, sc.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}
	return drawArea
}

func (sc *StepChart) drawLine(buf *ui.Buffer, lineData []float64, drawArea image.Rectangle, maxVal, height float64, scale int, style ui.Style) {

	for j := 0; j < len(lineData)-1; j++ {
		x1 := drawArea.Min.X + j*scale
		x2 := drawArea.Min.X + (j+1)*scale

		if x1 >= drawArea.Max.X {
			break
		}

		val1 := lineData[j]
		val2 := lineData[j+1]

		y1 := drawArea.Max.Y - 1 - int((val1/maxVal)*(height-1))
		y2 := drawArea.Max.Y - 1 - int((val2/maxVal)*(height-1))

		y1 = clamp(y1, drawArea.Min.Y, drawArea.Max.Y-1)
		y2 = clamp(y2, drawArea.Min.Y, drawArea.Max.Y-1)

		startH := x1
		if j > 0 {
			prevVal := lineData[j-1]
			yPrev := drawArea.Max.Y - 1 - int((prevVal/maxVal)*(height-1))
			yPrev = clamp(yPrev, drawArea.Min.Y, drawArea.Max.Y-1)
			if yPrev != y1 {
				startH = x1 + 1
			}
		}

		sc.drawSegment(buf, startH, x2, y1, y2, drawArea.Max.X, style)
	}

	lastIdx := len(lineData) - 1
	if lastIdx >= 0 {
		x1 := drawArea.Min.X + lastIdx*scale
		val := lineData[lastIdx]
		y := drawArea.Max.Y - 1 - int((val/maxVal)*(height-1))
		y = clamp(y, drawArea.Min.Y, drawArea.Max.Y-1)

		startH := x1
		if lastIdx > 0 {
			prevVal := lineData[lastIdx-1]
			yPrev := drawArea.Max.Y - 1 - int((prevVal/maxVal)*(height-1))
			yPrev = clamp(yPrev, drawArea.Min.Y, drawArea.Max.Y-1)
			if yPrev != y {
				startH = x1 + 1
			}
		}

		for x := startH; x < drawArea.Max.X; x++ {
			buf.SetCell(ui.NewCell(ui.HORIZONTAL_LINE, style), image.Pt(x, y))
		}
	}
}

func (sc *StepChart) drawSegment(buf *ui.Buffer, startX, endX, y1, y2, maxX int, style ui.Style) {
	limitX := min(endX, maxX)
	for x := startX; x < limitX; x++ {
		buf.SetCell(ui.NewCell(ui.HORIZONTAL_LINE, style), image.Pt(x, y1))
	}

	if endX < maxX {
		sc.drawVerticalTransition(buf, endX, y1, y2, style)
	}
}

func (sc *StepChart) drawVerticalTransition(buf *ui.Buffer, x, y1, y2 int, style ui.Style) {
	if y1 == y2 {
		buf.SetCell(ui.NewCell(ui.HORIZONTAL_LINE, style), image.Pt(x, y1))
		return
	}

	minY := min(y1, y2)
	maxY := max(y1, y2)

	for y := minY + 1; y < maxY; y++ {
		buf.SetCell(ui.NewCell(ui.VERTICAL_LINE, style), image.Pt(x, y))
	}

	sc.drawCorners(buf, x, y1, y2, style)
}

func (sc *StepChart) drawCorners(buf *ui.Buffer, x, y1, y2 int, style ui.Style) {
	var startChar rune
	if y2 < y1 {
		startChar = ui.BOTTOM_RIGHT // ┘
	} else {
		startChar = ui.TOP_RIGHT // ┐
	}
	buf.SetCell(ui.NewCell(startChar, style), image.Pt(x, y1))

	var endChar rune
	if y1 > y2 {
		endChar = ui.TOP_LEFT // ┌
	} else {
		endChar = ui.BOTTOM_LEFT // └
	}
	buf.SetCell(ui.NewCell(endChar, style), image.Pt(x, y2))
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
