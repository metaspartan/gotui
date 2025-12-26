package widgets

import (
	"image"

	ui "github.com/metaspartan/gotui/v5"
)

// Sparkline represents a single sparkline.
type Sparkline struct {
	Data            []float64
	Title           string
	TitleStyle      ui.Style
	LineColor       ui.Color
	BackgroundColor ui.Color
	MaxVal          float64
	MaxHeight       int
}

// SparklineGroup represents a group of sparklines.
type SparklineGroup struct {
	ui.Block
	Sparklines []*Sparkline
}

// NewSparkline returns a new Sparkline.
func NewSparkline() *Sparkline {
	return &Sparkline{
		TitleStyle:      ui.Theme.Sparkline.Title,
		LineColor:       ui.Theme.Sparkline.Line,
		BackgroundColor: ui.ColorClear,
	}
}

// NewSparklineGroup returns a new SparklineGroup.
func NewSparklineGroup(sls ...*Sparkline) *SparklineGroup {
	return &SparklineGroup{
		Block:      *ui.NewBlock(),
		Sparklines: sls,
	}
}

// Draw draws the sparkline group to the buffer.
func (sg *SparklineGroup) Draw(buf *ui.Buffer) {
	sg.Block.Draw(buf)
	sparklineHeight := sg.Inner.Dy() / len(sg.Sparklines)
	for i, sl := range sg.Sparklines {
		heightOffset := (sparklineHeight * (i + 1))
		barHeight := sparklineHeight
		if i == len(sg.Sparklines)-1 {
			heightOffset = sg.Inner.Dy()
			barHeight = sg.Inner.Dy() - (sparklineHeight * i)
		}
		if sl.Title != "" {
			barHeight--
		}
		maxVal := sl.MaxVal
		if maxVal == 0 {
			maxVal, _ = ui.GetMaxFloat64FromSlice(sl.Data)
		}
		lineColor := sl.LineColor
		if lineColor == ui.ColorClear || lineColor == 0 {
			lineColor = ui.ColorWhite
		}
		dataLen := len(sl.Data)
		width := sg.Inner.Dx()
		startIndex := 0
		if dataLen > width {
			startIndex = dataLen - width
		}
		for j := 0; j < width && startIndex+j < dataLen; j++ {
			data := sl.Data[startIndex+j]
			height := int((data / maxVal) * float64(barHeight))
			if sl.MaxHeight > 0 && height > sl.MaxHeight {
				height = sl.MaxHeight
			}
			sparkChar := ui.BARS[len(ui.BARS)-1]
			for k := 0; k < height; k++ {
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(lineColor, sl.BackgroundColor)),
					image.Pt(j+sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset-k),
				)
			}
			if height == 0 {
				sparkChar = ui.BARS[1]
				buf.SetCell(
					ui.NewCell(sparkChar, ui.NewStyle(lineColor, sl.BackgroundColor)),
					image.Pt(j+sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset),
				)
			}
		}
		if sl.Title != "" {
			buf.SetString(
				ui.TrimString(sl.Title, sg.Inner.Dx()),
				sl.TitleStyle,
				image.Pt(sg.Inner.Min.X, sg.Inner.Min.Y-1+heightOffset-barHeight),
			)
		}
	}
}
