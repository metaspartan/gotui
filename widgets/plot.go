package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Plot has two modes: line(default) and scatter.
// Plot also has two marker types: braille(default) and dot.
// A single braille character is a 2x4 grid of dots, so using braille
// gives 2x X resolution and 4x Y resolution over dot mode.
type Plot struct {
	ui.Block

	Data       [][]float64
	DataLabels []string
	MaxVal     float64

	LineColors []ui.Color
	AxesColor  ui.Color // TODO
	ShowAxes   bool
	Fill       bool // Fills the area under the line (like nvtop)

	Marker          PlotMarker
	DotMarkerRune   rune
	PlotType        PlotType
	HorizontalScale int
	DrawDirection   DrawDirection // TODO
}

const (
	xAxisLabelsHeight = 1
	yAxisLabelsWidth  = 4
	xAxisLabelsGap    = 2
	yAxisLabelsGap    = 1
)

type PlotType uint

const (
	LineChart PlotType = iota
	ScatterPlot
)

type PlotMarker uint

const (
	MarkerBraille PlotMarker = iota
	MarkerDot
)

type DrawDirection uint

const (
	DrawLeft DrawDirection = iota
	DrawRight
)

func NewPlot() *Plot {
	return &Plot{
		Block:           *ui.NewBlock(),
		LineColors:      ui.Theme.Plot.Lines,
		AxesColor:       ui.Theme.Plot.Axes,
		Marker:          MarkerBraille,
		DotMarkerRune:   ui.DOT,
		Data:            [][]float64{},
		HorizontalScale: 1,
		DrawDirection:   DrawRight,
		ShowAxes:        true,
		PlotType:        LineChart,
		Fill:            false,
	}
}

func (plt *Plot) renderBraille(buf *ui.Buffer, drawArea image.Rectangle, maxVal float64) {
	canvas := ui.NewCanvas()
	canvas.SetRect(drawArea.Min.X, drawArea.Min.Y, drawArea.Max.X, drawArea.Max.Y)
	canvas.Border = false

	switch plt.PlotType {
	case ScatterPlot:
		for i, line := range plt.Data {
			for j, val := range line {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				canvas.SetPoint(
					image.Pt(
						(drawArea.Min.X+(j*plt.HorizontalScale))*2,
						(drawArea.Max.Y-height-1)*4,
					),
					ui.SelectColor(plt.LineColors, i),
				)
			}
		}
	case LineChart:
		for i, line := range plt.Data {
			previousHeight := int((line[1] / maxVal) * float64(drawArea.Dy()-1))
			for j, val := range line[1:] {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))

				// Calculate coordinates for the line segment
				x1 := (drawArea.Min.X + (j * plt.HorizontalScale)) * 2
				y1 := (drawArea.Max.Y - previousHeight - 1) * 4
				x2 := (drawArea.Min.X + ((j + 1) * plt.HorizontalScale)) * 2
				y2 := (drawArea.Max.Y - height - 1) * 4

				color := ui.SelectColor(plt.LineColors, i)
				if color == ui.ColorClear || color == 0 {
					color = ui.ColorWhite
				}

				canvas.SetLine(
					image.Pt(x1, y1),
					image.Pt(x2, y2),
					color,
				)

				if plt.Fill {
					// Fill area under the line
					// We iterate from x1 to x2 and draw a vertical line from the curve to the bottom
					bottomY := (drawArea.Max.Y-1)*4 + 3 // Bottom of the braille canvas

					// Avoid divide by zero if x1 == x2 (vertical line)
					if x2 > x1 {
						slope := float64(y2-y1) / float64(x2-x1)
						for x := x1; x < x2; x++ {
							// Interpolate Y
							y := float64(y1) + slope*float64(x-x1)
							canvas.SetLine(
								image.Pt(x, int(y)),
								image.Pt(x, bottomY),
								color,
							)
						}
					} else {
						// Just draw vertical line
						canvas.SetLine(
							image.Pt(x1, y1),
							image.Pt(x1, bottomY),
							color,
						)
					}
				}

				previousHeight = height
			}
		}
	}

	canvas.Draw(buf)
}

