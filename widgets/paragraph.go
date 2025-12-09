package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui"
)

type Paragraph struct {
	ui.Block
	Text              string
	TextStyle         ui.Style
	WrapText          bool
	VerticalAlignment ui.Alignment
	TextAlignment     ui.Alignment
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

	cells := ui.ParseStyles(p.Text, p.TextStyle)
	if p.WrapText {
		cells = ui.WrapCells(cells, uint(p.Inner.Dx()))
	}

	rows := ui.SplitCells(cells, '\n')

	totalRows := len(rows)
	height := p.Inner.Dy()
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

		// Calculate Row Width
		last := cellWithX[len(cellWithX)-1]
		rowWidth := last.X + rw.RuneWidth(last.Cell.Rune)

		// Alternative: calculate offset
		xOffset := 0
		switch p.TextAlignment {
		case ui.AlignCenter:
			xOffset = (p.Inner.Dx() - rowWidth) / 2
		case ui.AlignRight:
			xOffset = p.Inner.Dx() - rowWidth
		}

		for _, cx := range cellWithX {
			x, cell := cx.X, cx.Cell
			buf.SetCell(cell, image.Pt(x+xOffset, y).Add(p.Inner.Min))
		}
	}
}
