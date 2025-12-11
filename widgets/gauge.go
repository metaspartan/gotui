package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

type Gauge struct {
	ui.Block
	Percent    int
	BarColor   ui.Color
	Label      string
	LabelStyle ui.Style
}

func NewGauge() *Gauge {
	return &Gauge{
		Block:      *ui.NewBlock(),
		BarColor:   ui.Theme.Gauge.Bar,
		LabelStyle: ui.Theme.Gauge.Label,
	}
}

func (g *Gauge) Draw(buf *ui.Buffer) {
	g.Block.Draw(buf)

	label := g.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", g.Percent)
	}

	// plot bar
	barWidth := int((float64(g.Percent) / 100) * float64(g.Inner.Dx()))
	buf.Fill(
		ui.NewCell(' ', ui.NewStyle(ui.ColorClear, g.BarColor)),
		image.Rect(g.Inner.Min.X, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
	)

	// plot label
	labelXCoordinate := g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)
	if labelYCoordinate < g.Inner.Max.Y {
		for i, char := range label {
			style := g.LabelStyle
			if labelXCoordinate+i+1 <= g.Inner.Min.X+barWidth {
				style = ui.NewStyle(g.BarColor, ui.ColorClear, ui.ModifierReverse)
			}
			buf.SetCell(ui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
