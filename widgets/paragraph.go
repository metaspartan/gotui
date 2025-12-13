package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui/v4"
)

type Paragraph struct {
	ui.Block
	Text              string
	TextStyle         ui.Style
	WrapText          bool
	VerticalAlignment ui.VerticalAlignment
	TextAlignment     ui.Alignment
	Gradient          ui.Gradient
}

func NewParagraph() *Paragraph {
	return &Paragraph{
		Block:             *ui.NewBlock(),
		TextStyle:         ui.Theme.Paragraph.Text,
		WrapText:          true,
		VerticalAlignment: ui.AlignTop,
		TextAlignment:     ui.AlignLeft,
	}
}

func (p *Paragraph) Draw(buf *ui.Buffer) {
	p.Block.Draw(buf)

	cells := p.computeCells()
	rows := ui.SplitCells(cells, '\n')
	p.drawRows(buf, rows)
}

func (p *Paragraph) computeCells() []ui.Cell {
	var cells []ui.Cell
	if p.Gradient.Enabled && p.Gradient.Direction == 0 {
		cells = ui.ApplyGradientToText(p.Text, p.Gradient.Start, p.Gradient.End)
	} else {
		cells = ui.ParseStyles(p.Text, p.TextStyle)
	}

	if p.WrapText {
		cells = ui.WrapCells(cells, uint(p.Inner.Dx()))
	}
	return cells
}

func (p *Paragraph) drawRows(buf *ui.Buffer, rows [][]ui.Cell) {
	height := p.Inner.Dy()
	totalRows := len(rows)
	topPadding := 0

	switch p.VerticalAlignment {
	case ui.AlignMiddle:
		topPadding = (height - totalRows) / 2
	case ui.AlignBottom:
		topPadding = height - totalRows
	}
	if topPadding < 0 {
		topPadding = 0
	}

	var gradientColors []ui.Color
	if p.Gradient.Enabled && p.Gradient.Direction == 1 {
		gradientColors = ui.GenerateGradient(p.Gradient.Start, p.Gradient.End, height)
	}

	for i, row := range rows {
		y := i + topPadding
		if y >= height {
			break
		}

		row = ui.TrimCells(row, p.Inner.Dx())
		cellWithX := ui.BuildCellWithXArray(row)

		if len(cellWithX) == 0 {
			continue
		}

		xOffset := p.calculateXOffset(cellWithX)

		for _, cx := range cellWithX {
			x, cell := cx.X, cx.Cell
			if p.Gradient.Enabled && p.Gradient.Direction == 1 && y < len(gradientColors) {
				cell.Style = ui.NewStyle(gradientColors[y])
			}
			buf.SetCell(cell, image.Pt(x+xOffset, y).Add(p.Inner.Min))
		}
	}
}

func (p *Paragraph) calculateXOffset(cellWithX []ui.CellWithX) int {
	last := cellWithX[len(cellWithX)-1]
	rowWidth := last.X + rw.RuneWidth(last.Cell.Rune)

	switch p.TextAlignment {
	case ui.AlignCenter:
		return (p.Inner.Dx() - rowWidth) / 2
	case ui.AlignRight:
		return p.Inner.Dx() - rowWidth
	}
	return 0
}
