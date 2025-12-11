// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

type Heatmap struct {
	ui.Block
	Data      [][]float64
	CellWidth int
	CellGap   int
	XLabels   []string
	YLabels   []string
	Colors    []ui.Color // Gradient colors from low to high
	TextColor ui.Style
}

func NewHeatmap() *Heatmap {
	return &Heatmap{
		Block:     *ui.NewBlock(),
		CellWidth: 3,
		CellGap:   1,
		Colors:    []ui.Color{ui.ColorBlack, ui.ColorRed, ui.ColorYellow, ui.ColorWhite}, // Default gradient
		TextColor: ui.Theme.Paragraph.Text,
	}
}

func (h *Heatmap) Draw(buf *ui.Buffer) {
	h.Block.Draw(buf)

	if len(h.Data) == 0 {
		return
	}

	// Calculate max value for normalization
	maxVal := 0.0
	for _, row := range h.Data {
		for _, val := range row {
			if val > maxVal {
				maxVal = val
			}
		}
	}

	// Draw Data
	y := h.Inner.Min.Y
	for r, row := range h.Data {
		if y >= h.Inner.Max.Y {
			break
		}
		x := h.Inner.Min.X

		// Draw Y Label if present
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

			// Determine color
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

			// Draw cell (gap is handled by just moving x forward more than width, or drawing spaces?
			// usually Draw handles filling a rect.
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

	// Draw X Labels (Simplified: assumed to be at bottom or top? Let's put them at the bottom inside if room)
	// Or maybe just below the data?
}
