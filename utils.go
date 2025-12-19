package gotui

import (
	"fmt"
	"math"
	"reflect"

	rw "github.com/mattn/go-runewidth"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

// InterfaceSlice converts a slice of any type to a slice of interface{}.
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	return ret
}

// TrimString trims a string to a given width.
func TrimString(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if rw.StringWidth(s) > w {
		return rw.Truncate(s, w, string(ELLIPSES))
	}
	return s
}

// SelectColor selects a color from a slice of colors based on an index.
func SelectColor(colors []Color, index int) Color {
	return colors[index%len(colors)]
}

// SelectStyle selects a style from a slice of styles based on an index.
func SelectStyle(styles []Style, index int) Style {
	return styles[index%len(styles)]
}

// SumIntSlice sums a slice of ints.
func SumIntSlice(slice []int) int {
	sum := 0
	for _, val := range slice {
		sum += val
	}
	return sum
}

// SumFloat64Slice sums a slice of float64s.
func SumFloat64Slice(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

// GetMaxIntFromSlice returns the maximum value from a slice of ints.
func GetMaxIntFromSlice(slice []int) (int, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	var max int
	for _, val := range slice {
		if val > max {
			max = val
		}
	}
	return max, nil
}

// GetMaxFloat64FromSlice returns the maximum value from a slice of float64s.
func GetMaxFloat64FromSlice(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	var max float64
	for _, val := range slice {
		if val > max {
			max = val
		}
	}
	return max, nil
}

// GetMaxFloat64From2dSlice returns the maximum value from a 2D slice of float64s.
func GetMaxFloat64From2dSlice(slices [][]float64) (float64, error) {
	if len(slices) == 0 {
		return 0, fmt.Errorf("cannot get max value from empty slice")
	}
	var max float64
	for _, slice := range slices {
		for _, val := range slice {
			if val > max {
				max = val
			}
		}
	}
	return max, nil
}

// RoundFloat64 rounds a float64 to the nearest integer.
func RoundFloat64(x float64) float64 {
	return math.Floor(x + 0.5)
}

// FloorFloat64 floors a float64.
func FloorFloat64(x float64) float64 {
	return math.Floor(x)
}

// AbsInt returns the absolute value of an int.
func AbsInt(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// MinFloat64 returns the minimum of two float64s.
func MinFloat64(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// MaxFloat64 returns the maximum of two float64s.
func MaxFloat64(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

// MaxInt returns the maximum of two ints.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the minimum of two ints.
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// WrapCells wraps a slice of cells to a given width.
func WrapCells(cells []Cell, width uint) []Cell {
	str := CellsToString(cells)
	wrapped := wordwrap.WrapString(str, width)
	wrappedCells := []Cell{}
	i := 0
	for _, _rune := range wrapped {
		if _rune == '\n' {
			wrappedCells = append(wrappedCells, Cell{_rune, StyleClear})
		} else {
			wrappedCells = append(wrappedCells, Cell{_rune, cells[i].Style})
		}
		i++
	}
	return wrappedCells
}

// RunesToStyledCells converts a slice of runes to a slice of cells with a given style.
func RunesToStyledCells(runes []rune, style Style) []Cell {
	cells := []Cell{}
	for _, _rune := range runes {
		cells = append(cells, Cell{_rune, style})
	}
	return cells
}

// CellsToString converts a slice of cells to a string.
func CellsToString(cells []Cell) string {
	runes := make([]rune, len(cells))
	for i, cell := range cells {
		runes[i] = cell.Rune
	}
	return string(runes)
}

// TrimCells trims a slice of cells to a given width.
func TrimCells(cells []Cell, w int) []Cell {
	s := CellsToString(cells)
	s = TrimString(s, w)
	runes := []rune(s)
	newCells := []Cell{}
	for i, r := range runes {
		newCells = append(newCells, Cell{r, cells[i].Style})
	}
	return newCells
}

// SplitCells splits a slice of cells by a rune.
func SplitCells(cells []Cell, r rune) [][]Cell {
	splitCells := [][]Cell{}
	temp := []Cell{}
	for _, cell := range cells {
		if cell.Rune == r {
			splitCells = append(splitCells, temp)
			temp = []Cell{}
		} else {
			temp = append(temp, cell)
		}
	}
	if len(temp) > 0 {
		splitCells = append(splitCells, temp)
	}
	return splitCells
}

type CellWithX struct {
	X    int
	Cell Cell
}

// BuildCellWithXArray builds an array of CellWithX from a slice of Cells.
func BuildCellWithXArray(cells []Cell) []CellWithX {
	cellWithXArray := make([]CellWithX, len(cells))
	index := 0
	for i, cell := range cells {
		cellWithXArray[i] = CellWithX{X: index, Cell: cell}
		index += rw.RuneWidth(cell.Rune)
	}
	return cellWithXArray
}
