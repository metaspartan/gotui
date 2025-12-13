package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Logo widget renders a hardcoded logo (Gotui)
type Logo struct {
	ui.Block
	Gradient      bool
	GradientStart ui.Color
	GradientEnd   ui.Color
}

func NewLogo() *Logo {
	return &Logo{
		Block:         *ui.NewBlock(),
		Gradient:      true,
		GradientStart: ui.NewRGBColor(100, 100, 255), // Default Blue-ish
		GradientEnd:   ui.NewRGBColor(255, 100, 200), // Default Pink-ish
	}
}

func (l *Logo) Draw(buf *ui.Buffer) {
	l.Block.Draw(buf)

	// ASCII Art / Block Art for "GOTUI"
	// 5 chars high
	// 5x3 Grid per letter approx?
	// Let's define the specific "GOTUI" shape using a coordinate map or simple bool grid.
	// 5 chars: G, O, T, U, I

	// Custom 5-pixel high font for GOTUI
	// G:
	//  ####
	// #
	// #  ##
	// #   #
	//  ###

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
	if l.Gradient {
		gradientColors = ui.GenerateGradient(l.GradientStart, l.GradientEnd, logoWidth)
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
				if l.Gradient && c < len(gradientColors) {
					style = ui.NewStyle(gradientColors[c])
				}
				buf.SetCell(ui.NewCell(char, style), image.Pt(x, y))
			}
		}
	}
}
