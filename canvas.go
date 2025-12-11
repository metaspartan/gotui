package gotui

import (
	"image"

	"github.com/metaspartan/gotui/drawille"
	// "log"
)

type Canvas struct {
	Block
	drawille.Canvas
}

func NewCanvas() *Canvas {
	return &Canvas{
		Block:  *NewBlock(),
		Canvas: *drawille.NewCanvas(),
	}
}

func (c *Canvas) SetPoint(p image.Point, color Color) {
	c.Canvas.SetPoint(p, drawille.Color(color))
}

func (c *Canvas) SetLine(p0, p1 image.Point, color Color) {
	c.Canvas.SetLine(p0, p1, drawille.Color(color))
}

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
