package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"

	ui "github.com/metaspartan/gotui/v4"
)

// Button represents a clickable button widget.
type Button struct {
	ui.Block
	Text        string
	TextStyle   ui.Style
	ActiveStyle ui.Style
	IsActive    bool
	OnClick     func()
}

// NewButton returns a new Button with the given text.
func NewButton(text string) *Button {
	return &Button{
		Block:       *ui.NewBlock(),
		Text:        text,
		TextStyle:   ui.NewStyle(ui.ColorWhite),
		ActiveStyle: ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
	}
}

// Draw draws the button to the buffer.
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
