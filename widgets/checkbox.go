package widgets

import (
	"image"

	gotui "github.com/metaspartan/gotui/v4"
	ui "github.com/metaspartan/gotui/v4"
)

// Checkbox represents a checkbox widget.
type Checkbox struct {
	gotui.Block
	Label       string
	Checked     bool
	CheckedRune rune
	TextStyle   gotui.Style
}

// NewCheckbox returns a new Checkbox with the given label.
func NewCheckbox(label string) *Checkbox {
	return &Checkbox{
		Block:       *gotui.NewBlock(),
		Label:       label,
		CheckedRune: 'x',
		TextStyle:   gotui.NewStyle(gotui.ColorWhite),
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
