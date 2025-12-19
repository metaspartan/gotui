package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// TabPane represents a widget that displays a list of tabs.
type TabPane struct {
	ui.Block
	TabNames         []string
	ActiveTabIndex   int
	ActiveTabStyle   ui.Style
	InactiveTabStyle ui.Style

	PadLeft   int
	PadRight  int
	TabGap    int
	Separator string
}

// NewTabPane returns a new TabPane.
func NewTabPane(names ...string) *TabPane {
	return &TabPane{
		Block:            *ui.NewBlock(),
		TabNames:         names,
		ActiveTabStyle:   ui.Theme.Tab.Active,
		InactiveTabStyle: ui.Theme.Tab.Inactive,
		PadLeft:          1,
		PadRight:         1,
		TabGap:           0,
		Separator:        "",
	}
}

// FocusLeft moves the focus to the left tab.
func (tp *TabPane) FocusLeft() {
	if tp.ActiveTabIndex > 0 {
		tp.ActiveTabIndex--
	}
}

// FocusRight moves the focus to the right tab.
func (tp *TabPane) FocusRight() {
	if tp.ActiveTabIndex < len(tp.TabNames)-1 {
		tp.ActiveTabIndex++
	}
}

// Draw draws the tab pane to the buffer.
func (tp *TabPane) Draw(buf *ui.Buffer) {
	tp.Block.Draw(buf)

	xCoordinate := tp.Inner.Min.X

	for i, name := range tp.TabNames {
		style := tp.InactiveTabStyle
		if i == tp.ActiveTabIndex {
			style = tp.ActiveTabStyle
		}

		for j := 0; j < tp.PadLeft; j++ {
			if xCoordinate < tp.Inner.Max.X {
				buf.SetCell(ui.NewCell(' ', style), image.Pt(xCoordinate, tp.Inner.Min.Y))
				xCoordinate++
			}
		}

		for _, r := range name {
			if xCoordinate < tp.Inner.Max.X {
				buf.SetCell(ui.NewCell(r, style), image.Pt(xCoordinate, tp.Inner.Min.Y))
				xCoordinate++
			}
		}

		for j := 0; j < tp.PadRight; j++ {
			if xCoordinate < tp.Inner.Max.X {
				buf.SetCell(ui.NewCell(' ', style), image.Pt(xCoordinate, tp.Inner.Min.Y))
				xCoordinate++
			}
		}

		if i < len(tp.TabNames)-1 {
			xCoordinate += tp.TabGap

			if tp.Separator != "" {
				for _, r := range tp.Separator {
					if xCoordinate < tp.Inner.Max.X {
						buf.SetCell(ui.NewCell(r, tp.InactiveTabStyle), image.Pt(xCoordinate, tp.Inner.Min.Y))
						xCoordinate++
					}
				}
				xCoordinate += tp.TabGap
			}
		}
	}
}

// ResolveClick returns the index of the tab that was clicked, or -1 if no tab was clicked.
func (tp *TabPane) ResolveClick(p image.Point) int {
	if !p.In(tp.Inner) {
		return -1
	}

	relativeX := p.X - tp.Inner.Min.X
	currentX := 0

	for i, name := range tp.TabNames {
		tabWidth := tp.PadLeft + len([]rune(name)) + tp.PadRight

		if relativeX >= currentX && relativeX < currentX+tabWidth {
			return i
		}

		currentX += tabWidth

		if i < len(tp.TabNames)-1 {
			currentX += tp.TabGap
			if tp.Separator != "" {
				currentX += len([]rune(tp.Separator))
				currentX += tp.TabGap
			}
		}
	}

	return -1
}
