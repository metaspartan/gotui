package gotui

import (
	"image"
	"io"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/metaspartan/gotui/v4/drawille"
)

// ---------- STRUCTS & INTERFACES ----------

type Backend struct {
	Screen         tcell.Screen
	ScreenshotMode bool
}

type ttyAdapter struct {
	rw            io.ReadWriter
	width, height int
	resizeCB      func()
}

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

type TTYHandle interface {
	io.ReadWriter
}

type InitConfig struct {
	CustomTTY      TTYHandle
	Width, Height  int
	SimulationMode bool
	SimulationSize image.Point
}

// ---------- EVENTS ----------

type EventType uint

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

// ---------- STYLE & COLORS ----------

type Color = tcell.Color

type Modifier = tcell.AttrMask

type Style struct {
	Fg       Color
	Bg       Color
	Modifier Modifier
}

type Gradient struct {
	Enabled   bool
	Start     Color
	End       Color
	Direction int // 0 = Horizontal, 1 = Vertical
}

// ---------- BUFFER & CELLS ----------

type Cell struct {
	Rune  rune
	Style Style
}

type Buffer struct {
	image.Rectangle
	Cells []Cell
}

// ---------- LAYOUT -----------

type Alignment uint

// ---------- GRID ------------

type gridItemType uint

type Grid struct {
	Block
	Items []*GridItem
}

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

// ---------- BLOCK ------------

type Block struct {
	Border          bool
	BorderStyle     Style
	BackgroundColor Color
	FillBorder      bool

	BorderLeft, BorderRight, BorderTop, BorderBottom bool

	BorderCollapse bool
	BorderRounded  bool
	BorderType     BorderType

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int

	image.Rectangle
	Inner image.Rectangle

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

	BorderGradient Gradient

	sync.Mutex
}

// ---------- CANVAS -----------

type Canvas struct {
	Block
	drawille.Canvas
}

// ---------- THEMES -----------

type RootTheme struct {
	Default Style

	Block BlockTheme

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

type TableTheme struct {
	Text Style
}
