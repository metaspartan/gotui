package gotui

import (
	"image"
	"io"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/metaspartan/gotui/v4/drawille"
)

// Backend represents the backend that the library uses to draw to the screen.
type Backend struct {
	Screen         tcell.Screen
	ScreenshotMode bool
}

type ttyAdapter struct {
	rw            io.ReadWriter
	width, height int
	resizeCB      func()
}

// Drawable represents a widget that can be drawn to the buffer.
type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

// EventHandler represents a widget that can handle events.
type EventHandler interface {
	HandleEvent(Event) bool
}

// Widget represents a Drawable that also handles events.
type Widget interface {
	Drawable
	EventHandler
}

// TTYHandle represents a handle to a TTY.
type TTYHandle interface {
	io.ReadWriter
}

// InitConfig represents the configuration for initializing the library.
type InitConfig struct {
	CustomTTY      TTYHandle
	Width, Height  int
	SimulationMode bool
	SimulationSize image.Point
}

// EventType represents the type of an event.
type EventType uint

// Event represents an event that occurred.
type Event struct {
	Type    EventType
	ID      string
	Payload interface{}
}
type Mouse struct {
	Drag bool
	X    int
	Y    int
}
type Resize struct {
	Width  int
	Height int
}
type Color = tcell.Color
type Modifier = tcell.AttrMask

// Style represents the style of a cell.
type Style struct {
	Fg       Color
	Bg       Color
	Modifier Modifier
}

// Gradient represents a color gradient.
type Gradient struct {
	Enabled   bool
	Start     Color
	End       Color
	Stops     []Color
	Direction int
}

// Cell represents a single cell in the terminal.
type Cell struct {
	Rune  rune
	Style Style
}

// Buffer represents a buffer of cells.
type Buffer struct {
	image.Rectangle
	Cells []Cell
}

// Alignment represents the alignment of text.
type Alignment uint

// gridItemType represents the type of a grid item.
type gridItemType uint

// Grid allows you to lay out widgets in a grid.
type Grid struct {
	Block
	Items []*GridItem
}

// GridItem represents an item in the grid.
type GridItem struct {
	Type        gridItemType
	XRatio      float64
	YRatio      float64
	WidthRatio  float64
	HeightRatio float64
	Entry       interface{}
	IsLeaf      bool
	ratio       float64
}

// Block is the base struct for all widgets.
type Block struct {
	Border                                               bool
	BorderStyle                                          Style
	BackgroundColor                                      Color
	FillBorder                                           bool
	BorderLeft, BorderRight, BorderTop, BorderBottom     bool
	BorderCollapse                                       bool
	BorderRounded                                        bool
	BorderType                                           BorderType
	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int
	image.Rectangle
	Inner                image.Rectangle
	Title                string
	TitleLeft            string
	TitleRight           string
	TitleStyle           Style
	TitleAlignment       Alignment
	TitleBottom          string
	TitleBottomLeft      string
	TitleBottomRight     string
	TitleBottomStyle     Style
	TitleBottomAlignment Alignment
	BorderGradient       Gradient
	BorderSet            *BorderSet
	sync.Mutex
}

// Canvas is a widget that allows drawing points and lines using braille characters.
type Canvas struct {
	Block
	drawille.Canvas
}
type RootTheme struct {
	Default         Style
	Block           BlockTheme
	BarChart        BarChartTheme
	Gauge           GaugeTheme
	Plot            PlotTheme
	List            ListTheme
	Tree            TreeTheme
	Paragraph       ParagraphTheme
	PieChart        PieChartTheme
	Sparkline       SparklineTheme
	StackedBarChart StackedBarChartTheme
	Tab             TabTheme
	Table           TableTheme
}
type BlockTheme struct {
	Title  Style
	Border Style
}
type BarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}
type GaugeTheme struct {
	Bar   Color
	Label Style
}
type PlotTheme struct {
	Lines []Color
	Axes  Color
}
type ListTheme struct {
	Text Style
}
type TreeTheme struct {
	Text      Style
	Collapsed rune
	Expanded  rune
}
type ParagraphTheme struct {
	Text Style
}
type PieChartTheme struct {
	Slices []Color
}
type SparklineTheme struct {
	Title Style
	Line  Color
}
type StackedBarChartTheme struct {
	Bars   []Color
	Nums   []Style
	Labels []Style
}
type TabTheme struct {
	Active   Style
	Inactive Style
}

// TableTheme represents the theme for a table.
type TableTheme struct {
	Text Style
}
