package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui"
)

// TabPane is a renderable widget which can be used to conditionally render certain tabs/views.
// TabPane shows a list of Tab names.
// The currently selected tab can be found through the `ActiveTabIndex` field.
type TabPane struct {
	ui.Block
	TabNames         []string
	ActiveTabIndex   int
	ActiveTabStyle   ui.Style
	InactiveTabStyle ui.Style
}

func NewTabPane(names ...string) *TabPane {
	return &TabPane{
		Block:            *ui.NewBlock(),
		TabNames:         names,
		ActiveTabStyle:   ui.Theme.Tab.Active,
		InactiveTabStyle: ui.Theme.Tab.Inactive,
	}
}

func (tp *TabPane) FocusLeft() {
	if tp.ActiveTabIndex > 0 {
		tp.ActiveTabIndex--
	}
}

func (tp *TabPane) FocusRight() {
	if tp.ActiveTabIndex < len(tp.TabNames)-1 {
		tp.ActiveTabIndex++
	}
}

func (tp *TabPane) Draw(buf *ui.Buffer) {
	tp.Block.Draw(buf)

	xCoordinate := tp.Inner.Min.X
	for i, name := range tp.TabNames {
		ColorPair := tp.InactiveTabStyle
		if i == tp.ActiveTabIndex {
			ColorPair = tp.ActiveTabStyle
		}
		buf.SetString(
			ui.TrimString(name, tp.Inner.Max.X-xCoordinate),
			ColorPair,
			image.Pt(xCoordinate, tp.Inner.Min.Y),
		)

		xCoordinate += 1 + len(name)

		if i < len(tp.TabNames)-1 && xCoordinate < tp.Inner.Max.X {
			buf.SetCell(
				ui.NewCell(ui.VERTICAL_LINE, ui.NewStyle(ui.ColorWhite)),
				image.Pt(xCoordinate, tp.Inner.Min.Y),
			)
		}

		xCoordinate += 2
	}
}
