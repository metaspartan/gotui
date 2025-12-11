package gotui

import (
	"image"
	"sync"

	"github.com/gdamore/tcell/v2"
)

type Drawable interface {
	GetRect() image.Rectangle
	SetRect(int, int, int, int)
	Draw(*Buffer)
	sync.Locker
}

func Render(items ...Drawable) {
	if len(items) == 0 {
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

		Screen.SetContent(
			x, y,
			cell.Rune,
			nil,
			style,
		)
	}
	Screen.Show()
}
