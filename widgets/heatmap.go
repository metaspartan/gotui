package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Heatmap represents a widget that displays a heat map.
type Heatmap struct {
	ui.Block
	Data      [][]float64
	CellWidth int
	CellGap   int
	XLabels   []string
	YLabels   []string
	Colors    []ui.Color
	TextColor ui.Style
}

// NewHeatmap returns a new Heatmap.
func NewHeatmap() *Heatmap {
	return &Heatmap{
		Block:     *ui.NewBlock(),
		CellWidth: 3,
		CellGap:   1,
		Colors:    []ui.Color{ui.ColorBlack, ui.ColorRed, ui.ColorYellow, ui.ColorWhite},
		TextColor: ui.Theme.Paragraph.Text,
	}
}

// Draw draws the heatmap to the buffer.
func (h *Heatmap) Draw(buf *ui.Buffer) {
	h.Block.Draw(buf)
	if len(h.Data) == 0 {
		return
	}
	maxVal := 0.0
	for _, row := range h.Data {
		for _, val := range row {
			if val > maxVal {
				maxVal = val
			}
		}
	}
	y := h.Inner.Min.Y
	for r, row := range h.Data {
		if y >= h.Inner.Max.Y {
			break
		}
		x := h.Inner.Min.X
		if r < len(h.YLabels) {
			buf.SetString(
				h.YLabels[r],
				h.TextColor,
				image.Pt(x, y),
			)
			x += len(h.YLabels[r]) + 1
		}
		for _, val := range row {
			if x+h.CellWidth > h.Inner.Max.X {
				break
			}
			colorIdx := 0
			if maxVal > 0 {
				percent := val / maxVal
				colorIdx = int(percent * float64(len(h.Colors)-1))
			}
			if colorIdx >= len(h.Colors) {
				colorIdx = len(h.Colors) - 1
			}
			if colorIdx < 0 {
				colorIdx = 0
			}
			style := ui.NewStyle(ui.ColorWhite, h.Colors[colorIdx])
			for i := 0; i < h.CellWidth; i++ {
				buf.SetCell(
					ui.NewCell(' ', style),
					image.Pt(x+i, y),
				)
			}
			x += h.CellWidth + h.CellGap
		}
		y++
	}
}
