package widgets

import (
	"image"
	"strings"
	"sync"

	rw "github.com/mattn/go-runewidth"
	ui "github.com/metaspartan/gotui/v4"
)

// TextArea represents a widget that displays a text area.
type TextArea struct {
	ui.Block
	Text        string
	TextStyle   ui.Style
	CursorStyle ui.Style
	Cursor      image.Point
	ShowCursor  bool
	topLine     int
	leftCol     int
	sync.Mutex
}

// NewTextArea returns a new TextArea.
func NewTextArea() *TextArea {
	return &TextArea{
		Block:       *ui.NewBlock(),
		TextStyle:   ui.Theme.Paragraph.Text,
		CursorStyle: ui.NewStyle(ui.ColorBlack, ui.ColorWhite),
		ShowCursor:  true,
		Cursor:      image.Point{0, 0},
	}
}

// Draw draws the text area to the buffer.
func (ta *TextArea) Draw(buf *ui.Buffer) {
	ta.Block.Draw(buf)
	lines := strings.Split(ta.Text, "\n")
	innerRect := ta.Inner
	height := innerRect.Dy()
	width := innerRect.Dx()
	ta.adjustScroll(lines, height)
	ta.drawText(buf, lines, width, height)
	ta.drawCursor(buf, lines, width, height)
}
func (ta *TextArea) adjustScroll(lines []string, height int) {
	if ta.Cursor.Y < ta.topLine {
		ta.topLine = ta.Cursor.Y
	}
	if ta.Cursor.Y >= ta.topLine+height {
		ta.topLine = ta.Cursor.Y - height + 1
	}
}
func (ta *TextArea) drawText(buf *ui.Buffer, lines []string, width, height int) {
	innerRect := ta.Inner
	for y := 0; y < height; y++ {
		lineIdx := ta.topLine + y
		if lineIdx >= len(lines) {
			break
		}
		line := lines[lineIdx]
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
}
func (ta *TextArea) drawCursor(buf *ui.Buffer, lines []string, width, height int) {
	if ta.ShowCursor {
		cursorY := ta.Cursor.Y - ta.topLine
		cursorX := 0
		if ta.Cursor.Y < len(lines) {
			line := []rune(lines[ta.Cursor.Y])
			for i := 0; i < ta.Cursor.X && i < len(line); i++ {
				cursorX += rw.RuneWidth(line[i])
			}
		}
		innerRect := ta.Inner
		if cursorY >= 0 && cursorY < height && cursorX >= 0 && cursorX < width {
			p := image.Pt(innerRect.Min.X+cursorX, innerRect.Min.Y+cursorY)
			cell := buf.GetCell(p)
			if cell.Rune == 0 {
				cell.Rune = ' '
			}
			cell.Style = ta.CursorStyle
			buf.SetCell(cell, p)
		}
	}
}

// MoveCursor moves the cursor by the given amount.
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
	lineLen := 0
	if newY < len(lines) {
		lineLen = len([]rune(lines[newY]))
	}
	if newX < 0 {
		newX = 0
	}
	if newX > lineLen {
		newX = lineLen
	}
	ta.Cursor = image.Point{newX, newY}
}

// InsertRune inserts a rune at the cursor position.
func (ta *TextArea) InsertRune(r rune) {
	ta.Lock()
	defer ta.Unlock()
	lines := strings.Split(ta.Text, "\n")
	if ta.Cursor.Y >= len(lines) {
		if len(lines) == 0 {
			lines = []string{""}
		}
	}
	line := []rune(lines[ta.Cursor.Y])
	newLine := make([]rune, len(line)+1)
	copy(newLine, line[:ta.Cursor.X])
	newLine[ta.Cursor.X] = r
	copy(newLine[ta.Cursor.X+1:], line[ta.Cursor.X:])
	lines[ta.Cursor.Y] = string(newLine)
	ta.Text = strings.Join(lines, "\n")
	ta.Cursor.X++
}

// InsertNewline inserts a newline at the cursor position.
func (ta *TextArea) InsertNewline() {
	ta.Lock()
	defer ta.Unlock()
	lines := strings.Split(ta.Text, "\n")
	line := []rune(lines[ta.Cursor.Y])
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

// DeleteRune deletes the rune at the cursor position.
func (ta *TextArea) DeleteRune() {
	ta.Lock()
	defer ta.Unlock()
	if ta.Cursor.X == 0 && ta.Cursor.Y == 0 {
		return
	}
	lines := strings.Split(ta.Text, "\n")
	if ta.Cursor.X > 0 {
		line := []rune(lines[ta.Cursor.Y])
		newLine := make([]rune, len(line)-1)
		copy(newLine, line[:ta.Cursor.X-1])
		copy(newLine[ta.Cursor.X-1:], line[ta.Cursor.X:])
		lines[ta.Cursor.Y] = string(newLine)
		ta.Cursor.X--
	} else {
		prevLineIdx := ta.Cursor.Y - 1
		currentLine := lines[ta.Cursor.Y]
		prevLine := lines[prevLineIdx]
		newCursorX := len([]rune(prevLine))
		lines[prevLineIdx] = prevLine + currentLine
		lines = append(lines[:ta.Cursor.Y], lines[ta.Cursor.Y+1:]...)
		ta.Cursor.Y--
		ta.Cursor.X = newCursorX
	}
	ta.Text = strings.Join(lines, "\n")
}
