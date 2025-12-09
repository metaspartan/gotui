package gotui

// Border Connections Bitmask
const (
	BorderTop    = 1 // 0001
	BorderRight  = 2 // 0010
	BorderBottom = 4 // 0100
	BorderLeft   = 8 // 1000
)

const (
	// Add CROSS if missing from symbols
	CROSS = '┼'
)

var borderMap = map[rune]int{
	HORIZONTAL_LINE: BorderLeft | BorderRight,
	VERTICAL_LINE:   BorderTop | BorderBottom,
	TOP_LEFT:        BorderRight | BorderBottom,
	TOP_RIGHT:       BorderLeft | BorderBottom,
	BOTTOM_LEFT:     BorderRight | BorderTop,
	BOTTOM_RIGHT:    BorderLeft | BorderTop,
	VERTICAL_RIGHT:  BorderTop | BorderBottom | BorderRight,  // ├
	VERTICAL_LEFT:   BorderTop | BorderBottom | BorderLeft,   // ┤
	HORIZONTAL_DOWN: BorderLeft | BorderRight | BorderBottom, // ┬
	HORIZONTAL_UP:   BorderLeft | BorderRight | BorderTop,    // ┴
	CROSS:           BorderTop | BorderBottom | BorderLeft | BorderRight,
	' ':             0,
}

var maskToRune = map[int]rune{
	(BorderLeft | BorderRight):                            HORIZONTAL_LINE,
	(BorderTop | BorderBottom):                            VERTICAL_LINE,
	(BorderRight | BorderBottom):                          TOP_LEFT,
	(BorderLeft | BorderBottom):                           TOP_RIGHT,
	(BorderRight | BorderTop):                             BOTTOM_LEFT,
	(BorderLeft | BorderTop):                              BOTTOM_RIGHT,
	(BorderTop | BorderBottom | BorderRight):              VERTICAL_RIGHT,
	(BorderTop | BorderBottom | BorderLeft):               VERTICAL_LEFT,
	(BorderLeft | BorderRight | BorderBottom):             HORIZONTAL_DOWN,
	(BorderLeft | BorderRight | BorderTop):                HORIZONTAL_UP,
	(BorderTop | BorderBottom | BorderLeft | BorderRight): CROSS,
	0: ' ',
}

// ResolveBorderRune merges an existing rune with a new one.
func ResolveBorderRune(existing, newRune rune) rune {
	m1, ok1 := borderMap[existing]
	m2, ok2 := borderMap[newRune]

	if ok1 && ok2 {
		combined := m1 | m2
		if r, found := maskToRune[combined]; found {
			return r
		}
	}
	return newRune
}
