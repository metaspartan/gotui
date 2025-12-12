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
	TextWrap      bool

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
		rowStyle := tb.TextStyle
		// get the row style if one exists
		if style, ok := tb.RowStyles[i]; ok {
			rowStyle = style
		}

		// Calculate height of this row
		rowHeight := 1
		if tb.TextWrap {
			for j, cellText := range row {
				if j < len(columnWidths) {
					width := columnWidths[j]
					cells := ui.ParseStyles(cellText, rowStyle)
					wrapped := ui.WrapCells(cells, uint(width))
					lines := ui.SplitCells(wrapped, '\n')
					if len(lines) > rowHeight {
						rowHeight = len(lines)
					}
				}
			}
		}

		tb.drawTableRow(buf, row, rowStyle, i, yCoordinate, rowHeight, columnWidths)

		yCoordinate += rowHeight

		// draw horizontal separator
		separatorStyle := tb.Block.BorderStyle
		horizontalCell := ui.NewCell(ui.HORIZONTAL_LINE, separatorStyle)
		if tb.RowSeparator && yCoordinate < tb.Inner.Max.Y && i != len(tb.Rows)-1 {
			buf.Fill(horizontalCell, image.Rect(tb.Inner.Min.X, yCoordinate, tb.Inner.Max.X, yCoordinate+1))
			yCoordinate++
		}
	}
}

func (tb *Table) drawTableRow(buf *ui.Buffer, row []string, rowStyle ui.Style, rowIndex, yCoordinate, rowHeight int, columnWidths []int) {
	colXCoordinate := tb.Inner.Min.X

	// Fill background for the entire row height
	if tb.FillRow {
		blankCell := ui.NewCell(' ', rowStyle)
		buf.Fill(blankCell, image.Rect(tb.Inner.Min.X, yCoordinate, tb.Inner.Max.X, yCoordinate+rowHeight))
	}

	// draw row cells
	for j := 0; j < len(row); j++ {
		if j >= len(columnWidths) {
			break
		}
		col := ui.ParseStyles(row[j], rowStyle)

		var lines [][]ui.Cell
		if tb.TextWrap {
			wrapped := ui.WrapCells(col, uint(columnWidths[j]))
			lines = ui.SplitCells(wrapped, '\n')
		} else {
			lines = [][]ui.Cell{col}
		}

		tb.drawTableCell(buf, lines, rowIndex, j, yCoordinate, rowHeight, colXCoordinate, columnWidths[j])

		colXCoordinate += columnWidths[j] + 1
	}

	// draw vertical separators for the full height
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
		// Fill vertical line down the rowHeight
		for h := 0; h < rowHeight; h++ {
			if yCoordinate+h < tb.Inner.Max.Y {
				buf.SetCell(verticalCell, image.Pt(separatorXCoordinate, yCoordinate+h))
			}
		}
		separatorXCoordinate++
	}
}

func (tb *Table) drawTableCell(buf *ui.Buffer, lines [][]ui.Cell, rowIndex, colIndex, yCoordinate, rowHeight, colXCoordinate, colWidth int) {
	for lineIdx := 0; lineIdx < rowHeight; lineIdx++ {
		currentY := yCoordinate + lineIdx
		if currentY >= tb.Inner.Max.Y {
			break
		}

		if lineIdx < len(lines) {
			line := lines[lineIdx]
			tb.drawTableLine(buf, line, currentY, colXCoordinate, colWidth)
		}
	}
}

func (tb *Table) drawTableLine(buf *ui.Buffer, line []ui.Cell, currentY, colXCoordinate, colWidth int) {
	if tb.TextWrap {
		switch tb.TextAlignment {
		case ui.AlignCenter:
			tb.drawCenterAligned(buf, line, currentY, colXCoordinate, colWidth)
		case ui.AlignRight:
			tb.drawRightAligned(buf, line, currentY, colXCoordinate, colWidth)
		default:
			tb.drawWrappedLeft(buf, line, currentY, colXCoordinate, colWidth)
		}
		return
	}

	if len(line) > colWidth || tb.TextAlignment == ui.AlignLeft {
		tb.drawLeftAligned(buf, line, currentY, colXCoordinate, colWidth)
	} else if tb.TextAlignment == ui.AlignCenter {
		tb.drawCenterAligned(buf, line, currentY, colXCoordinate, colWidth)
	} else if tb.TextAlignment == ui.AlignRight {
		tb.drawRightAligned(buf, line, currentY, colXCoordinate, colWidth)
	}
}

func (tb *Table) drawLeftAligned(buf *ui.Buffer, line []ui.Cell, currentY, colXCoordinate, colWidth int) {
	if len(line) > colWidth {
		for _, cx := range ui.BuildCellWithXArray(line) {
			k, cell := cx.X, cx.Cell
			if k == colWidth || colXCoordinate+k == tb.Inner.Max.X {
				cell.Rune = ui.ELLIPSES
				buf.SetCell(cell, image.Pt(colXCoordinate+k-1, currentY))
				break
			} else {
				buf.SetCell(cell, image.Pt(colXCoordinate+k, currentY))
			}
		}
	} else {
		for _, cx := range ui.BuildCellWithXArray(line) {
			k, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(colXCoordinate+k, currentY))
		}
	}
}

func (tb *Table) drawWrappedLeft(buf *ui.Buffer, line []ui.Cell, currentY, colXCoordinate, colWidth int) {
	for _, cx := range ui.BuildCellWithXArray(line) {
		k, cell := cx.X, cx.Cell
		if k < colWidth {
			buf.SetCell(cell, image.Pt(colXCoordinate+k, currentY))
		}
	}
}

func (tb *Table) drawCenterAligned(buf *ui.Buffer, line []ui.Cell, currentY, colXCoordinate, colWidth int) {
	xCoordinateOffset := (colWidth - len(line)) / 2
	stringXCoordinate := xCoordinateOffset + colXCoordinate
	for _, cx := range ui.BuildCellWithXArray(line) {
		k, cell := cx.X, cx.Cell
		buf.SetCell(cell, image.Pt(stringXCoordinate+k, currentY))
	}
}

func (tb *Table) drawRightAligned(buf *ui.Buffer, line []ui.Cell, currentY, colXCoordinate, colWidth int) {
	stringXCoordinate := ui.MinInt(colXCoordinate+colWidth, tb.Inner.Max.X) - len(line)
	for _, cx := range ui.BuildCellWithXArray(line) {
		k, cell := cx.X, cx.Cell
		buf.SetCell(cell, image.Pt(stringXCoordinate+k, currentY))
	}
}
