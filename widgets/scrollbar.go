package widgets

import (
	"image"
	"math"

	ui "github.com/metaspartan/gotui/v4"
)

type ScrollbarOrientation int

const (
	ScrollbarVertical ScrollbarOrientation = iota
	ScrollbarHorizontal
)

type Scrollbar struct {
	ui.Block
	Orientation ScrollbarOrientation

	// Max is the total size of the content (e.g., number of lines or characters)
	Max int
	// Current is the current scroll position
	Current int
	// PageSize is the size of the visible area
	PageSize int

	// Styles
	ThumbStyle ui.Style
	TrackStyle ui.Style

	// Runes
	ThumbRune rune
	TrackRune rune
	BeginRune rune
	EndRune   rune
}

func NewScrollbar() *Scrollbar {
	return &Scrollbar{
		Block:       *ui.NewBlock(),
		Orientation: ScrollbarVertical,
		Max:         100,
		Current:     0,
		PageSize:    10,
		ThumbStyle:  ui.NewStyle(ui.ColorWhite),
		TrackStyle:  ui.NewStyle(ui.ColorBlack),
		// Default to Ratatui DOUBLE_VERTICAL style
		ThumbRune: '█',
		TrackRune: '║',
		BeginRune: '▲',
		EndRune:   '▼',
	}
}

func (s *Scrollbar) Draw(buf *ui.Buffer) {
	// Don't draw block background, just the scrollbar components?
	// Or maybe the block is container? Let's assume the block IS the scrollbar area.
	s.Block.Draw(buf)

	if s.Max <= 0 {
		return
	}

	// Determine renderable area
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

	// Determine rendering space for track (excluding arrows if present)
	// Arrows always take 1 cell each?
	// Ratatui logic: if space allows.

	// Assume arrows take 1 cell if set
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
		return // Not enough space
	}

	// Calculate thumb size and position within trackLen
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

	// Clamp
	if thumbPos < 0 {
		thumbPos = 0
	}
	if thumbPos+thumbSize > trackLen {
		thumbPos = trackLen - thumbSize
	}

	// Render
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

		// Fill width of inner rect
		for x := 0; x < rect.Dx(); x++ {
			var char rune
			style := s.TrackStyle

			// Check boundaries
			if y < arrowStart {
				char = s.BeginRune
			} else if y >= totalSize-arrowEnd {
				char = s.EndRune
			} else {
				// In track area
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

		// Fill height of inner rect
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
