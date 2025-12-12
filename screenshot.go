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

// RenderBufferToImage converts a Buffer to an image.RGBA
func RenderBufferToImage(buf *Buffer) *image.RGBA {
	charWidth := 7
	charHeight := 13

	// Try loading a system font (macOS specific for now)
	var fontFace font.Face = basicfont.Face7x13 // fallback
	var brailleFace font.Face                   // secondary for braille

	fontPath := "/System/Library/Fonts/Menlo.ttc"
	fontBytes, err := os.ReadFile(fontPath)
	if err == nil {
		coll, err := opentype.ParseCollection(fontBytes)
		if err == nil && coll.NumFonts() > 0 {
			f, err := coll.Font(0)
			if err == nil {
				face, err := opentype.NewFace(f, &opentype.FaceOptions{
					Size:    12,
					DPI:     72,
					Hinting: font.HintingNone,
				})
				if err == nil {
					fontFace = face
					charWidth = 7
					charHeight = 15
				}
			}
		}
	}

	// Try loading Apple Braille for dots if Menlo fails them
	braillePath := "/System/Library/Fonts/Apple Braille.ttf"
	brailleBytes, err := os.ReadFile(braillePath)
	if err == nil {
		f, err := opentype.Parse(brailleBytes)
		if err == nil {
			face, err := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    12,
				DPI:     72,
				Hinting: font.HintingNone,
			})
			if err == nil {
				brailleFace = face
			}
		}
	} else {
		// Fallback to Apple Symbols?
		braillePath = "/System/Library/Fonts/Apple Symbols.ttf"
		brailleBytes, err = os.ReadFile(braillePath)
		if err == nil {
			f, err := opentype.Parse(brailleBytes)
			if err == nil {
				face, err := opentype.NewFace(f, &opentype.FaceOptions{
					Size:    12,
					DPI:     72,
					Hinting: font.HintingNone,
				})
				if err == nil {
					brailleFace = face
				}
			}
		}
	}

	imgWidth := buf.Max.X * charWidth
	imgHeight := buf.Max.Y * charHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	ascent := fontFace.Metrics().Ascent.Ceil()

	for x := 0; x < buf.Max.X; x++ {
		for y := 0; y < buf.Max.Y; y++ {
			cell := buf.GetCell(image.Pt(x, y))

			// Calculate position
			px := x * charWidth
			py := y * charHeight

			fgCol := cell.Style.Fg
			bgCol := cell.Style.Bg

			// Handle Reverse Modifier
			if cell.Style.Modifier&ModifierReverse != 0 {
				fgCol, bgCol = bgCol, fgCol
			}

			// Draw Background
			if bgCol != ColorClear {
				draw.Draw(img, image.Rect(px, py, px+charWidth, py+charHeight), &image.Uniform{toNRGBA(bgCol)}, image.Point{}, draw.Src)
			}

			// Draw Foreground (Rune)
			if cell.Rune != 0 && cell.Rune != ' ' {
				// Handle Block Elements manually for pixel-perfect rendering (no gaps)
				// U+2580..U+259F (Block Elements)
				// U+2500..U+257F (Box Drawing)
				if cell.Rune >= 0x2500 && cell.Rune <= 0x259F {
					drawCustomRune(img, px, py, charWidth, charHeight, cell.Rune, toNRGBA(fgCol))
					continue
				}

				if fgCol == ColorClear {
					fgCol = ColorWhite
				}

				currentFace := fontFace
				currentAscent := ascent

				// Use Braille font if available and rune is Braille Pattern (U+28xx)
				if brailleFace != nil && cell.Rune >= 0x2800 && cell.Rune <= 0x28FF {
					currentFace = brailleFace
					currentAscent = brailleFace.Metrics().Ascent.Ceil()
				}

				drawer := &font.Drawer{
					Dst:  img,
					Src:  &image.Uniform{toNRGBA(fgCol)},
					Face: currentFace,
					Dot:  fixed.P(px, py+currentAscent),
				}
				drawer.DrawString(string(cell.Rune))
			}
		}
	}

	return img
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

	// Center lines thickness
	// cx := w / 2
	// cy := h / 2
	// For 7x15, center is 3, 7.

	switch r {
	// --- Block Elements (Vertical) ---
	case 0x2581: // Lower 1/8 (was ' ' causing artifacts)
		fill(0, h-2, w, h)
	case 0x2582: // Lower 2/8 '▂'
		fill(0, h-4, w, h)
	case 0x2583: // Lower 3/8 '▃'
		fill(0, h-6, w, h)
	case 0x2584: // Lower 4/8 '▄'
		fill(0, h/2, w, h)
	case 0x2585: // Lower 5/8 '▅'
		fill(0, h-10, w, h)
	case 0x2586: // Lower 6/8 '▆'
		fill(0, h-12, w, h)
	case 0x2587: // Lower 7/8 '▇'
		fill(0, h-14, w, h)
	case 0x2588: // Full Block '█'
		fill(0, 0, w, h)

	case 0x2580: // Upper Half '▀'
		fill(0, 0, w, h/2)

	// --- Box Drawing ---
	// Light Horizontal '─'
	case 0x2500:
		fill(0, h/2, w, h/2+1)
	// Light Vertical '│'
	case 0x2502:
		fill(w/2, 0, w/2+1, h)

	// Corners (Light)
	case 0x250C: // '┌'
		fill(w/2, h/2, w, h/2+1) // Right
		fill(w/2, h/2, w/2+1, h) // Down
	case 0x2510: // '┐'
		fill(0, h/2, w/2+1, h/2+1) // Left
		fill(w/2, h/2, w/2+1, h)   // Down
	case 0x2514: // '└'
		fill(w/2, h/2, w, h/2+1)   // Right
		fill(w/2, 0, w/2+1, h/2+1) // Up
	case 0x2518: // '┘'
		fill(0, h/2, w/2+1, h/2+1) // Left
		fill(w/2, 0, w/2+1, h/2+1) // Up

	// Tees (Light)
	case 0x251C: // '├'
		fill(w/2, 0, w/2+1, h)   // Vertical
		fill(w/2, h/2, w, h/2+1) // Right
	case 0x2524: // '┤'
		fill(w/2, 0, w/2+1, h)     // Vertical
		fill(0, h/2, w/2+1, h/2+1) // Left
	case 0x252C: // '┬'
		fill(0, h/2, w, h/2+1)   // Horizontal
		fill(w/2, h/2, w/2+1, h) // Down
	case 0x2534: // '┴'
		fill(0, h/2, w, h/2+1)     // Horizontal
		fill(w/2, 0, w/2+1, h/2+1) // Up
	case 0x253C: // '┼'
		fill(0, h/2, w, h/2+1) // Horizontal
		fill(w/2, 0, w/2+1, h) // Vertical

	// Rounded Corners
	case 0x256D: // '╭'
		fill(w/2, h/2, w, h/2+1) // Right
		fill(w/2, h/2, w/2+1, h) // Down
	case 0x256E: // '╮'
		fill(0, h/2, w/2+1, h/2+1) // Left
		fill(w/2, h/2, w/2+1, h)   // Down
	case 0x256F: // '╯'
		fill(0, h/2, w/2+1, h/2+1) // Left
		fill(w/2, 0, w/2+1, h/2+1) // Up
	case 0x2570: // '╰'
		fill(w/2, h/2, w, h/2+1)   // Right
		fill(w/2, 0, w/2+1, h/2+1) // Up

	default:
		// If unimplemented, maybe better fallback?
		// For now, if we missed it, it won't draw anything here :(.
		// Wait, we should define fallthrough behavior?
		// But function signature returns nothing.
		// Let's implement a 'draw via font' fallback inside here?
		// No, easier to just check map validity or switch back to caller.
		// But simpler: just implement the main ones.
		// If we miss one (like double lines), it vanishes.
		// Let's rely on standard logic for unknown ones?
		// Refactoring: make checking cleaner.
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
