package widgets

import (
	"image"
	"image/color"

	ui "github.com/metaspartan/gotui/v5"
)

// Image represents a widget that displays an image.
type Image struct {
	ui.Block
	Image               image.Image
	Monochrome          bool
	MonochromeThreshold uint8
	MonochromeInvert    bool
}

// NewImage returns a new Image.
func NewImage(img image.Image) *Image {
	return &Image{
		Block:               *ui.NewBlock(),
		MonochromeThreshold: 128,
		Image:               img,
	}
}

// Draw draws the image to the buffer.
func (img *Image) Draw(buf *ui.Buffer) {
	img.Block.Draw(buf)

	if img.Image == nil {
		return
	}

	bufWidth := img.Inner.Dx()
	bufHeight := img.Inner.Dy()
	imageWidth := img.Image.Bounds().Dx()
	imageHeight := img.Image.Bounds().Dy()

	if img.Monochrome {
		if bufWidth > imageWidth/2 {
			bufWidth = imageWidth / 2
		}
		if bufHeight > imageHeight/2 {
			bufHeight = imageHeight / 2
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				ul := img.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ur := img.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					2*by*imageHeight/bufHeight/2,
					(2*by+1)*imageHeight/bufHeight/2,
				)
				ll := img.colorAverage(
					2*bx*imageWidth/bufWidth/2,
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				lr := img.colorAverage(
					(2*bx+1)*imageWidth/bufWidth/2,
					(2*bx+2)*imageWidth/bufWidth/2,
					(2*by+1)*imageHeight/bufHeight/2,
					(2*by+2)*imageHeight/bufHeight/2,
				)
				buf.SetCell(
					ui.NewCell(blocksChar(ul, ur, ll, lr, img.MonochromeThreshold, img.MonochromeInvert)),
					image.Pt(img.Inner.Min.X+bx, img.Inner.Min.Y+by),
				)
			}
		}
	} else {
		if bufWidth > imageWidth {
			bufWidth = imageWidth
		}
		if bufHeight > imageHeight {
			bufHeight = imageHeight
		}
		for bx := 0; bx < bufWidth; bx++ {
			for by := 0; by < bufHeight; by++ {
				c := img.colorAverage(
					bx*imageWidth/bufWidth,
					(bx+1)*imageWidth/bufWidth,
					by*imageHeight/bufHeight,
					(by+1)*imageHeight/bufHeight,
				)
				buf.SetCell(
					ui.NewCell(c.ch(), ui.NewStyle(c.fgColor(), ui.ColorBlack)),
					image.Pt(img.Inner.Min.X+bx, img.Inner.Min.Y+by),
				)
			}
		}
	}
}

func (img *Image) colorAverage(x0, x1, y0, y1 int) colorAverager {
	var c colorAverager
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			c = c.add(
				img.Image.At(
					x+img.Image.Bounds().Min.X,
					y+img.Image.Bounds().Min.Y,
				),
			)
		}
	}
	return c
}

type colorAverager struct {
	rsum, gsum, bsum, asum, count uint64
}

func (ca colorAverager) add(col color.Color) colorAverager {
	r, g, b, a := col.RGBA()
	return colorAverager{
		rsum:  ca.rsum + uint64(r),
		gsum:  ca.gsum + uint64(g),
		bsum:  ca.bsum + uint64(b),
		asum:  ca.asum + uint64(a),
		count: ca.count + 1,
	}
}

func (ca colorAverager) RGBA() (uint32, uint32, uint32, uint32) {
	if ca.count == 0 {
		return 0, 0, 0, 0
	}
	return uint32(ca.rsum/ca.count) & 0xffff,
		uint32(ca.gsum/ca.count) & 0xffff,
		uint32(ca.bsum/ca.count) & 0xffff,
		uint32(ca.asum/ca.count) & 0xffff
}

func (ca colorAverager) fgColor() ui.Color {
	if ca.count == 0 {
		return ui.ColorBlack
	}
	// Use true RGB color for better color accuracy
	r := uint8((ca.rsum / ca.count) >> 8)
	g := uint8((ca.gsum / ca.count) >> 8)
	b := uint8((ca.bsum / ca.count) >> 8)
	return ui.NewRGBColor(int32(r), int32(g), int32(b))
}

func (ca colorAverager) ch() rune {
	gray := color.GrayModel.Convert(ca).(color.Gray).Y
	switch {
	case gray < 51:
		return ui.SHADED_BLOCKS[0]
	case gray < 102:
		return ui.SHADED_BLOCKS[1]
	case gray < 153:
		return ui.SHADED_BLOCKS[2]
	case gray < 204:
		return ui.SHADED_BLOCKS[3]
	default:
		return ui.SHADED_BLOCKS[4]
	}
}

func (ca colorAverager) monochrome(threshold uint8, invert bool) bool {
	return ca.count != 0 && (color.GrayModel.Convert(ca).(color.Gray).Y < threshold != invert)
}

type paletteColor struct {
	rgba      color.RGBA
	attribute ui.Color
}

func (pc paletteColor) RGBA() (uint32, uint32, uint32, uint32) {
	return pc.rgba.RGBA()
}

var palette = color.Palette([]color.Color{
	paletteColor{color.RGBA{0, 0, 0, 255}, ui.ColorBlack},
	paletteColor{color.RGBA{255, 0, 0, 255}, ui.ColorRed},
	paletteColor{color.RGBA{0, 255, 0, 255}, ui.ColorGreen},
	paletteColor{color.RGBA{255, 255, 0, 255}, ui.ColorYellow},
	paletteColor{color.RGBA{0, 0, 255, 255}, ui.ColorBlue},
	paletteColor{color.RGBA{255, 0, 255, 255}, ui.ColorMagenta},
	paletteColor{color.RGBA{0, 255, 255, 255}, ui.ColorLightCyan},
	paletteColor{color.RGBA{255, 255, 255, 255}, ui.ColorWhite},
})

func blocksChar(ul, ur, ll, lr colorAverager, threshold uint8, invert bool) rune {
	index := 0
	if ul.monochrome(threshold, invert) {
		index |= 1
	}
	if ur.monochrome(threshold, invert) {
		index |= 2
	}
	if ll.monochrome(threshold, invert) {
		index |= 4
	}
	if lr.monochrome(threshold, invert) {
		index |= 8
	}
	return ui.IRREGULAR_BLOCKS[index]
}
