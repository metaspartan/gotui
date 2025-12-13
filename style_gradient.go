package gotui

import (
	"fmt"
)

func InterpolateColor(c1, c2 Color, step, steps int) Color {
	if steps <= 1 {
		return c1
	}
	if step >= steps {
		return c2
	}

	r1, g1, b1 := c1.RGB()
	r2, g2, b2 := c2.RGB()

	factor := float64(step) / float64(steps-1)

	r := int32(float64(r1) + factor*float64(r2-r1))
	g := int32(float64(g1) + factor*float64(g2-g1))
	b := int32(float64(b1) + factor*float64(b2-b1))

	return NewRGBColor(r, g, b)
}

func GenerateGradient(start, end Color, length int) []Color {
	colors := make([]Color, length)
	for i := 0; i < length; i++ {
		colors[i] = InterpolateColor(start, end, i, length)
	}
	return colors
}

func ApplyGradientToText(text string, start, end Color) []Cell {
	runes := []rune(text)
	colors := GenerateGradient(start, end, len(runes))
	cells := make([]Cell, len(runes))
	for i, r := range runes {
		cells[i] = Cell{
			Rune:  r,
			Style: NewStyle(colors[i]),
		}
	}
	return cells
}

func HexToColor(hex string) (Color, error) {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return Color(0), fmt.Errorf("invalid hex color length: %d, expected 6", len(hex))
	}

	var r, g, b int32
	n, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return Color(0), err
	}
	if n != 3 {
		return Color(0), fmt.Errorf("invalid hex color format")
	}

	return NewRGBColor(r, g, b), nil
}
