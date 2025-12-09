package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui"
)

type LineGauge struct {
	ui.Block
	Percent        int
	LineColor      ui.Color
	Label          string
	LabelStyle     ui.Style
	LabelAlignment ui.Alignment
}

func NewLineGauge() *LineGauge {
	return &LineGauge{
		Block:          *ui.NewBlock(),
		LineColor:      ui.Theme.Gauge.Bar,
		LabelStyle:     ui.Theme.Gauge.Label,
		LabelAlignment: ui.AlignCenter,
	}
}

func (g *LineGauge) Draw(buf *ui.Buffer) {
	g.Block.Draw(buf)

	label := g.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", g.Percent)
	}

	// Calculate filled width
	barWidth := int((float64(g.Percent) / 100) * float64(g.Inner.Dx()))

	// Center the line vertically
	y := g.Inner.Min.Y + (g.Inner.Dy() / 2)

	for i := 0; i < g.Inner.Dx(); i++ {
		x := g.Inner.Min.X + i
		if x >= g.Inner.Max.X {
			break
		}

		var char rune = '─' // Empty part
		style := g.LabelStyle
		style.Fg = ui.ColorWhite // Default empty color

		if i < barWidth {
			char = '━' // Filled part
			style.Fg = g.LineColor
		}

		if y < g.Inner.Max.Y {
			buf.SetCell(ui.NewCell(char, style), image.Pt(x, y))
		}
	}

	// Plot label
	if y < g.Inner.Max.Y {
		// Check if label is just whitespace
		isWhitespace := true
		for _, r := range label {
			if r != ' ' {
				isWhitespace = false
				break
			}
		}

		if !isWhitespace {
			labelLen := len(label)
			var labelXCoordinate int

			switch g.LabelAlignment {
			case ui.AlignLeft:
				labelXCoordinate = g.Inner.Min.X
			case ui.AlignRight:
				labelXCoordinate = g.Inner.Max.X - labelLen
			default: // Center
				labelXCoordinate = g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(labelLen)/2)
			}

			for i, char := range label {
				style := g.LabelStyle
				x := labelXCoordinate + i
				if x < g.Inner.Max.X && x >= g.Inner.Min.X {
					buf.SetCell(ui.NewCell(char, style), image.Pt(x, y))
				}
			}
		}
	}
}
