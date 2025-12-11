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
	Cells []Cell
}

func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,
		Cells:     make([]Cell, r.Dx()*r.Dy()),
	}
	buf.Fill(CellClear, r) // clears out area
	return buf
}

func (b *Buffer) GetCell(p image.Point) Cell {
	if !p.In(b.Rectangle) {
		return CellClear
	}
	// Index calculation: (y - Min.Y) * width + (x - Min.X)
	idx := (p.Y-b.Min.Y)*b.Dx() + (p.X - b.Min.X)
	if idx >= 0 && idx < len(b.Cells) {
		return b.Cells[idx]
	}
	return CellClear
}

func (b *Buffer) SetCell(c Cell, p image.Point) {
	if !p.In(b.Rectangle) {
		return
	}
	idx := (p.Y-b.Min.Y)*b.Dx() + (p.X - b.Min.X)
	if idx >= 0 && idx < len(b.Cells) {
		b.Cells[idx] = c
	}
}

func (b *Buffer) Fill(c Cell, rect image.Rectangle) {
	// Intersect the fill rect with the buffer bounds
	rect = rect.Intersect(b.Rectangle)
	if rect.Empty() {
		return
	}

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			b.SetCell(c, image.Pt(x, y))
		}
	}
}

func (b *Buffer) SetString(s string, style Style, p image.Point) {
	x := 0
	for _, char := range s {
		b.SetCell(Cell{char, style}, image.Pt(p.X+x, p.Y))
		x += rw.RuneWidth(char)
	}
}
