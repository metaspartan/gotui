package gotui

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// SaveImage renders items to an image and saves it to the specified path as PNG.
func SaveImage(path string, width, height int, items ...Drawable) error {
	img := Capture(width, height, items...)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// Capture renders the provided Drawables to an image.RGBA using the specified dimensions.
func Capture(width, height int, items ...Drawable) *image.RGBA {
	buf := NewBuffer(image.Rect(0, 0, width, height))
	for _, item := range items {
		item.Draw(buf)
	}
	return RenderBufferToImage(buf)
}

func RenderBufferToImage(buf *Buffer) *image.RGBA {
	charWidth, charHeight, fontFace, brailleFace, symbolFace, emojiFace := loadFonts()

	imgWidth := buf.Max.X * charWidth
	imgHeight := buf.Max.Y * charHeight
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	ascent := fontFace.Metrics().Ascent.Ceil()

	for x := 0; x < buf.Max.X; x++ {
		for y := 0; y < buf.Max.Y; y++ {
			drawCell(img, buf.GetCell(image.Pt(x, y)), x, y, charWidth, charHeight, ascent, fontFace, brailleFace, symbolFace, emojiFace)
		}
	}

	return img
}

func loadFonts() (int, int, font.Face, font.Face, font.Face, font.Face) {
	var fontFace font.Face = basicfont.Face7x13
	var brailleFace font.Face = basicfont.Face7x13
	var symbolFace font.Face
	var emojiFace font.Face
	w, h := 7, 13

	if face, err := loadFontFromFile("/System/Library/Fonts/Menlo.ttc", 0); err == nil {
		fontFace = face
		h = 15
	}

	if face, err := loadFontFromFile("/System/Library/Fonts/Apple Braille.ttf", 0); err == nil {
		brailleFace = face
	} else if face, err := loadFontFromFile("/System/Library/Fonts/Apple Symbols.ttf", 0); err == nil {
		brailleFace = face
	}

	if face, err := loadFontFromFile("/System/Library/Fonts/Apple Symbols.ttf", 0); err == nil {
		symbolFace = face
	}

	// Try loading Noto Emoji (Vector) for perfect screenshot rendering
	// Check multiple potential locations relative to where the binary might be running
	notoPaths := []string{
		"_fonts/NotoEmoji-Regular.ttf",
		"../_fonts/NotoEmoji-Regular.ttf",
		"../../_fonts/NotoEmoji-Regular.ttf",
		"../../../_fonts/NotoEmoji-Regular.ttf",
	}

	for _, path := range notoPaths {
		if face, err := loadFontFromFile(path, 0); err == nil {
			emojiFace = face
			break
		}
	}

	if emojiFace == nil {
		// Fallback to system fonts (Monochrome fallback)
		if face, err := loadFontFromFile("/Library/Fonts/Arial Unicode.ttf", 0); err == nil {
			emojiFace = face
		} else if face, err := loadFontFromFile("/System/Library/Fonts/Supplemental/Arial Unicode.ttf", 0); err == nil {
			emojiFace = face
		}
	}

	return w, h, fontFace, brailleFace, symbolFace, emojiFace
}

func loadFontFromFile(path string, index int) (font.Face, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(bytes) > 4 && string(bytes[:4]) == "ttcf" {
		coll, err := opentype.ParseCollection(bytes)
		if err != nil {
			return nil, err
		}
		f, err := coll.Font(index)
		if err != nil {
			return nil, err
		}
		return opentype.NewFace(f, &opentype.FaceOptions{
			Size:    12,
			DPI:     72,
			Hinting: font.HintingNone,
		})
	}

	f, err := opentype.Parse(bytes)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	})
}

func drawCell(img *image.RGBA, cell Cell, x, y, cw, ch, ascent int, fontFace, brailleFace, symbolFace, emojiFace font.Face) {
	px, py := x*cw, y*ch
	fgCol, bgCol := resolveColors(cell.Style)

	if bgCol != ColorClear {
		draw.Draw(img, image.Rect(px, py, px+cw, py+ch), &image.Uniform{toNRGBA(bgCol)}, image.Point{}, draw.Src)
	}

	if cell.Rune != 0 && cell.Rune != ' ' {
		if cell.Rune >= 0x2500 && cell.Rune <= 0x259F {
			drawCustomRune(img, px, py, cw, ch, cell.Rune, toNRGBA(fgCol))
			return
		}

		if fgCol == ColorClear {
			fgCol = ColorWhite
		}

		face, asc := selectFont(cell.Rune, fontFace, brailleFace, symbolFace, emojiFace, ascent)

		drawer := &font.Drawer{
			Dst:  img,
			Src:  &image.Uniform{toNRGBA(fgCol)},
			Face: face,
			Dot:  fixed.P(px, py+asc),
		}
		drawer.DrawString(string(cell.Rune))
	}
}

func resolveColors(s Style) (Color, Color) {
	fg, bg := s.Fg, s.Bg
	if s.Modifier&ModifierReverse != 0 {
		return bg, fg
	}
	return fg, bg
}

