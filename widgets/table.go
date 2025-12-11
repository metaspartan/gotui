package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

/*
Table is like:
┌ Awesome Table ───────────────────────────────────────────────┐
│  Col0          | Col1 | Col2 | Col3  | Col4  | Col5  | Col6  |
│──────────────────────────────────────────────────────────────│
│  Some Item #1  | AAA  | 123  | CCCCC | EEEEE | GGGGG | IIIII |
│──────────────────────────────────────────────────────────────│
│  Some Item #2  | BBB  | 456  | DDDDD | FFFFF | HHHHH | JJJJJ |
└──────────────────────────────────────────────────────────────┘
*/
type Table struct {
	ui.Block
	Rows          [][]string
	ColumnWidths  []int
	TextStyle     ui.Style
	RowSeparator  bool
	TextAlignment ui.Alignment
	RowStyles     map[int]ui.Style
	FillRow       bool

	// ColumnResizer is called on each Draw. Can be used for custom column sizing.
	ColumnResizer func()
}

func NewTable() *Table {
	return &Table{
		Block:         *ui.NewBlock(),
		TextStyle:     ui.Theme.Table.Text,
		RowSeparator:  true,
		RowStyles:     make(map[int]ui.Style),
		ColumnResizer: func() {},
	}
}

func (tb *Table) Draw(buf *ui.Buffer) {
	tb.Block.Draw(buf)

	tb.ColumnResizer()

	columnWidths := tb.ColumnWidths
	if len(columnWidths) == 0 {
		columnCount := len(tb.Rows[0])
		columnWidth := tb.Inner.Dx() / columnCount
		for i := 0; i < columnCount; i++ {
			columnWidths = append(columnWidths, columnWidth)
		}
	}

	yCoordinate := tb.Inner.Min.Y

	// draw rows
	for i := 0; i < len(tb.Rows) && yCoordinate < tb.Inner.Max.Y; i++ {
		row := tb.Rows[i]
		colXCoordinate := tb.Inner.Min.X

		rowStyle := tb.TextStyle
		// get the row style if one exists
		if style, ok := tb.RowStyles[i]; ok {
			rowStyle = style
		}

		if tb.FillRow {
			blankCell := ui.NewCell(' ', rowStyle)
			buf.Fill(blankCell, image.Rect(tb.Inner.Min.X, yCoordinate, tb.Inner.Max.X, yCoordinate+1))
		}

		// draw row cells
		for j := 0; j < len(row); j++ {
			col := ui.ParseStyles(row[j], rowStyle)
			// draw row cell
			if len(col) > columnWidths[j] || tb.TextAlignment == ui.AlignLeft {
				for _, cx := range ui.BuildCellWithXArray(col) {
					k, cell := cx.X, cx.Cell
					if k == columnWidths[j] || colXCoordinate+k == tb.Inner.Max.X {
						cell.Rune = ui.ELLIPSES
						buf.SetCell(cell, image.Pt(colXCoordinate+k-1, yCoordinate))
						break
					} else {
						buf.SetCell(cell, image.Pt(colXCoordinate+k, yCoordinate))
					}
				}
			} else if tb.TextAlignment == ui.AlignCenter {
				xCoordinateOffset := (columnWidths[j] - len(col)) / 2
				stringXCoordinate := xCoordinateOffset + colXCoordinate
				for _, cx := range ui.BuildCellWithXArray(col) {
					k, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(stringXCoordinate+k, yCoordinate))
				}
			} else if tb.TextAlignment == ui.AlignRight {
				stringXCoordinate := ui.MinInt(colXCoordinate+columnWidths[j], tb.Inner.Max.X) - len(col)
				for _, cx := range ui.BuildCellWithXArray(col) {
					k, cell := cx.X, cx.Cell
					buf.SetCell(cell, image.Pt(stringXCoordinate+k, yCoordinate))
				}
			}
			colXCoordinate += columnWidths[j] + 1
		}

		// draw vertical separators
		separatorStyle := tb.Block.BorderStyle

		separatorXCoordinate := tb.Inner.Min.X
		verticalCell := ui.NewCell(ui.VERTICAL_LINE, separatorStyle)
		for i, width := range columnWidths {
			if tb.FillRow && i < len(columnWidths)-1 {
				verticalCell.Style.Bg = rowStyle.Bg
			} else {
				verticalCell.Style.Bg = tb.Block.BorderStyle.Bg
			}

			separatorXCoordinate += width
			buf.SetCell(verticalCell, image.Pt(separatorXCoordinate, yCoordinate))
			separatorXCoordinate++
		}

		yCoordinate++

		// draw horizontal separator
		horizontalCell := ui.NewCell(ui.HORIZONTAL_LINE, separatorStyle)
		if tb.RowSeparator && yCoordinate < tb.Inner.Max.Y && i != len(tb.Rows)-1 {
			buf.Fill(horizontalCell, image.Rect(tb.Inner.Min.X, yCoordinate, tb.Inner.Max.X, yCoordinate+1))
			yCoordinate++
		}
	}
}
