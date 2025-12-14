package gotui

import "github.com/gdamore/tcell/v2"

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
func NewColorRGB(r, g, b int32) Color {
	return tcell.NewRGBColor(r, g, b)
}
func NewRGBColor(r, g, b int32) Color {
	return NewColorRGB(r, g, b)
}
