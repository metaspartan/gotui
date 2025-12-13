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
	Gradient   ui.Gradient
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

	barWidth := int((float64(g.Percent) / 100) * float64(g.Inner.Dx()))

	if g.Gradient.Enabled {
		colors := ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, barWidth)
		for i := 0; i < barWidth; i++ {
			if i < len(colors) {
				for y := g.Inner.Min.Y; y < g.Inner.Max.Y; y++ {
					buf.SetCell(
						ui.NewCell(' ', ui.NewStyle(ui.ColorClear, colors[i])),
						image.Pt(g.Inner.Min.X+i, y),
					)
				}
			}
		}
	} else {
		buf.Fill(
			ui.NewCell(' ', ui.NewStyle(ui.ColorClear, g.BarColor)),
			image.Rect(g.Inner.Min.X, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
		)
	}

	labelXCoordinate := g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)
	if labelYCoordinate < g.Inner.Max.Y {
		var gradientColors []ui.Color
		if g.Gradient.Enabled {
			gradientColors = ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, barWidth)
		}

		for i, char := range label {
			style := g.LabelStyle
			barX := labelXCoordinate + i - g.Inner.Min.X
			if barX >= 0 && barX < barWidth {
				if g.Gradient.Enabled && barX < len(gradientColors) {
					style = ui.NewStyle(ui.ColorWhite, gradientColors[barX])
				} else {
					style = ui.NewStyle(g.BarColor, ui.ColorClear, ui.ModifierReverse)
				}
			}
			buf.SetCell(ui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
