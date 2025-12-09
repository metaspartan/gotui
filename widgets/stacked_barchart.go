package widgets

import (
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui"
)

type StackedBarChart struct {
	ui.Block
	BarColors    []ui.Color
	LabelStyles  []ui.Style
	NumStyles    []ui.Style // only Fg and Modifier are used
	NumFormatter func(float64) string
	Data         [][]float64
	Labels       []string
	BarWidth     int
	BarGap       int
	MaxVal       float64
}

func NewStackedBarChart() *StackedBarChart {
	return &StackedBarChart{
		Block:        *ui.NewBlock(),
		BarColors:    ui.Theme.StackedBarChart.Bars,
		LabelStyles:  ui.Theme.StackedBarChart.Labels,
		NumStyles:    ui.Theme.StackedBarChart.Nums,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       1,
		BarWidth:     3,
	}
}

func (sbc *StackedBarChart) Draw(buf *ui.Buffer) {
	sbc.Block.Draw(buf)

	maxVal := sbc.MaxVal
	if maxVal == 0 {
		for _, data := range sbc.Data {
			maxVal = ui.MaxFloat64(maxVal, ui.SumFloat64Slice(data))
		}
	}

	barXCoordinate := sbc.Inner.Min.X

	for i, bar := range sbc.Data {
		// draw stacked bars
		stackedBarYCoordinate := 0
		for j, data := range bar {
			// draw each stacked bar
			height := int((data / maxVal) * float64(sbc.Inner.Dy()-1))
			for x := barXCoordinate; x < ui.MinInt(barXCoordinate+sbc.BarWidth, sbc.Inner.Max.X); x++ {
				for y := (sbc.Inner.Max.Y - 2) - stackedBarYCoordinate; y > (sbc.Inner.Max.Y-2)-stackedBarYCoordinate-height; y-- {
					c := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(sbc.BarColors, j)))
					buf.SetCell(c, image.Pt(x, y))
				}
			}

			// draw number
			numberXCoordinate := barXCoordinate + int((float64(sbc.BarWidth) / 2)) - 1
			buf.SetString(
				sbc.NumFormatter(data),
				ui.NewStyle(
					ui.SelectStyle(sbc.NumStyles, j+1).Fg,
					ui.SelectColor(sbc.BarColors, j),
					ui.SelectStyle(sbc.NumStyles, j+1).Modifier,
				),
				image.Pt(numberXCoordinate, (sbc.Inner.Max.Y-2)-stackedBarYCoordinate),
			)

			stackedBarYCoordinate += height
		}

		// draw label
		if i < len(sbc.Labels) {
			labelXCoordinate := barXCoordinate + ui.MaxInt(
				int((float64(sbc.BarWidth)/2))-int((float64(rw.StringWidth(sbc.Labels[i]))/2)),
				0,
			)
			buf.SetString(
				ui.TrimString(sbc.Labels[i], sbc.BarWidth),
				ui.SelectStyle(sbc.LabelStyles, i),
				image.Pt(labelXCoordinate, sbc.Inner.Max.Y-1),
			)
		}

		barXCoordinate += (sbc.BarWidth + sbc.BarGap)
	}
}
