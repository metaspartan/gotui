package gotui

import (
	"image"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

// Render renders the collection of items to the default backend.
func Render(items ...Drawable) {
	DefaultBackend.Render(items...)
}

// Render renders the collection of items to the backend's screen.
func (b *Backend) Render(items ...Drawable) {
	if b.Screen == nil || len(items) == 0 {
		return
	}
	// Calculate the union rectangle for all items
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

	// If ScreenshotMode is active, render to file and exit
	if b.ScreenshotMode {
		// Default size if not detected
		width, height := 1024/7, 768/13 // approx 146x59
		// Or 120x40
		width, height = 120, 60

		if err := SaveImage("screenshot.png", width, height, items...); err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	for i, cell := range buf.Cells {
		// Calculate Point from index
		x := (i % buf.Dx()) + buf.Min.X
		y := (i / buf.Dx()) + buf.Min.Y

		if cell.Rune == 0 {
			continue // skip empty cells if needed, or render them?
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
