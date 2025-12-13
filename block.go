package gotui

import (
	"image"
)

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
	var gradientSchema []Color

	if b.BorderGradient.Enabled {
		if b.BorderGradient.Direction == GradientVertical {
			gradientSchema = GenerateGradient(b.BorderGradient.Start, b.BorderGradient.End, b.Dy())
		} else {
			// Default to Horizontal
			gradientSchema = GenerateGradient(b.BorderGradient.Start, b.BorderGradient.End, b.Dx())
		}
	}

	drawRune := func(r rune, p image.Point) {
		if b.BorderCollapse {
			existing := buf.GetCell(p).Rune
			r = ResolveBorderRune(existing, r)
		}
		style := b.BorderStyle
		if b.BackgroundColor != ColorClear && b.FillBorder {
			style.Bg = b.BackgroundColor
		}

		if b.BorderGradient.Enabled {
			var idx int
			if b.BorderGradient.Direction == GradientVertical {
				idx = p.Y - b.Min.Y
			} else {
				idx = p.X - b.Min.X
			}

			if idx >= 0 && idx < len(gradientSchema) {
				style.Fg = gradientSchema[idx]
			}
		}

		buf.SetCell(Cell{r, style}, p)
	}

	b.drawBorderLines(drawRune)
	b.drawBorderCorners(drawRune)
}

func (b *Block) getBorderRunes() (top, bottom, left, right, tl, tr, bl, br rune) {
	// Default (BorderLine)
	top = HORIZONTAL_LINE
	bottom = HORIZONTAL_LINE
	left = VERTICAL_LINE
	right = VERTICAL_LINE
	tl = TOP_LEFT
	tr = TOP_RIGHT
	bl = BOTTOM_LEFT
	br = BOTTOM_RIGHT

	if b.BorderRounded {
		tl = ROUNDED_TOP_LEFT
		tr = ROUNDED_TOP_RIGHT
		bl = ROUNDED_BOTTOM_LEFT
		br = ROUNDED_BOTTOM_RIGHT
		return
	}

	switch b.BorderType {
	case BorderBlock: // "Thick/Unique" Block Style
		top = '▀'
		bottom = '▄'
		left = '▌'
		right = '▐'
		tl = '█'
		tr = '█'
		bl = '█'
		br = '█'
	case BorderDouble:
		top = '═'
		bottom = '═'
		left = '║'
		right = '║'
		tl = '╔'
		tr = '╗'
		bl = '╚'
		br = '╝'
	case BorderThick:
		top = '━'
		bottom = '━'
		left = '┃'
		right = '┃'
		tl = '┏'
		tr = '┓'
		bl = '┗'
		br = '┛'
	}
	return
}

func (b *Block) drawBorderLines(drawRune func(rune, image.Point)) {
	top, bottom, left, right, _, _, _, _ := b.getBorderRunes()

	if b.BorderTop {
		b.drawHorizontalBorder(drawRune, b.Min.Y, top)
	}
	if b.BorderBottom {
		b.drawHorizontalBorder(drawRune, b.Max.Y-1, bottom)
	}
	if b.BorderLeft {
		b.drawVerticalBorder(drawRune, b.Min.X, left)
	}
	if b.BorderRight {
		b.drawVerticalBorder(drawRune, b.Max.X-1, right)
	}
}

func (b *Block) drawHorizontalBorder(drawRune func(rune, image.Point), y int, r rune) {
	xStart := b.Min.X
	xEnd := b.Max.X
	if b.BorderLeft {
		xStart++
	}
	if b.BorderRight {
		xEnd--
	}
	for x := xStart; x < xEnd; x++ {
		drawRune(r, image.Pt(x, y))
	}
}

func (b *Block) drawVerticalBorder(drawRune func(rune, image.Point), x int, r rune) {
	yStart := b.Min.Y
	yEnd := b.Max.Y
	if b.BorderTop {
		yStart++
	}
	if b.BorderBottom {
		yEnd--
	}
	for y := yStart; y < yEnd; y++ {
		drawRune(r, image.Pt(x, y))
	}
}

func (b *Block) drawBorderCorners(drawRune func(rune, image.Point)) {
	_, _, _, _, tl, tr, bl, br := b.getBorderRunes()

	if b.BorderTop && b.BorderLeft {
		drawRune(tl, b.Min)
	}
	if b.BorderTop && b.BorderRight {
		drawRune(tr, image.Pt(b.Max.X-1, b.Min.Y))
	}
	if b.BorderBottom && b.BorderLeft {
		drawRune(bl, image.Pt(b.Min.X, b.Max.Y-1))
	}
	if b.BorderBottom && b.BorderRight {
		drawRune(br, b.Max.Sub(image.Pt(1, 1)))
	}
}

// Draw implements the Drawable interface.
func (b *Block) Draw(buf *Buffer) {
	b.drawBackground(buf)

	if b.Border {
		b.drawBorder(buf)
	}

	b.drawTitles(buf)
}

func (b *Block) drawBackground(buf *Buffer) {
	if b.BackgroundColor != ColorClear {
		bgCell := NewCell(' ', NewStyle(ColorClear, b.BackgroundColor))

		bgRect := b.Rectangle
		if !b.FillBorder && b.Border {
			if b.BorderTop {
				bgRect.Min.Y++
			}
			if b.BorderBottom {
				bgRect.Max.Y--
			}
			if b.BorderLeft {
				bgRect.Min.X++
			}
			if b.BorderRight {
				bgRect.Max.X--
			}
		}

		// Ensure bgRect is valid (Min <= Max)
		if bgRect.Min.X < bgRect.Max.X && bgRect.Min.Y < bgRect.Max.Y {
			buf.Fill(bgCell, bgRect)
		}
	}
}

func (b *Block) drawTitles(buf *Buffer) {
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

	// Top Left Title (Explicit)
	if b.TitleLeft != "" {
		buf.SetString(
			b.TitleLeft,
			b.TitleStyle,
			image.Pt(b.Min.X+2, b.Min.Y),
		)
	}

	// Top Right Title (Explicit)
	if b.TitleRight != "" {
		buf.SetString(
			b.TitleRight,
			b.TitleStyle,
			image.Pt(b.Max.X-len(b.TitleRight)-2, b.Min.Y),
		)
	}

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

	// Bottom Left Title (Explicit)
	if b.TitleBottomLeft != "" {
		buf.SetString(
			b.TitleBottomLeft,
			b.TitleBottomStyle,
			image.Pt(b.Min.X+2, b.Max.Y-1),
		)
	}

	// Bottom Right Title (Explicit)
	if b.TitleBottomRight != "" {
		buf.SetString(
			b.TitleBottomRight,
			b.TitleBottomStyle,
			image.Pt(b.Max.X-len(b.TitleBottomRight)-2, b.Max.Y-1),
		)
	}
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
