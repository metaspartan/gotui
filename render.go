
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
	for _, item := range items {
		buf := NewBuffer(item.GetRect())
		item.Lock()
		item.Draw(buf)
		item.Unlock()

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
	}
	Screen.Show()
}
