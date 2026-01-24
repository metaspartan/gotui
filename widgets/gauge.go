package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui/v5"
)

// Gauge represents a widget that displays a progress bar.
type Gauge struct {
	ui.Block
	Percent    int
	BarColor   ui.Color
	Label      string
	LabelStyle ui.Style
	Gradient   ui.Gradient
}

// NewGauge returns a new Gauge.
func NewGauge() *Gauge {
	return &Gauge{
		Block:      *ui.NewBlock(),
		BarColor:   ui.Theme.Gauge.Bar,
		LabelStyle: ui.Theme.Gauge.Label,
	}
}

// Draw draws the gauge to the buffer.
func (g *Gauge) Draw(buf *ui.Buffer) {
	g.Block.Draw(buf)

	label := g.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", g.Percent)
	}

	barWidth := int((float64(g.Percent) / 100) * float64(g.Inner.Dx()))

	g.drawBar(buf, barWidth)
	g.drawLabel(buf, label, barWidth)
}

func (g *Gauge) drawBar(buf *ui.Buffer, barWidth int) {
	if g.Gradient.Enabled {
		var gradientColors []ui.Color
		if g.Gradient.Direction == 1 {
			gradientColors = ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, g.Inner.Dy())
		} else if barWidth > 0 {
			gradientColors = ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, barWidth)
		}

		for i := range barWidth {
			for y := g.Inner.Min.Y; y < g.Inner.Max.Y; y++ {
				color := ui.ColorClear
				if g.Gradient.Direction == 1 {
					relativeY := y - g.Inner.Min.Y
					if relativeY < len(gradientColors) {
						color = gradientColors[relativeY]
					}
				} else {
					if i < len(gradientColors) {
						color = gradientColors[i]
					}
				}
				buf.SetCell(
					ui.NewCell(' ', ui.NewStyle(ui.ColorClear, color)),
					image.Pt(g.Inner.Min.X+i, y),
				)
			}
		}
	} else {
		buf.Fill(
			ui.NewCell(' ', ui.NewStyle(ui.ColorClear, g.BarColor)),
			image.Rect(g.Inner.Min.X, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
		)
	}
}

func (g *Gauge) drawLabel(buf *ui.Buffer, label string, barWidth int) {
	labelXCoordinate := g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)

	if labelYCoordinate >= g.Inner.Max.Y {
		return
	}

	var gradientColors []ui.Color
	if g.Gradient.Enabled {
		if g.Gradient.Direction == 1 {
			gradientColors = ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, g.Inner.Dy())
		} else {
			gradientColors = ui.GenerateGradient(g.Gradient.Start, g.Gradient.End, barWidth)
		}
	}

	for i, char := range label {
		style := g.LabelStyle
		barX := labelXCoordinate + i - g.Inner.Min.X

		if barX >= 0 && barX < barWidth {
			if g.Gradient.Enabled {
				if g.Gradient.Direction == 1 {
					relativeY := labelYCoordinate - g.Inner.Min.Y
					if relativeY >= 0 && relativeY < len(gradientColors) {
						style = ui.NewStyle(ui.ColorWhite, gradientColors[relativeY])
					}
				} else if barX < len(gradientColors) {
					style = ui.NewStyle(ui.ColorWhite, gradientColors[barX])
				}
			} else {
				style = ui.NewStyle(g.BarColor, ui.ColorClear, ui.ModifierReverse)
			}
		}
		buf.SetCell(ui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
	}
}
