package widgets

import (
	"image"
	"sync"

	rw "github.com/mattn/go-runewidth"
	ui "github.com/metaspartan/gotui/v5"
)

// Input represents a text input widget.
type Input struct {
	ui.Block
	Text        string
	TextStyle   ui.Style
	CursorStyle ui.Style
	Placeholder string
	EchoMode    EchoMode
	Cursor      int
	offset      int
	sync.Mutex
}
type EchoMode int

const (
	EchoNormal EchoMode = iota
	EchoPassword
)

// NewInput returns a new Input.
func NewInput() *Input {
	return &Input{
		Block:       *ui.NewBlock(),
		TextStyle:   ui.Theme.Paragraph.Text,
		CursorStyle: ui.NewStyle(ui.ColorBlack, ui.ColorWhite),
		EchoMode:    EchoNormal,
	}
}

// Draw draws the input to the buffer.
func (i *Input) Draw(buf *ui.Buffer) {
	i.Block.Draw(buf)
	rect := i.Inner
	width := rect.Dx()
	runes := []rune(i.Text)
	if i.EchoMode == EchoPassword {
		passwordRunes := make([]rune, len(runes))
		for j := range passwordRunes {
			passwordRunes[j] = '*'
		}
		runes = passwordRunes
	}
	if len(runes) == 0 && i.Placeholder != "" {
		i.drawPlaceholder(buf, rect)
	}
	cursorVisualX := i.calculateOffset(runes, width)
	i.drawText(buf, rect, runes)
	i.drawCursor(buf, rect, cursorVisualX)
}
func (i *Input) drawPlaceholder(buf *ui.Buffer, rect image.Rectangle) {
	buf.SetString(
		i.Placeholder,
		ui.NewStyle(ui.ColorGrey),
		image.Pt(rect.Min.X, rect.Min.Y),
	)
}
func (i *Input) calculateOffset(runes []rune, width int) int {
	if i.Cursor < 0 {
		i.Cursor = 0
	}
	if i.Cursor > len(runes) {
		i.Cursor = len(runes)
	}
	cursorVisualX := 0
	for j := 0; j < i.Cursor; j++ {
		cursorVisualX += rw.RuneWidth(runes[j])
	}
	if cursorVisualX-i.offset >= width {
		i.offset = cursorVisualX - width + 1
	}
	if cursorVisualX-i.offset < 0 {
		i.offset = cursorVisualX
	}
	totalWidth := 0
	for _, r := range runes {
		totalWidth += rw.RuneWidth(r)
	}
	if totalWidth < width {
		i.offset = 0
	}
	return cursorVisualX
}
func (i *Input) drawText(buf *ui.Buffer, rect image.Rectangle, runes []rune) {
	currentX := 0
	for _, r := range runes {
		w := rw.RuneWidth(r)
		if currentX >= i.offset {
			screenX := rect.Min.X + (currentX - i.offset)
			if screenX+w <= rect.Max.X {
				buf.SetCell(
					ui.NewCell(r, i.TextStyle),
					image.Pt(screenX, rect.Min.Y),
				)
			}
		}
		currentX += w
	}
}
func (i *Input) drawCursor(buf *ui.Buffer, rect image.Rectangle, cursorVisualX int) {
	screenCursorX := rect.Min.X + (cursorVisualX - i.offset)
	if screenCursorX >= rect.Min.X && screenCursorX < rect.Max.X {
		cell := buf.GetCell(image.Pt(screenCursorX, rect.Min.Y))
		cell.Style = i.CursorStyle
		if cell.Rune == 0 {
			cell.Rune = ' '
		}
		buf.SetCell(cell, image.Pt(screenCursorX, rect.Min.Y))
	}
}
func (i *Input) InsertRune(r rune) {
	i.Lock()
	defer i.Unlock()
	runes := []rune(i.Text)
	runes = append(runes, 0)
	copy(runes[i.Cursor+1:], runes[i.Cursor:])
	runes[i.Cursor] = r
	i.Text = string(runes)
	i.Cursor++
}
func (i *Input) Backspace() {
	i.Lock()
	defer i.Unlock()
	if i.Cursor > 0 {
		runes := []rune(i.Text)
		runes = append(runes[:i.Cursor-1], runes[i.Cursor:]...)
		i.Text = string(runes)
		i.Cursor--
	}
}
func (i *Input) MoveCursorLeft() {
	i.Lock()
	defer i.Unlock()
	if i.Cursor > 0 {
		i.Cursor--
	}
}
func (i *Input) MoveCursorRight() {
	i.Lock()
	defer i.Unlock()
	runes := []rune(i.Text)
	if i.Cursor < len(runes) {
		i.Cursor++
	}
}
