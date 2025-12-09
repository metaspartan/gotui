package gotui

import (
	"image"
	"sync"
)

// Block is the base struct inherited by most widgets.
// Block manages size, position, border, and title.
// It implements all 3 of the methods needed for the `Drawable` interface.
// Custom widgets will override the Draw method.
type Block struct {
	Border      bool
	BorderStyle Style

	BorderLeft, BorderRight, BorderTop, BorderBottom bool

	PaddingLeft, PaddingRight, PaddingTop, PaddingBottom int

	image.Rectangle
	Inner image.Rectangle

	Title                string
	TitleStyle           Style
	TitleAlignment       Alignment
	TitleBottom          string
	TitleBottomStyle     Style
	TitleBottomAlignment Alignment

	sync.Mutex
}

func NewBlock() *Block {
	return &Block{
		Border:       true,
		BorderStyle:  Theme.Block.Border,
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,

		TitleStyle:           Theme.Block.Title,
		TitleAlignment:       AlignLeft,
		TitleBottomStyle:     Theme.Block.Title,
		TitleBottomAlignment: AlignLeft,
	}
}

func (b *Block) drawBorder(buf *Buffer) {
	verticalCell := Cell{VERTICAL_LINE, b.BorderStyle}
	horizontalCell := Cell{HORIZONTAL_LINE, b.BorderStyle}

	// draw lines
	if b.BorderTop {
		buf.Fill(horizontalCell, image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Min.Y+1))
	}
	if b.BorderBottom {
		buf.Fill(horizontalCell, image.Rect(b.Min.X, b.Max.Y-1, b.Max.X, b.Max.Y))
	}
	if b.BorderLeft {
		buf.Fill(verticalCell, image.Rect(b.Min.X, b.Min.Y, b.Min.X+1, b.Max.Y))
	}
	if b.BorderRight {
		buf.Fill(verticalCell, image.Rect(b.Max.X-1, b.Min.Y, b.Max.X, b.Max.Y))
	}

	// draw corners
	if b.BorderTop && b.BorderLeft {
		buf.SetCell(Cell{TOP_LEFT, b.BorderStyle}, b.Min)
	}
	if b.BorderTop && b.BorderRight {
		buf.SetCell(Cell{TOP_RIGHT, b.BorderStyle}, image.Pt(b.Max.X-1, b.Min.Y))
	}
	if b.BorderBottom && b.BorderLeft {
		buf.SetCell(Cell{BOTTOM_LEFT, b.BorderStyle}, image.Pt(b.Min.X, b.Max.Y-1))
	}
	if b.BorderBottom && b.BorderRight {
		buf.SetCell(Cell{BOTTOM_RIGHT, b.BorderStyle}, b.Max.Sub(image.Pt(1, 1)))
	}
}

// Draw implements the Drawable interface.
func (b *Block) Draw(buf *Buffer) {
	if b.Border {
		b.drawBorder(buf)
	}

	// Top Title
	titleX := b.Min.X + 2
	switch b.TitleAlignment {
	case AlignCenter:
		titleX = b.Min.X + (b.Max.X-b.Min.X-len(b.Title))/2
	case AlignRight:
		titleX = b.Max.X - len(b.Title) - 2
	}

	buf.SetString(
		b.Title,
		b.TitleStyle,
		image.Pt(titleX, b.Min.Y),
	)

	// Bottom Title
	bottomTitleX := b.Min.X + 2
	switch b.TitleBottomAlignment {
	case AlignCenter:
		bottomTitleX = b.Min.X + (b.Max.X-b.Min.X-len(b.TitleBottom))/2
	case AlignRight:
		bottomTitleX = b.Max.X - len(b.TitleBottom) - 2
	}

	buf.SetString(
		b.TitleBottom,
		b.TitleBottomStyle,
		image.Pt(bottomTitleX, b.Max.Y-1),
	)
}

// SetRect implements the Drawable interface.
func (b *Block) SetRect(x1, y1, x2, y2 int) {
	b.Rectangle = image.Rect(x1, y1, x2, y2)
	b.Inner = image.Rect(
		b.Min.X+1+b.PaddingLeft,
		b.Min.Y+1+b.PaddingTop,
		b.Max.X-1-b.PaddingRight,
		b.Max.Y-1-b.PaddingBottom,
	)
}

// GetRect implements the Drawable interface.
func (b *Block) GetRect() image.Rectangle {
	return b.Rectangle
}
