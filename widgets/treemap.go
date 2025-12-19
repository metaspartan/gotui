package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// TreeMapNode represents a node in the tree map.
type TreeMapNode struct {
	Value      float64
	Label      string
	Children   []*TreeMapNode
	Style      ui.Style
	X, Y, W, H int
}

// TreeMap represents a widget that displays a tree map.
type TreeMap struct {
	ui.Block
	Root      *TreeMapNode
	TextColor ui.Color
}

// NewTreeMap returns a new TreeMap.
func NewTreeMap() *TreeMap {
	return &TreeMap{
		Block:     *ui.NewBlock(),
		TextColor: ui.ColorWhite,
	}
}

// Draw draws the tree map to the buffer.
func (tm *TreeMap) Draw(buf *ui.Buffer) {
	tm.Block.Draw(buf)
	if tm.Root == nil {
		return
	}
	tm.layout(tm.Root, tm.Inner)
	tm.renderNode(buf, tm.Root)
}
func (tm *TreeMap) layout(node *TreeMapNode, area image.Rectangle) {
	node.X = area.Min.X
	node.Y = area.Min.Y
	node.W = area.Dx()
	node.H = area.Dy()
	if len(node.Children) == 0 {
		return
	}
	totalValue := 0.0
	for _, child := range node.Children {
		totalValue += child.Value
	}
	if totalValue == 0 {
		return
	}
	x := area.Min.X
	y := area.Min.Y
	width := area.Dx()
	height := area.Dy()
	horizontalSplit := width > height
	currentPos := 0.0
	for _, child := range node.Children {
		ratio := child.Value / totalValue
		var childArea image.Rectangle
		if horizontalSplit {
			w := int(float64(width) * ratio)
			childArea = image.Rect(x+int(currentPos), y, x+int(currentPos)+w, y+height)
			currentPos += float64(w)
		} else {
			h := int(float64(height) * ratio)
			childArea = image.Rect(x, y+int(currentPos), x+width, y+int(currentPos)+h)
			currentPos += float64(h)
		}
		tm.layout(child, childArea)
	}
}
func (tm *TreeMap) renderNode(buf *ui.Buffer, node *TreeMapNode) {
	rect := image.Rect(node.X, node.Y, node.X+node.W, node.Y+node.H)
	if len(node.Children) == 0 {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				cell := ui.NewCell(' ', node.Style)
				buf.SetCell(cell, image.Pt(x, y))
			}
		}
		if node.Label != "" && rect.Dx() > len(node.Label) && rect.Dy() > 1 {
			cx := rect.Min.X + (rect.Dx()-len(node.Label))/2
			cy := rect.Min.Y + rect.Dy()/2
			buf.SetString(node.Label, ui.NewStyle(tm.TextColor, node.Style.Bg), image.Pt(cx, cy))
		}
	}
	for _, child := range node.Children {
		tm.renderNode(buf, child)
	}
}