func (plt *Plot) renderDot(buf *ui.Buffer, drawArea image.Rectangle, maxVal float64) {
	switch plt.PlotType {
	case ScatterPlot:
		for i, line := range plt.Data {
			for j, val := range line {
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				point := image.Pt(drawArea.Min.X+(j*plt.HorizontalScale), drawArea.Max.Y-1-height)
				if point.In(drawArea) {
					buf.SetCell(
						ui.NewCell(plt.DotMarkerRune, ui.NewStyle(ui.SelectColor(plt.LineColors, i))),
						point,
					)
				}
			}
		}
	case LineChart:
		for i, line := range plt.Data {
			for j := 0; j < len(line) && j*plt.HorizontalScale < drawArea.Dx(); j++ {
				val := line[j]
				height := int((val / maxVal) * float64(drawArea.Dy()-1))
				buf.SetCell(
					ui.NewCell(plt.DotMarkerRune, ui.NewStyle(ui.SelectColor(plt.LineColors, i))),
					image.Pt(drawArea.Min.X+(j*plt.HorizontalScale), drawArea.Max.Y-1-height),
				)
			}
		}
	}
}

func (plt *Plot) plotAxes(buf *ui.Buffer, maxVal float64) {
	// draw origin cell
	buf.SetCell(
		ui.NewCell(ui.BOTTOM_LEFT, ui.NewStyle(ui.ColorWhite)),
		image.Pt(plt.Inner.Min.X+yAxisLabelsWidth, plt.Inner.Max.Y-xAxisLabelsHeight-1),
	)
	// draw x axis line
	for i := yAxisLabelsWidth + 1; i < plt.Inner.Dx(); i++ {
		buf.SetCell(
			ui.NewCell(ui.HORIZONTAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(i+plt.Inner.Min.X, plt.Inner.Max.Y-xAxisLabelsHeight-1),
		)
	}
	// draw y axis line
	for i := 0; i < plt.Inner.Dy()-xAxisLabelsHeight-1; i++ {
		buf.SetCell(
			ui.NewCell(ui.VERTICAL_DASH, ui.NewStyle(ui.ColorWhite)),
			image.Pt(plt.Inner.Min.X+yAxisLabelsWidth, i+plt.Inner.Min.Y),
		)
	}
	// draw x axis labels
	// draw 0
	buf.SetString(
		"0",
		ui.NewStyle(ui.ColorWhite),
		image.Pt(plt.Inner.Min.X+yAxisLabelsWidth, plt.Inner.Max.Y-1),
	)
	// draw rest
	for x := plt.Inner.Min.X + yAxisLabelsWidth + (xAxisLabelsGap)*plt.HorizontalScale + 1; x < plt.Inner.Max.X-1; {
		label := fmt.Sprintf(
			"%d",
			(x-(plt.Inner.Min.X+yAxisLabelsWidth)-1)/(plt.HorizontalScale)+1,
		)
		buf.SetString(
			label,
			ui.NewStyle(ui.ColorWhite),
			image.Pt(x, plt.Inner.Max.Y-1),
		)
		x += (len(label) + xAxisLabelsGap) * plt.HorizontalScale
	}
	// draw y axis labels
	verticalScale := maxVal / float64(plt.Inner.Dy()-xAxisLabelsHeight-1)
	for i := 0; i*(yAxisLabelsGap+1) < plt.Inner.Dy()-1; i++ {
		buf.SetString(
			fmt.Sprintf("%.2f", float64(i)*verticalScale*(yAxisLabelsGap+1)),
			ui.NewStyle(ui.ColorWhite),
			image.Pt(plt.Inner.Min.X, plt.Inner.Max.Y-(i*(yAxisLabelsGap+1))-2),
		)
	}
}

func (plt *Plot) Draw(buf *ui.Buffer) {
	plt.Block.Draw(buf)

	maxVal := plt.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64From2dSlice(plt.Data)
	}

	if plt.ShowAxes {
		plt.plotAxes(buf, maxVal)
	}

	drawArea := plt.Inner
	if plt.ShowAxes {
		drawArea = image.Rect(
			plt.Inner.Min.X+yAxisLabelsWidth+1, plt.Inner.Min.Y,
			plt.Inner.Max.X, plt.Inner.Max.Y-xAxisLabelsHeight-1,
		)
	}

	switch plt.Marker {
	case MarkerBraille:
		plt.renderBraille(buf, drawArea, maxVal)
	case MarkerDot:
		plt.renderDot(buf, drawArea, maxVal)
	}
}