func selectFont(r rune, def, braille, symbol, emoji font.Face, defAscent int) (font.Face, int) {
	if braille != nil && r >= 0x2800 && r <= 0x28FF {
		return braille, braille.Metrics().Ascent.Ceil()
	}
	if emoji != nil && r >= 0x1F000 {
		return emoji, emoji.Metrics().Ascent.Ceil()
	}
	if symbol != nil && r >= 0x2000 {
		return symbol, symbol.Metrics().Ascent.Ceil()
	}
	return def, defAscent
}

// drawCustomRune manually renders block/box characters to avoid font gaps
func drawCustomRune(img *image.RGBA, x, y, w, h int, r rune, col color.NRGBA) {
	// Helper to fill rect
	fill := func(x0, y0, x1, y1 int) {
		// Clamp coords
		if x0 < 0 {
			x0 = 0
		}
		if y0 < 0 {
			y0 = 0
		}
		// Absolute coords in image
		rect := image.Rect(x+x0, y+y0, x+x1, y+y1)
		draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Src)
	}

	if r >= 0x2580 && r <= 0x259F {
		drawBlockElement(fill, w, h, r)
		return
	}
	if r >= 0x2500 && r <= 0x257F {
		drawBoxDrawing(fill, w, h, r)
		return
	}
}

func drawBlockElement(fill func(int, int, int, int), w, h int, r rune) {
	switch r {
	case 0x2581: // Lower 1/8
		fill(0, h-2, w, h)
	case 0x2582: // Lower 2/8
		fill(0, h-4, w, h)
	case 0x2583: // Lower 3/8
		fill(0, h-6, w, h)
	case 0x2584: // Lower 4/8
		fill(0, h/2, w, h)
	case 0x2585: // Lower 5/8
		fill(0, h-10, w, h)
	case 0x2586: // Lower 6/8
		fill(0, h-12, w, h)
	case 0x2587: // Lower 7/8
		fill(0, h-14, w, h)
	case 0x2588: // Full Block
		fill(0, 0, w, h)
	case 0x2580: // Upper Half
		fill(0, 0, w, h/2)
	}
}

func drawBoxDrawing(fill func(int, int, int, int), w, h int, r rune) {
	switch r {
	case 0x2500:
		fill(0, h/2, w, h/2+1)
		return
	case 0x2502:
		fill(w/2, 0, w/2+1, h)
		return
	}

	if r >= 0x250C && r <= 0x2518 {
		drawBoxCorners(fill, w, h, r)
		return
	}
	if r >= 0x251C && r <= 0x253C {
		drawBoxTees(fill, w, h, r)
		return
	}

	// Rounded Corners
	drawBoxRounded(fill, w, h, r)
}

func drawBoxCorners(fill func(int, int, int, int), w, h int, r rune) {
	switch r {
	case 0x250C:
		fill(w/2, h/2, w, h/2+1)
		fill(w/2, h/2, w/2+1, h)
	case 0x2510:
		fill(0, h/2, w/2+1, h/2+1)
		fill(w/2, h/2, w/2+1, h)
	case 0x2514:
		fill(w/2, h/2, w, h/2+1)
		fill(w/2, 0, w/2+1, h/2+1)
	case 0x2518:
		fill(0, h/2, w/2+1, h/2+1)
		fill(w/2, 0, w/2+1, h/2+1)
	}
}

func drawBoxTees(fill func(int, int, int, int), w, h int, r rune) {
	switch r {
	case 0x251C:
		fill(w/2, 0, w/2+1, h)
		fill(w/2, h/2, w, h/2+1)
	case 0x2524:
		fill(w/2, 0, w/2+1, h)
		fill(0, h/2, w/2+1, h/2+1)
	case 0x252C:
		fill(0, h/2, w, h/2+1)
		fill(w/2, h/2, w/2+1, h)
	case 0x2534:
		fill(0, h/2, w, h/2+1)
		fill(w/2, 0, w/2+1, h/2+1)
	case 0x253C:
		fill(0, h/2, w, h/2+1)
		fill(w/2, 0, w/2+1, h)
	}
}

func drawBoxRounded(fill func(int, int, int, int), w, h int, r rune) {
	switch r {
	case 0x256D:
		fill(w/2, h/2, w, h/2+1)
		fill(w/2, h/2, w/2+1, h)
	case 0x256E:
		fill(0, h/2, w/2+1, h/2+1)
		fill(w/2, h/2, w/2+1, h)
	case 0x256F:
		fill(0, h/2, w/2+1, h/2+1)
		fill(w/2, 0, w/2+1, h/2+1)
	case 0x2570:
		fill(w/2, h/2, w, h/2+1)
		fill(w/2, 0, w/2+1, h/2+1)
	}
}

// toNRGBA converts the library's Color to stdlib color.NRGBA
func toNRGBA(c Color) color.NRGBA {
	// If default/clear, assume black for now (or handle transparency if RGBA allowed)
	if c == ColorClear {
		return color.NRGBA{0, 0, 0, 255}
	}

	r, g, b := c.RGB()
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
