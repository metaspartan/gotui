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

// GenerateGradient generates a gradient from start to end color.
func GenerateGradient(start, end Color, length int) []Color {
	if length <= 0 {
		return []Color{}
	}
	colors := make([]Color, length)
	for i := 0; i < length; i++ {
		colors[i] = InterpolateColor(start, end, i, length)
	}
	return colors
}

// ApplyGradientToText applies a gradient to a string of text.
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

// HexToColor converts a hex color string to a Color.
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

// GenerateMultiGradient generates a gradient that transitions through multiple colors.
func GenerateMultiGradient(length int, colors ...Color) []Color {
	if length <= 0 || len(colors) == 0 {
		return []Color{}
	}
	if len(colors) == 1 || length == 1 {
		res := make([]Color, length)
		for i := range res {
			res[i] = colors[0]
		}
		return res
	}

	result := make([]Color, length)
	segments := len(colors) - 1

	for i := 0; i < length; i++ {
		t := float64(i) / float64(length-1) * float64(segments)
		segmentIdx := int(t)
		if segmentIdx >= segments {
			segmentIdx = segments - 1
		}
		segmentT := t - float64(segmentIdx)

		startColor := colors[segmentIdx]
		endColor := colors[segmentIdx+1]

		r1, g1, b1 := startColor.RGB()
		r2, g2, b2 := endColor.RGB()

		r := int32(float64(r1) + segmentT*float64(r2-r1))
		g := int32(float64(g1) + segmentT*float64(g2-g1))
		b := int32(float64(b1) + segmentT*float64(b2-b1))

		result[i] = NewRGBColor(r, g, b)
	}

	return result
}
