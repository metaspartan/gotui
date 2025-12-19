package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Checkbox represents a checkbox widget.
type Checkbox struct {
	ui.Block
	Label       string
	Checked     bool
	CheckedRune rune
	TextStyle   ui.Style
}

// NewCheckbox returns a new Checkbox with the given label.
func NewCheckbox(label string) *Checkbox {
	return &Checkbox{
		Block:       *ui.NewBlock(),
		Label:       label,
		CheckedRune: 'x',
		TextStyle:   ui.NewStyle(ui.ColorWhite),
	}
}

// Draw draws the checkbox to the buffer.
func (c *Checkbox) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)

	x := c.Inner.Min.X
	y := c.Inner.Min.Y

	buf.SetString("[ ]", c.TextStyle, image.Pt(x, y))

	if c.Checked {
		buf.SetCell(ui.NewCell(c.CheckedRune, c.TextStyle), image.Pt(x+1, y))
	}

	if c.Label != "" {
		buf.SetString(c.Label, c.TextStyle, image.Pt(x+4, y))
	}
}

func (c *Checkbox) Toggle() {
	c.Checked = !c.Checked
}
