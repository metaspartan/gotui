package gotui

import (
	"image"

	rw "github.com/mattn/go-runewidth"
)

// Cell represents a viewable terminal cell
type Cell struct {
	Rune  rune
	Style Style
}

var CellClear = Cell{
	Rune:  ' ',
	Style: StyleClear,
}

// NewCell takes 1 to 2 arguments
// 1st argument = rune
// 2nd argument = optional style
func NewCell(rune rune, args ...interface{}) Cell {
	style := StyleClear
	if len(args) == 1 {
		style = args[0].(Style)
	}
	return Cell{
		Rune:  rune,
		Style: style,
	}
}

// Buffer represents a section of a terminal and is a renderable rectangle of cells.
type Buffer struct {
	image.Rectangle
	CellMap map[image.Point]Cell
}

func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,
		CellMap:   make(map[image.Point]Cell),
	}
	buf.Fill(CellClear, r) // clears out area
	return buf
}

func (b *Buffer) GetCell(p image.Point) Cell {
	return b.CellMap[p]
}

func (b *Buffer) SetCell(c Cell, p image.Point) {
	b.CellMap[p] = c
}

func (b *Buffer) Fill(c Cell, rect image.Rectangle) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			b.SetCell(c, image.Pt(x, y))
		}
	}
}

func (b *Buffer) SetString(s string, style Style, p image.Point) {
	runes := []rune(s)
	x := 0
	for _, char := range runes {
		b.SetCell(Cell{char, style}, image.Pt(p.X+x, p.Y))
		x += rw.RuneWidth(char)
	}
}
