package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui/v4"
)

type Button struct {
	ui.Block
	Text        string
	TextStyle   ui.Style
	ActiveStyle ui.Style
	IsActive    bool
}

func NewButton(text string) *Button {
	return &Button{
		Block:       *ui.NewBlock(),
		Text:        text,
		TextStyle:   ui.NewStyle(ui.ColorWhite),
		ActiveStyle: ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
	}
}

func (b *Button) Draw(buf *ui.Buffer) {
	b.Block.Draw(buf)

	style := b.TextStyle
	if b.IsActive {
		style = b.ActiveStyle
	}

	prefix := "❰ "
	suffix := " ❱"
	str := prefix + b.Text + suffix

	textWidth := rw.StringWidth(str)
	innerDx := b.Inner.Dx()

	x := b.Inner.Min.X + (innerDx-textWidth)/2
	y := b.Inner.Min.Y + (b.Inner.Dy()-1)/2

	if x < b.Inner.Min.X {
		x = b.Inner.Min.X
	}

	buf.SetString(str, style, image.Pt(x, y))
}

func (b *Button) Activate() {
	b.IsActive = true
}

func (b *Button) Deactivate() {
	b.IsActive = false
}
