package gotui

import (
	"image"
	"os"

	"github.com/gdamore/tcell/v2"
)

func Render(items ...Drawable) {
	DefaultBackend.Render(items...)
}

func (b *Backend) Render(items ...Drawable) {
	if b.Screen == nil || len(items) == 0 {
		return
	}
	minX, minY := items[0].GetRect().Min.X, items[0].GetRect().Min.Y
	maxX, maxY := items[0].GetRect().Max.X, items[0].GetRect().Max.Y

	for _, item := range items {
		r := item.GetRect()
		if r.Min.X < minX {
			minX = r.Min.X
		}
		if r.Min.Y < minY {
			minY = r.Min.Y
		}
		if r.Max.X > maxX {
			maxX = r.Max.X
		}
		if r.Max.Y > maxY {
			maxY = r.Max.Y
		}
	}

	buf := NewBuffer(image.Rect(minX, minY, maxX, maxY))

	for _, item := range items {
		item.Lock()
		item.Draw(buf)
		item.Unlock()
	}

	if b.ScreenshotMode {
		width, height := 120, 60

		if err := SaveImage("screenshot.png", width, height, items...); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	for i, cell := range buf.Cells {
		x := (i % buf.Dx()) + buf.Min.X
		y := (i / buf.Dx()) + buf.Min.Y

		if cell.Rune == 0 {
			continue
		}

		style := tcell.StyleDefault.
			Foreground(cell.Style.Fg).
			Background(cell.Style.Bg).
			Attributes(cell.Style.Modifier)

		b.Screen.SetContent(
			x, y,
			cell.Rune,
			nil,
			style,
		)
	}
	b.Screen.Show()
}
