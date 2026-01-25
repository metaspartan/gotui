package gotui

import (
	"strings"

	"github.com/gdamore/tcell/v3"
)

const (
	tokenFg              = "fg"
	tokenBg              = "bg"
	tokenModifier        = "mod"
	tokenItemSeparator   = ","
	tokenValueSeparator  = ":"
	tokenBeginStyledText = '['
	tokenEndStyledText   = ']'
	tokenBeginStyle      = '('
	tokenEndStyle        = ')'
)

type parserState uint

const (
	parserStateDefault parserState = iota
	parserStateStyleItems
	parserStateStyledText
)

var StyleParserColorMap = map[string]Color{
	"red":        ColorRed,
	"blue":       ColorBlue,
	"black":      ColorBlack,
	"cyan":       ColorLightCyan,
	"yellow":     ColorYellow,
	"white":      ColorWhite,
	"clear":      ColorClear,
	"green":      ColorGreen,
	"magenta":    ColorMagenta,
	"grey":       ColorGrey,
	"darkgrey":   ColorDarkGrey,
	"lightgrey":  ColorLightGrey,
	"silver":     ColorSilver,
	"orange":     ColorOrange,
	"purple":     ColorPurple,
	"pink":       ColorPink,
	"coral":      ColorCoral,
	"crimson":    ColorCrimson,
	"gold":       ColorGold,
	"teal":       ColorTeal,
	"turquoise":  ColorTurquoise,
	"indigo":     ColorIndigo,
	"violet":     ColorViolet,
	"olive":      ColorOlive,
	"navy":       ColorNavy,
	"aliceblue":  ColorAliceBlue,
	"beige":      ColorBeige,
	"brown":      ColorBrown,
	"darkblue":   ColorDarkBlue,
	"darkcyan":   ColorDarkCyan,
	"darkgreen":  ColorDarkGreen,
	"darkred":    ColorDarkRed,
	"hotpink":    ColorHotPink,
	"lightblue":  ColorLightBlue,
	"lightgreen": ColorLightGreen,
	"lime":       ColorLime,
	"maroon":     ColorMaroon,
	"mintcream":  ColorMintCream,
	"mistyrose":  ColorMistyRose,
	"orchid":     ColorOrchid,
	"plum":       ColorPlum,
	"salmon":     ColorSalmon,
	"seagreen":   ColorSeaGreen,
	"skyblue":    ColorSkyBlue,
	"slateblue":  ColorSlateBlue,
	"tan":        ColorTan,
	"tomato":     ColorTomato,
	"wheat":      ColorWheat,
}
var modifierMap = map[string]Modifier{
	"bold":    ModifierBold,
	"reverse": ModifierReverse,
	"dim":     ModifierDim,
	"blink":   ModifierBlink,
	"italic":  ModifierItalic,
	"strike":  ModifierStrike,
	// NOTE: "underline" was removed in tcell v3
}

func readStyle(runes []rune, defaultStyle Style) Style {
	style := defaultStyle
	split := strings.SplitSeq(string(runes), tokenItemSeparator)
	for item := range split {
		pair := strings.Split(item, tokenValueSeparator)
		if len(pair) == 2 {
			switch pair[0] {
			case tokenFg:
				if c, ok := StyleParserColorMap[pair[1]]; ok {
					style.Fg = c
				} else {
					style.Fg = tcell.GetColor(pair[1])
				}
			case tokenBg:
				if c, ok := StyleParserColorMap[pair[1]]; ok {
					style.Bg = c
				} else {
					style.Bg = tcell.GetColor(pair[1])
				}
			case tokenModifier:
				style.Modifier = modifierMap[pair[1]]
			}
		}
	}
	return style
}
func ParseStyles(s string, defaultStyle Style) []Cell {
	cells := []Cell{}
	runes := []rune(s)
	state := parserStateDefault
	styledText := []rune{}
	styleItems := []rune{}
	squareCount := 0
	reset := func() {
		styledText = []rune{}
		styleItems = []rune{}
		state = parserStateDefault
		squareCount = 0
	}
	rollback := func() {
		cells = append(cells, RunesToStyledCells(styledText, defaultStyle)...)
		cells = append(cells, RunesToStyledCells(styleItems, defaultStyle)...)
		reset()
	}
	chop := func(s []rune) []rune {
		return s[1 : len(s)-1]
	}
	for i, _rune := range runes {
		switch state {
		case parserStateDefault:
			if _rune == tokenBeginStyledText {
				state = parserStateStyledText
				squareCount = 1
				styledText = append(styledText, _rune)
			} else {
				cells = append(cells, Cell{_rune, defaultStyle})
			}
		case parserStateStyledText:
			switch {
			case squareCount == 0:
				switch _rune {
				case tokenBeginStyle:
					state = parserStateStyleItems
					styleItems = append(styleItems, _rune)
				default:
					rollback()
					switch _rune {
					case tokenBeginStyledText:
						state = parserStateStyledText
						squareCount = 1
						styleItems = append(styleItems, _rune)
					default:
						cells = append(cells, Cell{_rune, defaultStyle})
					}
				}
			case len(runes) == i+1:
				rollback()
				styledText = append(styledText, _rune)
			case _rune == tokenBeginStyledText:
				squareCount++
				styledText = append(styledText, _rune)
			case _rune == tokenEndStyledText:
				squareCount--
				styledText = append(styledText, _rune)
			default:
				styledText = append(styledText, _rune)
			}
		case parserStateStyleItems:
			styleItems = append(styleItems, _rune)
			if _rune == tokenEndStyle {
				style := readStyle(chop(styleItems), defaultStyle)
				cells = append(cells, RunesToStyledCells(chop(styledText), style)...)
				reset()
			} else if len(runes) == i+1 {
				rollback()
			}
		}
	}
	return cells
}
