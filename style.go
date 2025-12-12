package gotui

import "github.com/gdamore/tcell/v2"

// NewStyle takes 1 to 3 arguments
// 1st argument = Fg
// 2nd argument = optional Bg
// 3rd argument = optional Modifier
func NewStyle(fg Color, args ...interface{}) Style {
	bg := ColorClear
	modifier := ModifierClear
	if len(args) >= 1 {
		bg = args[0].(Color)
	}
	if len(args) == 2 {
		modifier = args[1].(Modifier)
	}
	return Style{
		fg,
		bg,
		modifier,
	}
}

// NewColorRGB returns a new Color with the given RGB values
func NewColorRGB(r, g, b int32) Color {
	return tcell.NewRGBColor(r, g, b)
}

// NewRGBColor is a convenience alias for NewColorRGB
func NewRGBColor(r, g, b int32) Color {
	return NewColorRGB(r, g, b)
}
