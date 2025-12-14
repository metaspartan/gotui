package widgets

import (
	ui "github.com/metaspartan/gotui/v4"
	"image"
)

type FunnelChart struct {
	ui.Block
	Data          []float64
	Labels        []string
	Colors        []ui.Color
	UniformHeight bool
}

func NewFunnelChart() *FunnelChart {
	return &FunnelChart{
		Block:         *ui.NewBlock(),
		UniformHeight: true,
		Colors:        ui.Theme.BarChart.Bars,
	}
}
func (fc *FunnelChart) Draw(buf *ui.Buffer) {
	fc.Block.Draw(buf)
	if len(fc.Data) == 0 {
		return
	}
	maxVal, _ := ui.GetMaxFloat64FromSlice(fc.Data)
	if maxVal == 0 {
		return
	}
	totalHeight := fc.Inner.Dy()
	sectionHeight := 0
	if fc.UniformHeight {
		sectionHeight = totalHeight / len(fc.Data)
	}
	canvas := ui.NewCanvas()
	canvas.SetRect(fc.Inner.Min.X, fc.Inner.Min.Y, fc.Inner.Max.X, fc.Inner.Max.Y)
	canvas.Border = false
	width := float64(fc.Inner.Dx() * 2)
	currentY := 0.0
	for i, val := range fc.Data {
		h := float64(sectionHeight)
		if !fc.UniformHeight {
			h = float64(totalHeight) * (val / ui.SumFloat64Slice(fc.Data))
		}
		sectionTopY := float64(fc.Inner.Min.Y)*4 + currentY*4
		sectionBottomY := sectionTopY + h*4
		wVal := (val / maxVal) * width
		wTop := wVal
		var wBottom float64
		if i < len(fc.Data)-1 {
			wNext := (fc.Data[i+1] / maxVal) * width
			wBottom = wNext
		} else {
			wBottom = wTop * 0.5
		}
		canvasW := float64(fc.Inner.Dx() * 2)
		x1_top := (canvasW - wTop) / 2.0
		x2_top := x1_top + wTop
		x1_bot := (canvasW - wBottom) / 2.0
		x2_bot := x1_bot + wBottom
		color := ui.SelectColor(fc.Colors, i)
		if color == ui.ColorClear || color == 0 {
			color = ui.ColorWhite
		}
		canvas.SetLine(image.Pt(int(x1_top), int(sectionTopY)), image.Pt(int(x2_top), int(sectionTopY)), color)
		canvas.SetLine(image.Pt(int(x1_bot), int(sectionBottomY)), image.Pt(int(x2_bot), int(sectionBottomY)), color)
		canvas.SetLine(image.Pt(int(x1_top), int(sectionTopY)), image.Pt(int(x1_bot), int(sectionBottomY)), color)
		canvas.SetLine(image.Pt(int(x2_top), int(sectionTopY)), image.Pt(int(x2_bot), int(sectionBottomY)), color)
		for y := int(sectionTopY); y < int(sectionBottomY); y++ {
			progress := float64(y-int(sectionTopY)) / float64(int(sectionBottomY)-int(sectionTopY))
			currX1 := x1_top + (x1_bot-x1_top)*progress
			currX2 := x2_top + (x2_bot-x2_top)*progress
			canvas.SetLine(image.Pt(int(currX1), y), image.Pt(int(currX2), y), color)
		}
		if i < len(fc.Labels) {
			ly := fc.Inner.Min.Y + int(currentY) + int(h)/2
			lx := fc.Inner.Min.X + fc.Inner.Dx()/2 - len(fc.Labels[i])/2
			buf.SetString(fc.Labels[i], ui.NewStyle(ui.ColorWhite, ui.ColorClear), image.Pt(lx, ly))
		}
		currentY += h
	}
	canvas.Draw(buf)
}
