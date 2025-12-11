package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

type FunnelChart struct {
	ui.Block
	Data   []float64
	Labels []string
	Colors []ui.Color
	// If false, all trapezoids have height proportional to value
	// If true, uniform height
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

	// Use utility function
	maxVal, _ := ui.GetMaxFloat64FromSlice(fc.Data)

	if maxVal == 0 {
		return
	}
	// For funnel, top is usually maxVal (or 100%), bottom is minVal.
	// But actually, we want the width of each section to be proportional to its value relative to MaxVal.

	totalHeight := fc.Inner.Dy()
	sectionHeight := 0
	if fc.UniformHeight {
		sectionHeight = totalHeight / len(fc.Data)
	}

	canvas := ui.NewCanvas()
	canvas.SetRect(fc.Inner.Min.X, fc.Inner.Min.Y, fc.Inner.Max.X, fc.Inner.Max.Y)
	canvas.Border = false

	// We'll use canvas for trapezoids to get smoother lines if possible,
	// but filling them might be tricky with just lines.
	// Let's stick to simple block rendering or canvas for edges?
	// For "modern" look, lines + fill is nice.
	// Let's us drawille canvas for edges and fill.

	// Braille canvas resolution
	width := float64(fc.Inner.Dx() * 2)

	currentY := 0.0

	// Let's assume top is full width if it's the biggest?
	// Common funnel: Top is widest. Width is proportional to value.

	for i, val := range fc.Data {
		h := float64(sectionHeight)
		if !fc.UniformHeight {
			h = float64(totalHeight) * (val / ui.SumFloat64Slice(fc.Data)) // Proportional height?
		}

		// Calculate top and bottom widths for this section (trapezoid)
		// Usually funnel sections don't separate, they flow.
		// So TopWidth of Iter[i] = BottomWidth of Iter[i-1] ?
		// But if Data is [100, 80, 60], then:
		// S1: Top=100, Bottom=??? (Interpolate to next?) or just Trapezoid representing 100?
		// Typically Funnel Chart: Each stage represents a value.
		// Width at center of stage = value?
		// Or: Top Width = Value[i], Bottom Width = Value[i+1]?
		// If discrete stages: Top Width = Value[i], Bottom Width = Value[i] (Rectangle) -> Pyramid?
		// Let's do: Top Width = Value[i], Bottom Width = Value[i] * 0.8? No.
		// Let's do: Trapezoid where Top edge is Value[i] and Bottom edge is Value[i+1] (or 0 for last).

		// Wait, a funnel usually shows drop off.
		// Value[i] is the volume at stage i.
		// So Stage i should have width proportional to Value[i].
		// But to make it look connected:
		// Top of Stage i needs to connect to Bottom of Stage i-1.
		// So Top Width = (i==0) ? MaxWidth : Value[i-1 relative to Max] * MaxAllWidth?
		// Actually, standard Funnel:
		// Stage 0: 100 users.
		// Stage 1: 50 users.
		// We want Stage 0 to look wider than Stage 1.
		// Visual: Trapezoid for Stage 0?
		// Top Base = Max Width. Bottom Base = Width proportional to 50?
		// No, Stage 0 represents 100.
		// Stage 1 represents 50.
		// If we stack them:
		// Rect 0 (Width 100)
		// Rect 1 (Width 50)
		// Centered.
		// Then connect edges.

		sectionTopY := float64(fc.Inner.Min.Y)*4 + currentY*4
		sectionBottomY := sectionTopY + h*4

		// Widths in braille units (2x cell width)
		// Scale: MaxVal -> full widget width
		wVal := (val / maxVal) * width

		// If we want a smooth funnel, we might interpolate.
		// But discrete stages usually look like stacked blocks with angled sides.
		// Let's just draw centered trapezoids where:
		// Top Width = Previous Value's Width (or this value if first) ??
		// No, that implies flow.

		// Let's implement specific style:
		// Top Edge Width = wVal
		// Bottom Edge Width = wNext (or wVal * 0.8 for last one to make it look like a funnel tip)

		wTop := wVal
		wBottom := wVal
		if i < len(fc.Data)-1 {
			wNext := (fc.Data[i+1] / maxVal) * width
			// wBottom = (wTop + wNext) / 2.0 // Average to smooth transition?
			// Or just set to wNext?
			// If we set Bottom = wNext, then Stage i transitions from Val[i] to Val[i+1].
			wBottom = wNext
		} else {
			// wBottom = 0 // Pointy tip?
			// Or just keep it straight?
			wBottom = wTop * 0.5
		}

		// Center them
		canvasW := float64(fc.Inner.Dx() * 2)
		x1_top := (canvasW - wTop) / 2.0
		x2_top := x1_top + wTop

		x1_bot := (canvasW - wBottom) / 2.0
		x2_bot := x1_bot + wBottom

		// Draw trapezoid outline
		color := ui.SelectColor(fc.Colors, i)
		if color == ui.ColorClear || color == 0 {
			color = ui.ColorWhite
		}
		canvas.SetLine(image.Pt(int(x1_top), int(sectionTopY)), image.Pt(int(x2_top), int(sectionTopY)), color)       // Top
		canvas.SetLine(image.Pt(int(x1_bot), int(sectionBottomY)), image.Pt(int(x2_bot), int(sectionBottomY)), color) // Bottom
		canvas.SetLine(image.Pt(int(x1_top), int(sectionTopY)), image.Pt(int(x1_bot), int(sectionBottomY)), color)    // Left
		canvas.SetLine(image.Pt(int(x2_top), int(sectionTopY)), image.Pt(int(x2_bot), int(sectionBottomY)), color)    // Right

		// Fill?
		// Simple scanline fill for trapezoid
		// For y from top to bottom, interp x start/end
		for y := int(sectionTopY); y < int(sectionBottomY); y++ {
			progress := float64(y-int(sectionTopY)) / float64(int(sectionBottomY)-int(sectionTopY))

			currX1 := x1_top + (x1_bot-x1_top)*progress
			currX2 := x2_top + (x2_bot-x2_top)*progress

			// Fill line
			// Optimization: Draw simple horizontal line
			canvas.SetLine(image.Pt(int(currX1), y), image.Pt(int(currX2), y), color)
		}

		// Label?
		if i < len(fc.Labels) {
			// Center label in segment
			// Need to use buffer coordinates
			ly := fc.Inner.Min.Y + int(currentY) + int(h)/2
			lx := fc.Inner.Min.X + fc.Inner.Dx()/2 - len(fc.Labels[i])/2
			buf.SetString(fc.Labels[i], ui.NewStyle(ui.ColorWhite, ui.ColorClear), image.Pt(lx, ly))
		}

		currentY += h
	}

	canvas.Draw(buf)
}
