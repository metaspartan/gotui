package widgets

import (
	ui "github.com/metaspartan/gotui/v4"
	"image"
)

type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

type flexItem struct {
	Widget     ui.Drawable
	FixedSize  int
	Proportion int
	Focus      bool
}
type Flex struct {
	ui.Block
	Items     []*flexItem
	Direction FlexDirection
}

func NewFlex() *Flex {
	return &Flex{
		Block:     *ui.NewBlock(),
		Items:     make([]*flexItem, 0),
		Direction: FlexColumn,
	}
}
func (f *Flex) AddItem(widget ui.Drawable, fixedSize, proportion int, focus bool) {
	f.Items = append(f.Items, &flexItem{
		Widget:     widget,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
	})
}
func (f *Flex) Draw(buf *ui.Buffer) {
	f.Block.Draw(buf)
	rect := f.Inner
	totalSize := rect.Dx()
	if f.Direction == FlexColumn {
		totalSize = rect.Dy()
	}
	usedSize := 0
	totalProportion := 0
	for _, item := range f.Items {
		if item.FixedSize > 0 {
			usedSize += item.FixedSize
		} else {
			totalProportion += item.Proportion
		}
	}
	remainingSize := totalSize - usedSize
	if remainingSize < 0 {
		remainingSize = 0
	}
	currentPos := 0
	if f.Direction == FlexRow {
		currentPos = rect.Min.X
	} else {
		currentPos = rect.Min.Y
	}
	for _, item := range f.Items {
		size := 0
		if item.FixedSize > 0 {
			size = item.FixedSize
		} else if totalProportion > 0 {
			size = remainingSize * item.Proportion / totalProportion
		}
		var childRect image.Rectangle
		if f.Direction == FlexRow {
			if currentPos+size > rect.Max.X {
				size = rect.Max.X - currentPos
			}
			childRect = image.Rect(currentPos, rect.Min.Y, currentPos+size, rect.Max.Y)
			currentPos += size
		} else {
			if currentPos+size > rect.Max.Y {
				size = rect.Max.Y - currentPos
			}
			childRect = image.Rect(rect.Min.X, currentPos, rect.Max.X, currentPos+size)
			currentPos += size
		}
		if size > 0 {
			item.Widget.SetRect(childRect.Min.X, childRect.Min.Y, childRect.Max.X, childRect.Max.Y)
			item.Widget.Draw(buf)
		}
	}
}
