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

	BorderCollapse bool

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
		Border:         true,
		BorderStyle:    Theme.Block.Border,
		BorderLeft:     true,
		BorderRight:    true,
		BorderTop:      true,
		BorderBottom:   true,
		BorderCollapse: false,

		TitleStyle:           Theme.Block.Title,
		TitleAlignment:       AlignLeft,
		TitleBottomStyle:     Theme.Block.Title,
		TitleBottomAlignment: AlignLeft,
	}
}

func (b *Block) drawBorder(buf *Buffer) {
	// Helper to draw a rune with optional merge
	drawRune := func(r rune, p image.Point) {
		if b.BorderCollapse {
			existing := buf.GetCell(p).Rune
			r = ResolveBorderRune(existing, r)
		}
		buf.SetCell(Cell{r, b.BorderStyle}, p)
	}

	// draw lines
	if b.BorderTop {
		xStart := b.Min.X
		xEnd := b.Max.X
		if b.BorderLeft {
			xStart++
		}
		if b.BorderRight {
			xEnd--
		}

		for x := xStart; x < xEnd; x++ {
			drawRune(HORIZONTAL_LINE, image.Pt(x, b.Min.Y))
		}
	}
	if b.BorderBottom {
		xStart := b.Min.X
		xEnd := b.Max.X
		if b.BorderLeft {
			xStart++
		}
		if b.BorderRight {
			xEnd--
		}

		for x := xStart; x < xEnd; x++ {
			drawRune(HORIZONTAL_LINE, image.Pt(x, b.Max.Y-1))
		}
	}
	if b.BorderLeft {
		yStart := b.Min.Y
		yEnd := b.Max.Y
		if b.BorderTop {
			yStart++
		}
		if b.BorderBottom {
			yEnd--
		}

		for y := yStart; y < yEnd; y++ {
			drawRune(VERTICAL_LINE, image.Pt(b.Min.X, y))
		}
	}
	if b.BorderRight {
		yStart := b.Min.Y
		yEnd := b.Max.Y
		if b.BorderTop {
			yStart++
		}
		if b.BorderBottom {
			yEnd--
		}

		for y := yStart; y < yEnd; y++ {
			drawRune(VERTICAL_LINE, image.Pt(b.Max.X-1, y))
		}
	}

	// draw corners
	if b.BorderTop && b.BorderLeft {
		drawRune(TOP_LEFT, b.Min)
	}
	// Handle cases where only one border exists at corner?
	// If BorderTop=true, BorderLeft=false. We drew the line starting at MinX. So it's covered.
	// If BorderTop=true, BorderLeft=true. We skipped MinX in loop. We draw Corner here. Covered.

	// BUT what (Top=false, Left=true)?
	// Left loop starts at MinY.
	// Corner logic skipped. So we have `|` at (0,0). Correct.

	if b.BorderTop && b.BorderRight {
		drawRune(TOP_RIGHT, image.Pt(b.Max.X-1, b.Min.Y))
	}
	if b.BorderBottom && b.BorderLeft {
		drawRune(BOTTOM_LEFT, image.Pt(b.Min.X, b.Max.Y-1))
	}
	if b.BorderBottom && b.BorderRight {
		drawRune(BOTTOM_RIGHT, b.Max.Sub(image.Pt(1, 1)))
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
