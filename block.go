package gotui

import (
	"image"
)

// NewBlock returns a new Block.
func NewBlock() *Block {
	return &Block{
		Border:               true,
		BorderStyle:          Theme.Block.Border,
		BorderLeft:           true,
		BorderRight:          true,
		BorderTop:            true,
		BorderBottom:         true,
		BorderCollapse:       false,
		TitleStyle:           Theme.Block.Title,
		TitleAlignment:       AlignLeft,
		TitleBottomStyle:     Theme.Block.Title,
		TitleBottomAlignment: AlignLeft,
	}
}

// drawBorder draws the border of the block to the buffer.
func (b *Block) drawBorder(buf *Buffer) {
	var gradientSchema []Color
	if b.BorderGradient.Enabled {
		if len(b.BorderGradient.Stops) > 0 {
			if b.BorderGradient.Direction == GradientVertical {
				gradientSchema = GenerateMultiGradient(b.Dy(), b.BorderGradient.Stops...)
			} else {
				gradientSchema = GenerateMultiGradient(b.Dx(), b.BorderGradient.Stops...)
			}
		} else {
			if b.BorderGradient.Direction == GradientVertical {
				gradientSchema = GenerateGradient(b.BorderGradient.Start, b.BorderGradient.End, b.Dy())
			} else {
				gradientSchema = GenerateGradient(b.BorderGradient.Start, b.BorderGradient.End, b.Dx())
			}
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
	top = HORIZONTAL_LINE
	bottom = HORIZONTAL_LINE
	left = VERTICAL_LINE
	right = VERTICAL_LINE
	tl = TOP_LEFT
	tr = TOP_RIGHT
	bl = BOTTOM_LEFT
	br = BOTTOM_RIGHT
	if b.BorderSet != nil {
		top = b.BorderSet.Top
		bottom = b.BorderSet.Bottom
		left = b.BorderSet.Left
		right = b.BorderSet.Right
		tl = b.BorderSet.TopLeft
		tr = b.BorderSet.TopRight
		bl = b.BorderSet.BottomLeft
		br = b.BorderSet.BottomRight
		return
	}
	if b.BorderRounded {
		tl = ROUNDED_TOP_LEFT
		tr = ROUNDED_TOP_RIGHT
		bl = ROUNDED_BOTTOM_LEFT
		br = ROUNDED_BOTTOM_RIGHT
		return
	}
	switch b.BorderType {
	case BorderBlock:
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

// Draw draws the block to the buffer.
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
		if bgRect.Min.X < bgRect.Max.X && bgRect.Min.Y < bgRect.Max.Y {
			buf.Fill(bgCell, bgRect)
		}
	}
}
func (b *Block) drawTitles(buf *Buffer) {
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
	if b.TitleLeft != "" {
		buf.SetString(
			b.TitleLeft,
			b.TitleStyle,
			image.Pt(b.Min.X+2, b.Min.Y),
		)
	}
	if b.TitleRight != "" {
		buf.SetString(
			b.TitleRight,
			b.TitleStyle,
			image.Pt(b.Max.X-len(b.TitleRight)-2, b.Min.Y),
		)
	}
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
	if b.TitleBottomLeft != "" {
		buf.SetString(
			b.TitleBottomLeft,
			b.TitleBottomStyle,
			image.Pt(b.Min.X+2, b.Max.Y-1),
		)
	}
	if b.TitleBottomRight != "" {
		buf.SetString(
			b.TitleBottomRight,
			b.TitleBottomStyle,
			image.Pt(b.Max.X-len(b.TitleBottomRight)-2, b.Max.Y-1),
		)
	}
}

// SetRect sets the rectangle of the block.
func (b *Block) SetRect(x1, y1, x2, y2 int) {
	b.Rectangle = image.Rect(x1, y1, x2, y2)
	innerMinX := b.Min.X + 1 + b.PaddingLeft
	innerMinY := b.Min.Y + 1 + b.PaddingTop
	innerMaxX := b.Max.X - 1 - b.PaddingRight
	innerMaxY := b.Max.Y - 1 - b.PaddingBottom
	if innerMinX > innerMaxX {
		mid := b.Min.X + (b.Max.X-b.Min.X)/2
		innerMinX = mid
		innerMaxX = mid
	}
	if innerMinY > innerMaxY {
		mid := b.Min.Y + (b.Max.Y-b.Min.Y)/2
		innerMinY = mid
		innerMaxY = mid
	}
	b.Inner = image.Rect(innerMinX, innerMinY, innerMaxX, innerMaxY)
}

// GetRect returns the rectangle of the block.
func (b *Block) GetRect() image.Rectangle {
	return b.Rectangle
}
