package gotui

import "github.com/gdamore/tcell/v2"

// Color is an integer from -1 to 255
// -1 = ColorClear
// 0-255 = Xterm colors
type Color = tcell.Color

// ColorClear clears the Fg or Bg color of a Style
const ColorClear Color = tcell.ColorDefault

// Basic terminal colors
const (
	// Standard Colors
	ColorBlack   Color = tcell.ColorBlack
	ColorRed     Color = tcell.ColorRed
	ColorGreen   Color = tcell.ColorGreen
	ColorYellow  Color = tcell.ColorYellow
	ColorBlue    Color = tcell.ColorBlue
	ColorMagenta Color = tcell.ColorDarkMagenta // Termui legacy mapping
	ColorCyan    Color = tcell.ColorLightCyan   // Termui legacy mapping
	ColorWhite   Color = tcell.ColorWhite

	// Extended Colors (Common)
	ColorGrey       Color = tcell.ColorGrey
	ColorDarkGrey   Color = tcell.ColorDarkGrey
	ColorLightGrey  Color = tcell.ColorLightGrey
	ColorSilver     Color = tcell.ColorSilver
	ColorOrange     Color = tcell.ColorOrange
	ColorPurple     Color = tcell.ColorPurple
	ColorPink       Color = tcell.ColorPink
	ColorCoral      Color = tcell.ColorCoral
	ColorCrimson    Color = tcell.ColorCrimson
	ColorGold       Color = tcell.ColorGold
	ColorTeal       Color = tcell.ColorTeal
	ColorTurquoise  Color = tcell.ColorTurquoise
	ColorIndigo     Color = tcell.ColorIndigo
	ColorViolet     Color = tcell.ColorViolet
	ColorOlive      Color = tcell.ColorOlive
	ColorNavy       Color = tcell.ColorNavy
	ColorAliceBlue  Color = tcell.ColorAliceBlue
	ColorBeige      Color = tcell.ColorBeige
	ColorBrown      Color = tcell.ColorBrown
	ColorDarkBlue   Color = tcell.ColorDarkBlue
	ColorDarkCyan   Color = tcell.ColorDarkCyan
	ColorDarkGreen  Color = tcell.ColorDarkGreen
	ColorDarkRed    Color = tcell.ColorDarkRed
	ColorHotPink    Color = tcell.ColorHotPink
	ColorLightBlue  Color = tcell.ColorLightBlue
	ColorLightGreen Color = tcell.ColorLightGreen
	ColorLime       Color = tcell.ColorLime
	ColorMaroon     Color = tcell.ColorMaroon
	ColorMintCream  Color = tcell.ColorMintCream
	ColorMistyRose  Color = tcell.ColorMistyRose
	ColorOrchid     Color = tcell.ColorOrchid
	ColorPlum       Color = tcell.ColorPlum
	ColorSalmon     Color = tcell.ColorSalmon
	ColorSeaGreen   Color = tcell.ColorSeaGreen
	ColorSkyBlue    Color = tcell.ColorSkyblue
	ColorSlateBlue  Color = tcell.ColorSlateBlue
	ColorTan        Color = tcell.ColorTan
	ColorTomato     Color = tcell.ColorTomato
	ColorWheat      Color = tcell.ColorWheat
)

type Modifier = tcell.AttrMask

const (
	// ModifierClear clears any modifiers
	ModifierClear     Modifier = 0
	ModifierBold      Modifier = tcell.AttrBold
	ModifierUnderline Modifier = tcell.AttrUnderline
	ModifierReverse   Modifier = tcell.AttrReverse
)

// Style represents the style of one terminal cell
type Style struct {
	Fg       Color
	Bg       Color
	Modifier Modifier
}

// StyleClear represents a default Style, with no colors or modifiers
var StyleClear = Style{
	Fg:       ColorClear,
	Bg:       ColorClear,
	Modifier: ModifierClear,
}

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
