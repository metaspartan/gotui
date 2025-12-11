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
	if Screen == nil {
		return
	}
	if len(items) == 0 {
		return
	}

	if len(items) == 0 {
		return
	}
	// Calculate the union rectangle for all items
	if len(items) == 0 {
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

	for point, cell := range buf.CellMap {
		if point.In(buf.Rectangle) {
			style := tcell.StyleDefault.
				Foreground(cell.Style.Fg).
				Background(cell.Style.Bg).
				Attributes(cell.Style.Modifier)

			Screen.SetContent(
				point.X, point.Y,
				cell.Rune,
				nil,
				style,
			)
		}
	}
	Screen.Show()
}
