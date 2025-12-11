package widgets

import (
	"image"
	"math"

	ui "github.com/metaspartan/gotui/v4"
)

const (
	piechartOffsetUp = -.5 * math.Pi // the northward angle
	resolutionFactor = .0001         // circle resolution: precision vs. performance
	fullCircle       = 2.0 * math.Pi // the full circle angle
	xStretch         = 2.0           // horizontal adjustment
)

// PieChartLabel callback
type PieChartLabel func(dataIndex int, currentValue float64) string

type PieChart struct {
	ui.Block
	Data           []float64     // list of data items
	Colors         []ui.Color    // colors to by cycled through
	LabelFormatter PieChartLabel // callback function for labels
	AngleOffset    float64       // which angle to start drawing at? (see piechartOffsetUp)
	InnerRadius    float64       // 0.0 to 1.0, defines the size of the hole
}

// NewPieChart Creates a new pie chart with reasonable defaults and no labels.
func NewPieChart() *PieChart {
	return &PieChart{
		Block:       *ui.NewBlock(),
		Colors:      ui.Theme.PieChart.Slices,
		AngleOffset: piechartOffsetUp,
		InnerRadius: 0.0,
	}
}

func (pc *PieChart) Draw(buf *ui.Buffer) {
	pc.Block.Draw(buf)

	center := pc.Inner.Min.Add(pc.Inner.Size().Div(2))
	radius := ui.MinFloat64(float64(pc.Inner.Dx()/2/xStretch), float64(pc.Inner.Dy()/2))
	innerRadius := radius * pc.InnerRadius

	// compute slice sizes
	sum := ui.SumFloat64Slice(pc.Data)
	sliceSizes := make([]float64, len(pc.Data))
	for i, v := range pc.Data {
		sliceSizes[i] = v / sum * fullCircle
	}

	borderCircle := &circle{center, radius}
	innerCircle := &circle{center, innerRadius}
	middleCircle := circle{Point: center, radius: (radius + innerRadius) / 2.0} // Middle of the donut ring

	// draw sectors
	phi := pc.AngleOffset
	for i, size := range sliceSizes {
		for j := 0.0; j < size; j += resolutionFactor {
			borderPoint := borderCircle.at(phi + j)
			innerPoint := innerCircle.at(phi + j)

			// If InnerRadius is 0, innerPoint should be center, which it is since radius is 0.

			line := line{P1: innerPoint, P2: borderPoint}
			line.draw(ui.NewCell(ui.SHADED_BLOCKS[1], ui.NewStyle(ui.SelectColor(pc.Colors, i))), buf)
		}
		phi += size
	}

	// draw labels
	if pc.LabelFormatter != nil {
		phi = pc.AngleOffset
		for i, size := range sliceSizes {
			labelPoint := middleCircle.at(phi + size/2.0)
			if len(pc.Data) == 1 {
				labelPoint = center
			}
			buf.SetString(
				pc.LabelFormatter(i, pc.Data[i]),
				ui.NewStyle(ui.SelectColor(pc.Colors, i)),
				image.Pt(labelPoint.X, labelPoint.Y),
			)
			phi += size
		}
	}
}

type circle struct {
	image.Point
	radius float64
}

// computes the point at a given angle phi
func (c circle) at(phi float64) image.Point {
	x := c.X + int(ui.RoundFloat64(xStretch*c.radius*math.Cos(phi)))
	y := c.Y + int(ui.RoundFloat64(c.radius*math.Sin(phi)))
	return image.Point{X: x, Y: y}
}

// computes the perimeter of a circle
func (c circle) perimeter() float64 {
	return 2.0 * math.Pi * c.radius
}

// a line between two points
type line struct {
	P1, P2 image.Point
}

// draws the line
func (l line) draw(cell ui.Cell, buf *ui.Buffer) {
	isLeftOf := func(p1, p2 image.Point) bool {
		return p1.X <= p2.X
	}
	isTopOf := func(p1, p2 image.Point) bool {
		return p1.Y <= p2.Y
	}
	p1, p2 := l.P1, l.P2
	buf.SetCell(ui.NewCell('*', cell.Style), l.P2)
	width, height := l.size()
	if width > height { // paint left to right
		if !isLeftOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isTopOf(p1, p2) {
			flip = -1.0
		}
		for x := p1.X; x <= p2.X; x++ {
			ratio := float64(height) / float64(width)
			factor := float64(x - p1.X)
			y := ratio * factor * flip
			buf.SetCell(cell, image.Pt(x, int(ui.RoundFloat64(y))+p1.Y))
		}
	} else { // paint top to bottom
		if !isTopOf(p1, p2) {
			p1, p2 = p2, p1
		}
		flip := 1.0
		if !isLeftOf(p1, p2) {
			flip = -1.0
		}
		for y := p1.Y; y <= p2.Y; y++ {
			ratio := float64(width) / float64(height)
			factor := float64(y - p1.Y)
			x := ratio * factor * flip
			buf.SetCell(cell, image.Pt(int(ui.RoundFloat64(x))+p1.X, y))
		}
	}
}

// width and height of a line
func (l line) size() (w, h int) {
	return ui.AbsInt(l.P2.X - l.P1.X), ui.AbsInt(l.P2.Y - l.P1.Y)
}
