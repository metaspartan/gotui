package widgets

import (
	"image"
	"math"

	ui "github.com/metaspartan/gotui/v5"
)

// ScrollbarOrientation represents the orientation of the scrollbar.
type ScrollbarOrientation int

const (
	ScrollbarVertical ScrollbarOrientation = iota
	ScrollbarHorizontal
)

// Scrollbar represents a widget that displays a scrollbar.
type Scrollbar struct {
	ui.Block
	Orientation ScrollbarOrientation
	Max         int
	Current     int
	PageSize    int
	ThumbStyle  ui.Style
	TrackStyle  ui.Style
	ThumbRune   rune
	TrackRune   rune
	BeginRune   rune
	EndRune     rune
}

// NewScrollbar returns a new Scrollbar.
func NewScrollbar() *Scrollbar {
	return &Scrollbar{
		Block:       *ui.NewBlock(),
		Orientation: ScrollbarVertical,
		Max:         100,
		Current:     0,
		PageSize:    10,
		ThumbStyle:  ui.NewStyle(ui.ColorWhite),
		TrackStyle:  ui.NewStyle(ui.ColorBlack),
		ThumbRune:   '█',
		TrackRune:   '║',
		BeginRune:   '▲',
		EndRune:     '▼',
	}
}

// Draw draws the scrollbar to the buffer.
func (s *Scrollbar) Draw(buf *ui.Buffer) {
	s.Block.Draw(buf)
	if s.Max <= 0 {
		return
	}
	rect := s.Inner
	totalSize := 0
	if s.Orientation == ScrollbarVertical {
		totalSize = rect.Dy()
	} else {
		totalSize = rect.Dx()
	}
	if totalSize <= 0 {
		return
	}
	arrowStart := 0
	arrowEnd := 0
	if s.BeginRune != 0 {
		arrowStart = 1
	}
	if s.EndRune != 0 {
		arrowEnd = 1
	}
	trackLen := totalSize - arrowStart - arrowEnd
	if trackLen <= 0 {
		return
	}
	viewportRatio := float64(s.PageSize) / float64(s.Max)
	if viewportRatio > 1.0 {
		viewportRatio = 1.0
	}
	thumbSize := int(math.Max(1.0, float64(trackLen)*viewportRatio))
	moveableSpace := trackLen - thumbSize
	scrollRatio := 0.0
	if s.Max > s.PageSize {
		scrollRatio = float64(s.Current) / float64(s.Max-s.PageSize)
	}
	thumbPos := int(scrollRatio * float64(moveableSpace))
	if thumbPos < 0 {
		thumbPos = 0
	}
	if thumbPos+thumbSize > trackLen {
		thumbPos = trackLen - thumbSize
	}
	if s.Orientation == ScrollbarVertical {
		s.drawVertical(buf, rect, totalSize, arrowStart, arrowEnd, thumbPos, thumbSize)
	} else {
		s.drawHorizontal(buf, rect, totalSize, arrowStart, arrowEnd, thumbPos, thumbSize)
	}
}
func (s *Scrollbar) drawVertical(buf *ui.Buffer, rect image.Rectangle, totalSize, arrowStart, arrowEnd, thumbPos, thumbSize int) {
	for y := 0; y < rect.Dy(); y++ {
		py := rect.Min.Y + y
		px := rect.Min.X
		for x := 0; x < rect.Dx(); x++ {
			var char rune
			style := s.TrackStyle
			if y < arrowStart {
				char = s.BeginRune
			} else if y >= totalSize-arrowEnd {
				char = s.EndRune
			} else {
				trackY := y - arrowStart
				char = s.TrackRune
				if trackY >= thumbPos && trackY < thumbPos+thumbSize {
					style = s.ThumbStyle
					char = s.ThumbRune
				}
			}
			if char != 0 {
				buf.SetCell(ui.NewCell(char, style), image.Pt(px+x, py))
			}
		}
	}
}
func (s *Scrollbar) drawHorizontal(buf *ui.Buffer, rect image.Rectangle, totalSize, arrowStart, arrowEnd, thumbPos, thumbSize int) {
	for x := 0; x < rect.Dx(); x++ {
		px := rect.Min.X + x
		py := rect.Min.Y
		for y := 0; y < rect.Dy(); y++ {
			var char rune
			style := s.TrackStyle
			if x < arrowStart {
				char = s.BeginRune
			} else if x >= totalSize-arrowEnd {
				char = s.EndRune
			} else {
				trackX := x - arrowStart
				char = s.TrackRune
				if trackX >= thumbPos && trackX < thumbPos+thumbSize {
					style = s.ThumbStyle
					char = s.ThumbRune
				}
			}
			if char != 0 {
				buf.SetCell(ui.NewCell(char, style), image.Pt(px, py+y))
			}
		}
	}
}
