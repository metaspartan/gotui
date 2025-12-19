package widgets

import (
	"image"

	gotui "github.com/metaspartan/gotui/v4"
	ui "github.com/metaspartan/gotui/v4"
)

// Logo represents a widget that displays the gotui logo.
type Logo struct {
	gotui.Block
	Gradient gotui.Gradient
	image    image.Image
}

// NewLogo returns a new Logo.
func NewLogo() *Logo {
	return &Logo{
		Block: *gotui.NewBlock(),
		Gradient: gotui.Gradient{
			Enabled: false,
			Start:   gotui.NewRGBColor(100, 100, 255),
			End:     gotui.NewRGBColor(255, 100, 200),
		},
	}
}

// Draw draws the logo to the buffer.
func (l *Logo) Draw(buf *gotui.Buffer) {
	l.Block.Draw(buf)
	logoDefinition := []string{
		" ██████   ██████  ████████ ██    ██ ██ ",
		"██       ██    ██    ██    ██    ██ ██ ",
		"██   ███ ██    ██    ██    ██    ██ ██ ",
		"██    ██ ██    ██    ██    ██    ██ ██ ",
		" ██████   ██████     ██     ██████  ██ ",
	}
	logoWidth := len([]rune(logoDefinition[0]))
	logoHeight := len(logoDefinition)
	xStart := l.Inner.Min.X + (l.Inner.Dx()-logoWidth)/2
	yStart := l.Inner.Min.Y + (l.Inner.Dy()-logoHeight)/2
	var gradientColors []ui.Color
	if l.Gradient.Enabled {
		if l.Gradient.Direction == 1 {
			gradientColors = ui.GenerateGradient(l.Gradient.Start, l.Gradient.End, logoHeight)
		} else {
			gradientColors = ui.GenerateGradient(l.Gradient.Start, l.Gradient.End, logoWidth)
		}
	}
	for r, line := range logoDefinition {
		y := yStart + r
		if y >= l.Inner.Max.Y {
			break
		}
		if y < l.Inner.Min.Y {
			continue
		}
		for c, char := range []rune(line) {
			x := xStart + c
			if x >= l.Inner.Max.X {
				break
			}
			if x < l.Inner.Min.X {
				continue
			}
			if char != ' ' {
				style := ui.NewStyle(ui.Theme.Gauge.Bar)
				if l.Gradient.Enabled {
					if l.Gradient.Direction == 1 {
						if r < len(gradientColors) {
							style = ui.NewStyle(gradientColors[r])
						}
					} else if c < len(gradientColors) {
						style = ui.NewStyle(gradientColors[c])
					}
				}
				buf.SetCell(ui.NewCell(char, style), image.Pt(x, y))
			}
		}
	}
}
