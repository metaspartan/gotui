package widgets

import (
	"fmt"
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui/v4"
)

type BarChart struct {
	ui.Block
	BarColors    []ui.Color
	LabelStyles  []ui.Style
	NumStyles    []ui.Style // only Fg and Modifier are used
	NumFormatter func(float64) string
	Data         []float64
	Labels       []string
	BarWidth     int
	BarGap       int
	MaxVal       float64
}

func NewBarChart() *BarChart {
	return &BarChart{
		Block:        *ui.NewBlock(),
		BarColors:    ui.Theme.BarChart.Bars,
		NumStyles:    ui.Theme.BarChart.Nums,
		LabelStyles:  ui.Theme.BarChart.Labels,
		NumFormatter: func(n float64) string { return fmt.Sprint(n) },
		BarGap:       1,
		BarWidth:     3,
	}
}

func (bc *BarChart) Draw(buf *ui.Buffer) {
	bc.Block.Draw(buf)

	maxVal := bc.MaxVal
	if maxVal == 0 {
		maxVal, _ = ui.GetMaxFloat64FromSlice(bc.Data)
	}

	barXCoordinate := bc.Inner.Min.X

	for i, data := range bc.Data {
		if data > 0 {
			// draw bar
			height := int((data / maxVal) * float64(bc.Inner.Dy()-1))
			for x := barXCoordinate; x < ui.MinInt(barXCoordinate+bc.BarWidth, bc.Inner.Max.X); x++ {
				for y := bc.Inner.Max.Y - 2; y > (bc.Inner.Max.Y-2)-height; y-- {
					c := ui.NewCell(' ', ui.NewStyle(ui.ColorClear, ui.SelectColor(bc.BarColors, i)))
					buf.SetCell(c, image.Pt(x, y))
				}
			}
		}
		// draw label
		if i < len(bc.Labels) {
			labelXCoordinate := barXCoordinate +
				int((float64(bc.BarWidth) / 2)) -
				int((float64(rw.StringWidth(bc.Labels[i])) / 2))
			buf.SetString(
				bc.Labels[i],
				ui.SelectStyle(bc.LabelStyles, i),
				image.Pt(labelXCoordinate, bc.Inner.Max.Y-1),
			)
		}

		// draw number
		numberXCoordinate := barXCoordinate + int((float64(bc.BarWidth) / 2))
		if numberXCoordinate <= bc.Inner.Max.X {
			buf.SetString(
				bc.NumFormatter(data),
				ui.NewStyle(
					ui.SelectStyle(bc.NumStyles, i+1).Fg,
					ui.SelectColor(bc.BarColors, i),
					ui.SelectStyle(bc.NumStyles, i+1).Modifier,
				),
				image.Pt(numberXCoordinate, bc.Inner.Max.Y-2),
			)
		}

		barXCoordinate += (bc.BarWidth + bc.BarGap)
	}
}
