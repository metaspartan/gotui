package gotui

import "github.com/gdamore/tcell/v2"

// ---------- GLOBALS ----------

var (
	DefaultBackend = &Backend{}
)

var Screen tcell.Screen

var ScreenshotMode bool

var StandardColors = []Color{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorMagenta,
	ColorLightCyan,
	ColorWhite,
}

var StandardStyles = []Style{
	NewStyle(ColorRed),
	NewStyle(ColorGreen),
	NewStyle(ColorYellow),
	NewStyle(ColorBlue),
	NewStyle(ColorMagenta),
	NewStyle(ColorLightCyan),
	NewStyle(ColorWhite),
}

var Theme = RootTheme{
	Default: NewStyle(ColorWhite),

	Block: BlockTheme{
		Title:  NewStyle(ColorWhite),
		Border: NewStyle(ColorWhite),
	},

	BarChart: BarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Paragraph: ParagraphTheme{
		Text: NewStyle(ColorWhite),
	},

	PieChart: PieChartTheme{
		Slices: StandardColors,
	},

	List: ListTheme{
		Text: NewStyle(ColorWhite),
	},

	Tree: TreeTheme{
		Text:      NewStyle(ColorWhite),
		Collapsed: COLLAPSED,
		Expanded:  EXPANDED,
	},

	StackedBarChart: StackedBarChartTheme{
		Bars:   StandardColors,
		Nums:   StandardStyles,
		Labels: StandardStyles,
	},

	Gauge: GaugeTheme{
		Bar:   ColorWhite,
		Label: NewStyle(ColorWhite),
	},

	Sparkline: SparklineTheme{
		Title: NewStyle(ColorWhite),
		Line:  ColorWhite,
	},

	Plot: PlotTheme{
		Lines: StandardColors,
		Axes:  ColorWhite,
	},

	Table: TableTheme{
		Text: NewStyle(ColorWhite),
	},

	Tab: TabTheme{
		Active:   NewStyle(ColorRed),
		Inactive: NewStyle(ColorWhite),
	},
}

var StyleClear = Style{
	Fg:       ColorClear,
	Bg:       ColorClear,
	Modifier: ModifierClear,
}

var CellClear = Cell{
	Rune:  ' ',
	Style: StyleClear,
}

// ---------- CONSTANTS ----------

const (
	KeyboardEvent EventType = iota
	MouseEvent
	ResizeEvent
)

const ColorClear Color = tcell.ColorDefault

const (
	ColorBlack     Color = tcell.ColorBlack
	ColorRed       Color = tcell.ColorRed
	ColorGreen     Color = tcell.ColorGreen
	ColorYellow    Color = tcell.ColorYellow
	ColorBlue      Color = tcell.ColorBlue
	ColorMagenta   Color = tcell.ColorDarkMagenta
	ColorLightCyan Color = tcell.ColorLightCyan
	ColorWhite     Color = tcell.ColorWhite

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

const (
	ModifierClear     Modifier = 0
	ModifierBold      Modifier = tcell.AttrBold
	ModifierUnderline Modifier = tcell.AttrUnderline
	ModifierReverse   Modifier = tcell.AttrReverse
)

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

const (
	AlignTop    Alignment = 0
	AlignMiddle Alignment = 1
	AlignBottom Alignment = 2
)

const (
	col gridItemType = 0
	row gridItemType = 1
)
