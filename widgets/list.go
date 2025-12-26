package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui/v5"
)

// List represents a widget that displays a list of items.
type List struct {
	ui.Block
	Rows          []string
	WrapText      bool
	TextStyle     ui.Style
	SelectedStyle ui.Style
	TextAlignment ui.Alignment
	SelectedRow   int
	Gradient      ui.Gradient
	topRow        int
}

// NewList returns a new List.
func NewList() *List {
	return &List{
		Block:         *ui.NewBlock(),
		TextStyle:     ui.Theme.List.Text,
		SelectedStyle: ui.Theme.List.Text,
		TextAlignment: ui.AlignLeft,
	}
}

// Draw draws the list to the buffer.
func (l *List) Draw(buf *ui.Buffer) {
	l.Block.Draw(buf)

	if l.SelectedRow >= l.Inner.Dy()+l.topRow {
		l.topRow = l.SelectedRow - l.Inner.Dy() + 1
	} else if l.SelectedRow < l.topRow {
		l.topRow = l.SelectedRow
	}

	l.drawRows(buf)
	l.drawArrows(buf)
}

// drawRows draws the visible rows of the list to the buffer.
func (l *List) drawRows(buf *ui.Buffer) {
	var gradientColors []ui.Color
	if l.Gradient.Enabled && l.Gradient.Direction == 1 {
		gradientColors = ui.GenerateGradient(l.Gradient.Start, l.Gradient.End, l.Inner.Dy())
	}

	y := l.Inner.Min.Y
	for row := l.topRow; row < len(l.Rows) && y < l.Inner.Max.Y; row++ {
		y = l.drawRow(buf, row, y, gradientColors)
	}
}

// getRowCells prepares the cells for a given row, applying styles and gradients.
func (l *List) getRowCells(row int) []ui.Cell {
	var cells []ui.Cell
	if l.Gradient.Enabled && l.Gradient.Direction == 0 {
		cells = ui.ApplyGradientToText(l.Rows[row], l.Gradient.Start, l.Gradient.End)
	} else {
		cells = ui.ParseStyles(l.Rows[row], l.TextStyle)
		if row == l.SelectedRow {
			for i := 0; i < len(cells); i++ {
				if cells[i].Style.Fg == l.TextStyle.Fg && cells[i].Style.Bg == l.TextStyle.Bg {
					cells[i].Style = l.SelectedStyle
				}
			}
		}
	}
	return cells
}

// drawRow draws a single logical row, handling text wrapping and splitting.
func (l *List) drawRow(buf *ui.Buffer, row int, y int, gradientColors []ui.Color) int {
	cells := l.getRowCells(row)

	if l.WrapText {
		cells = ui.WrapCells(cells, uint(l.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')
	for _, rowCells := range rows {
		if y >= l.Inner.Max.Y {
			break
		}
		l.drawRowLine(buf, row, y, rowCells, gradientColors)
		y++
	}
	return y
}

// drawRowLine draws a single line of text within a row, applying alignment and cell styles.
func (l *List) drawRowLine(buf *ui.Buffer, row int, y int, rowCells []ui.Cell, gradientColors []ui.Color) {
	xOffset := 0
	rowWidth := 0
	for _, c := range rowCells {
		rowWidth += rw.RuneWidth(c.Rune)
	}

	switch l.TextAlignment {
	case ui.AlignCenter:
		xOffset = (l.Inner.Dx() - rowWidth) / 2
	case ui.AlignRight:
		xOffset = l.Inner.Dx() - rowWidth
	}

	if xOffset < 0 {
		xOffset = 0
	}

	x := l.Inner.Min.X + xOffset
	for _, cell := range rowCells {
		if x >= l.Inner.Max.X {
			break
		}

		if l.Gradient.Enabled && l.Gradient.Direction == 1 {
			relativeY := y - l.Inner.Min.Y
			if relativeY >= 0 && relativeY < len(gradientColors) {
				cell.Style.Fg = gradientColors[relativeY]
			}
		}

		if row == l.SelectedRow && l.Gradient.Enabled {
			cell.Style.Modifier |= l.SelectedStyle.Modifier
			if l.SelectedStyle.Bg != ui.ColorClear {
				cell.Style.Bg = l.SelectedStyle.Bg
			} else {
				cell.Style.Modifier |= ui.ModifierReverse
			}
		}

		if x >= l.Inner.Min.X {
			buf.SetCell(cell, image.Pt(x, y))
		}
		x += rw.RuneWidth(cell.Rune)
	}
}

// drawArrows draws scroll arrows if content is out of view.
func (l *List) drawArrows(buf *ui.Buffer) {
	if l.topRow > 0 {
		buf.SetCell(
			ui.NewCell(ui.UP_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Min.Y),
		)
	}

	if len(l.Rows) > l.topRow+l.Inner.Dy() {
		buf.SetCell(
			ui.NewCell(ui.DOWN_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Max.Y-1),
		)
	}
}

// ScrollUp scrolls the list up by one row.
func (l *List) ScrollUp() {
	l.ScrollAmount(-1)
}

// ScrollDown scrolls the list down by one row.
func (l *List) ScrollDown() {
	l.ScrollAmount(1)
}

// ScrollAmount scrolls the list by the given amount.
func (l *List) ScrollAmount(amount int) {
	if len(l.Rows)-int(l.SelectedRow) <= amount {
		l.SelectedRow = len(l.Rows) - 1
	} else if int(l.SelectedRow)+amount < 0 {
		l.SelectedRow = 0
	} else {
		l.SelectedRow += amount
	}
}

// ScrollPageUp scrolls the list up by one page.
func (l *List) ScrollPageUp() {
	if l.SelectedRow > l.topRow {
		l.SelectedRow = l.topRow
	} else {
		l.ScrollAmount(-l.Inner.Dy())
	}
}

// ScrollPageDown scrolls the list down by one page.
func (l *List) ScrollPageDown() {
	l.ScrollAmount(l.Inner.Dy())
}

// ScrollHalfPageUp scrolls the list up by half a page.
func (l *List) ScrollHalfPageUp() {
	l.ScrollAmount(-int(ui.FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

// ScrollHalfPageDown scrolls the list down by half a page.
func (l *List) ScrollHalfPageDown() {
	l.ScrollAmount(int(ui.FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

// ScrollTop scrolls the list to the top.
func (l *List) ScrollTop() {
	l.SelectedRow = 0
}

// ScrollBottom scrolls the list to the bottom.
func (l *List) ScrollBottom() {
	l.SelectedRow = len(l.Rows) - 1
}
