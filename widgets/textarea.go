package widgets

import (
	"image"
	"strings"
	"sync"

	rw "github.com/mattn/go-runewidth"
	ui "github.com/metaspartan/gotui"
)

type TextArea struct {
	ui.Block
	Text        string
	TextStyle   ui.Style
	CursorStyle ui.Style
	Cursor      image.Point // Cursor position in the text (Column, Line) 0-indexed
	ShowCursor  bool

	// Internal scroll offset
	topLine int
	leftCol int

	sync.Mutex
}

func NewTextArea() *TextArea {
	return &TextArea{
		Block:       *ui.NewBlock(),
		TextStyle:   ui.Theme.Paragraph.Text,
		CursorStyle: ui.NewStyle(ui.ColorBlack, ui.ColorWhite),
		ShowCursor:  true,
		Cursor:      image.Point{0, 0},
	}
}

func (ta *TextArea) Draw(buf *ui.Buffer) {
	ta.Block.Draw(buf)

	lines := strings.Split(ta.Text, "\n")

	// Calculate visible area
	innerRect := ta.Inner
	height := innerRect.Dy()
	width := innerRect.Dx()

	// Adjust scroll to keep cursor in view
	if ta.Cursor.Y < ta.topLine {
		ta.topLine = ta.Cursor.Y
	}
	if ta.Cursor.Y >= ta.topLine+height {
		ta.topLine = ta.Cursor.Y - height + 1
	}
	// TODO: Horizontal scrolling

	for y := 0; y < height; y++ {
		lineIdx := ta.topLine + y
		if lineIdx >= len(lines) {
			break
		}

		line := lines[lineIdx]
		// Convert to runes for safe handling
		runes := []rune(line)

		x := 0
		for _, r := range runes {
			if x >= width {
				break
			}
			w := rw.RuneWidth(r)
			if x+w > width {
				break
			}

			buf.SetCell(
				ui.NewCell(r, ta.TextStyle),
				image.Pt(innerRect.Min.X+x, innerRect.Min.Y+y),
			)
			x += w
		}
	}

	// Draw Cursor
	if ta.ShowCursor {
		cursorY := ta.Cursor.Y - ta.topLine
		// Calculate cursor X based on rune widths of the current line
		cursorX := 0
		if ta.Cursor.Y < len(lines) {
			line := []rune(lines[ta.Cursor.Y])
			for i := 0; i < ta.Cursor.X && i < len(line); i++ {
				cursorX += rw.RuneWidth(line[i])
			}
		}

		// If cursor is beyond line length (ghost cursor), explicitly handle it?
		// Simpler: Just rely on logical X if we enforce it within bounds.

		if cursorY >= 0 && cursorY < height && cursorX >= 0 && cursorX < width {
			// Get cell under cursor
			p := image.Pt(innerRect.Min.X+cursorX, innerRect.Min.Y+cursorY)
			cell := buf.GetCell(p)
			if cell.Rune == 0 {
				cell.Rune = ' '
			}
			// Apply cursor style
			cell.Style = ta.CursorStyle
			buf.SetCell(cell, p)
		}
	}
}

// MoveCursor moves the cursor safely
func (ta *TextArea) MoveCursor(dx, dy int) {
	ta.Lock()
	defer ta.Unlock()

	lines := strings.Split(ta.Text, "\n")

	newX := ta.Cursor.X + dx
	newY := ta.Cursor.Y + dy

	if newY < 0 {
		newY = 0
	}
	if newY >= len(lines) {
		newY = len(lines) - 1
	}

	// Clamp X to line length
	lineLen := 0
	if newY < len(lines) {
		lineLen = len([]rune(lines[newY]))
	}

	if newX < 0 {
		newX = 0 // Wrap to prev line? For now just clamp.
	}
	if newX > lineLen {
		newX = lineLen
	}

	ta.Cursor = image.Point{newX, newY}
}

// InsertRune inserts a rune at the current cursor position
func (ta *TextArea) InsertRune(r rune) {
	ta.Lock()
	defer ta.Unlock()

	lines := strings.Split(ta.Text, "\n")
	if ta.Cursor.Y >= len(lines) {
		// Should not happen unless empty
		if len(lines) == 0 {
			lines = []string{""}
		}
	}

	line := []rune(lines[ta.Cursor.Y])

	// Insert
	newLine := make([]rune, len(line)+1)
	copy(newLine, line[:ta.Cursor.X])
	newLine[ta.Cursor.X] = r
	copy(newLine[ta.Cursor.X+1:], line[ta.Cursor.X:])

	lines[ta.Cursor.Y] = string(newLine)
	ta.Text = strings.Join(lines, "\n")
	ta.Cursor.X++
}

// InsertNewline inserts a newline at the cursor
func (ta *TextArea) InsertNewline() {
	ta.Lock()
	defer ta.Unlock()

	lines := strings.Split(ta.Text, "\n")
	line := []rune(lines[ta.Cursor.Y])

	// Split line
	left := string(line[:ta.Cursor.X])
	right := string(line[ta.Cursor.X:])

	newLines := make([]string, len(lines)+1)
	copy(newLines, lines[:ta.Cursor.Y])
	newLines[ta.Cursor.Y] = left
	newLines[ta.Cursor.Y+1] = right
	copy(newLines[ta.Cursor.Y+2:], lines[ta.Cursor.Y+1:])

	ta.Text = strings.Join(newLines, "\n")
	ta.Cursor.Y++
	ta.Cursor.X = 0
}

// DeleteRune deletes the rune before the cursor (backspace)
func (ta *TextArea) DeleteRune() {
	ta.Lock()
	defer ta.Unlock()

	if ta.Cursor.X == 0 && ta.Cursor.Y == 0 {
		return
	}

	lines := strings.Split(ta.Text, "\n")

	if ta.Cursor.X > 0 {
		// Simple delete char
		line := []rune(lines[ta.Cursor.Y])
		newLine := make([]rune, len(line)-1)
		copy(newLine, line[:ta.Cursor.X-1])
		copy(newLine[ta.Cursor.X-1:], line[ta.Cursor.X:])
		lines[ta.Cursor.Y] = string(newLine)
		ta.Cursor.X--
	} else {
		// Merge with previous line
		prevLineIdx := ta.Cursor.Y - 1
		currentLine := lines[ta.Cursor.Y]
		prevLine := lines[prevLineIdx]

		newCursorX := len([]rune(prevLine))

		lines[prevLineIdx] = prevLine + currentLine
		// Remove current line
		lines = append(lines[:ta.Cursor.Y], lines[ta.Cursor.Y+1:]...)

		ta.Cursor.Y--
		ta.Cursor.X = newCursorX
	}
	ta.Text = strings.Join(lines, "\n")
}
