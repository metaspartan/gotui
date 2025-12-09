package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui"
)

type List struct {
	ui.Block
	Rows          []string
	WrapText      bool
	TextStyle     ui.Style
	SelectedStyle ui.Style
	TextAlignment ui.Alignment
	SelectedRow   int
	topRow        int
}

func NewList() *List {
	return &List{
		Block:         *ui.NewBlock(),
		TextStyle:     ui.Theme.List.Text,
		SelectedStyle: ui.Theme.List.Text,
		TextAlignment: ui.AlignLeft,
	}
}

func (l *List) Draw(buf *ui.Buffer) {
	l.Block.Draw(buf)

	point := l.Inner.Min

	// adjusts view into widget
	if l.SelectedRow >= l.Inner.Dy()+l.topRow {
		l.topRow = l.SelectedRow - l.Inner.Dy() + 1
	} else if l.SelectedRow < l.topRow {
		l.topRow = l.SelectedRow
	}

	// draw rows
	for row := l.topRow; row < len(l.Rows) && point.Y < l.Inner.Max.Y; row++ {
		cells := ui.ParseStyles(l.Rows[row], l.TextStyle)
		if l.WrapText {
			cells = ui.WrapCells(cells, uint(l.Inner.Dx()))
		}

		// Apply Selected Style
		if row == l.SelectedRow {
			for i := 0; i < len(cells); i++ {
				if cells[i].Style.Fg == l.TextStyle.Fg && cells[i].Style.Bg == l.TextStyle.Bg {
					cells[i].Style = l.SelectedStyle
				}
			}
		}

		rows := ui.SplitCells(cells, '\n')
		for _, rowCells := range rows {
			if point.Y >= l.Inner.Max.Y {
				break
			}

			// Calculate alignment offset
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

			x := point.X + xOffset
			for _, cell := range rowCells {
				if x >= l.Inner.Max.X {
					break
				}
				if x >= l.Inner.Min.X {
					buf.SetCell(cell, image.Pt(x, point.Y))
				}
				x += rw.RuneWidth(cell.Rune)
			}
			point.Y++
		}
	}

	// draw UP_ARROW if needed
	if l.topRow > 0 {
		buf.SetCell(
			ui.NewCell(ui.UP_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(l.Rows) > l.topRow+l.Inner.Dy() {
		buf.SetCell(
			ui.NewCell(ui.DOWN_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set l.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (l *List) ScrollAmount(amount int) {
	if len(l.Rows)-int(l.SelectedRow) <= amount {
		l.SelectedRow = len(l.Rows) - 1
	} else if int(l.SelectedRow)+amount < 0 {
		l.SelectedRow = 0
	} else {
		l.SelectedRow += amount
	}
}

func (l *List) ScrollUp() {
	l.ScrollAmount(-1)
}

func (l *List) ScrollDown() {
	l.ScrollAmount(1)
}

func (l *List) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if l.SelectedRow > l.topRow {
		l.SelectedRow = l.topRow
	} else {
		l.ScrollAmount(-l.Inner.Dy())
	}
}

func (l *List) ScrollPageDown() {
	l.ScrollAmount(l.Inner.Dy())
}

func (l *List) ScrollHalfPageUp() {
	l.ScrollAmount(-int(ui.FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

func (l *List) ScrollHalfPageDown() {
	l.ScrollAmount(int(ui.FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

func (l *List) ScrollTop() {
	l.SelectedRow = 0
}

func (l *List) ScrollBottom() {
	l.SelectedRow = len(l.Rows) - 1
}
