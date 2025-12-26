package gotui

import (
	"image"

	"github.com/metaspartan/gotui/v5/drawille"
)

// NewCanvas returns a new Canvas.
func NewCanvas() *Canvas {
	return &Canvas{
		Block:  *NewBlock(),
		Canvas: *drawille.NewCanvas(),
	}
}

// SetPoint sets the color of the point at the given coordinates.
func (c *Canvas) SetPoint(p image.Point, color Color) {
	c.Canvas.SetPoint(p, drawille.Color(color))
}

// SetLine draws a line from p0 to p1 with the given color.
func (c *Canvas) SetLine(p0, p1 image.Point, color Color) {
	c.Canvas.SetLine(p0, p1, drawille.Color(color))
}

// Draw draws the canvas to the buffer.
func (c *Canvas) Draw(buf *Buffer) {
	c.Block.Draw(buf)
	for point, cell := range c.Canvas.GetCells() {
		dest := point.Add(c.Inner.Min)
		if dest.In(c.Inner) {
			col := Color(cell.Color)
			if col == 0 || col == ColorClear {
				col = ColorWhite
			}
			convertedCell := Cell{
				cell.Rune,
				Style{
					col,
					ColorClear,
					ModifierClear,
				},
			}
			buf.SetCell(convertedCell, dest)
		}
	}
}
