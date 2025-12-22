package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Logo represents a widget that displays the gotui logo.
type Logo struct {
	ui.Block
	Gradient ui.Gradient
}

// NewLogo returns a new Logo.
func NewLogo() *Logo {
	return &Logo{
		Block: *ui.NewBlock(),
		Gradient: ui.Gradient{
			Enabled: false,
			Start:   ui.NewRGBColor(100, 100, 255),
			End:     ui.NewRGBColor(255, 100, 200),
		},
	}
}

// Draw draws the logo to the buffer.
func (l *Logo) Draw(buf *ui.Buffer) {
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
	gradientColors := l.generateGradientColors(logoWidth, logoHeight)
	l.drawLogoLines(buf, logoDefinition, xStart, yStart, gradientColors)
}

func (l *Logo) generateGradientColors(logoWidth, logoHeight int) []ui.Color {
	if !l.Gradient.Enabled {
		return nil
	}
	length := logoWidth
	if l.Gradient.Direction == 1 {
		length = logoHeight
	}
	if len(l.Gradient.Stops) > 0 {
		return ui.GenerateMultiGradient(length, l.Gradient.Stops...)
	}
	return ui.GenerateGradient(l.Gradient.Start, l.Gradient.End, length)
}

func (l *Logo) drawLogoLines(buf *ui.Buffer, lines []string, xStart, yStart int, gradientColors []ui.Color) {
	for r, line := range lines {
		y := yStart + r
		if y >= l.Inner.Max.Y || y < l.Inner.Min.Y {
			continue
		}
		for c, char := range []rune(line) {
			x := xStart + c
			if x >= l.Inner.Max.X || x < l.Inner.Min.X || char == ' ' {
				continue
			}
			style := l.getCharStyle(r, c, gradientColors)
			buf.SetCell(ui.NewCell(char, style), image.Pt(x, y))
		}
	}
}

func (l *Logo) getCharStyle(r, c int, gradientColors []ui.Color) ui.Style {
	if gradientColors == nil {
		return ui.NewStyle(ui.Theme.Gauge.Bar)
	}
	idx := c
	if l.Gradient.Direction == 1 {
		idx = r
	}
	if idx < len(gradientColors) {
		return ui.NewStyle(gradientColors[idx])
	}
	return ui.NewStyle(ui.Theme.Gauge.Bar)
}
