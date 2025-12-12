package gotui

import (
	"image"

	rw "github.com/mattn/go-runewidth"
)

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

func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,
		Cells:     make([]Cell, r.Dx()*r.Dy()),
	}
	buf.Fill(CellClear, r)
	return buf
}

func (b *Buffer) GetCell(p image.Point) Cell {
	if !p.In(b.Rectangle) {
		return CellClear
	}
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
	rect = rect.Intersect(b.Rectangle)
	if rect.Empty() {
		return
	}

	width := b.Dx()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		rowStart := (y - b.Min.Y) * width
		startIdx := rowStart + (rect.Min.X - b.Min.X)
		endIdx := rowStart + (rect.Max.X - b.Min.X)

		for i := startIdx; i < endIdx; i++ {
			b.Cells[i] = c
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
