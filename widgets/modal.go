package widgets

import (
	"image"

	rw "github.com/mattn/go-runewidth"
	ui "github.com/metaspartan/gotui/v4"
)

type Modal struct {
	ui.Block
	Text              string
	TextStyle         ui.Style
	Buttons           []*Button
	ActiveButtonIndex int
}

func NewModal(text string) *Modal {
	return &Modal{
		Block:             *ui.NewBlock(),
		Text:              text,
		TextStyle:         ui.NewStyle(ui.ColorWhite),
		Buttons:           make([]*Button, 0),
		ActiveButtonIndex: 0,
	}
}

func (m *Modal) CenterIn(x1, y1, x2, y2, width, height int) {
	totalW := x2 - x1
	totalH := y2 - y1
	if width > totalW {
		width = totalW
	}
	if height > totalH {
		height = totalH
	}

	mx := x1 + (totalW-width)/2
	my := y1 + (totalH-height)/2
	m.SetRect(mx, my, mx+width, my+height)
}

func (m *Modal) AddButton(text string, onClick func()) *Button {
	b := NewButton(text)
	b.Border = true
	m.Buttons = append(m.Buttons, b)
	return b
}

func (m *Modal) SetRect(x1, y1, x2, y2 int) {
	m.Block.SetRect(x1, y1, x2, y2)
	m.layoutButtons()
}

func (m *Modal) layoutButtons() {
	if len(m.Buttons) == 0 {
		return
	}

	buttonHeight := 3
	buttonY := m.Inner.Max.Y - buttonHeight - 1
	if buttonY < m.Inner.Min.Y {
		buttonY = m.Inner.Min.Y
	}

	totalWidth := 0
	gap := 2
	for _, b := range m.Buttons {
		w := 2 + 6 + rw.StringWidth(b.Text)
		totalWidth += w
	}
	totalWidth += (len(m.Buttons) - 1) * gap

	startX := m.Inner.Min.X + (m.Inner.Dx()-totalWidth)/2
	if startX < m.Inner.Min.X {
		startX = m.Inner.Min.X
	}

	currentX := startX
	for i, b := range m.Buttons {
		w := 2 + 6 + rw.StringWidth(b.Text)
		b.SetRect(currentX, buttonY, currentX+w, buttonY+buttonHeight)

		if i == m.ActiveButtonIndex {
			b.IsActive = true
		} else {
			b.IsActive = false
		}

		currentX += w + gap
	}
}

func (m *Modal) Draw(buf *ui.Buffer) {
	for y := m.Dy(); y < m.Dy()+m.Dy(); y++ {
	}
	for y := m.Min.Y; y < m.Max.Y; y++ {
		for x := m.Min.X; x < m.Max.X; x++ {
			buf.SetCell(ui.NewCell(' ', ui.NewStyle(ui.ColorWhite, m.BorderStyle.Bg)), image.Pt(x, y))
		}
	}
	for y := m.Min.Y; y < m.Max.Y; y++ {
		for x := m.Min.X; x < m.Max.X; x++ {
			buf.SetCell(ui.NewCell(' ', ui.NewStyle(ui.ColorWhite, m.BorderStyle.Bg)), image.Pt(x, y))
		}
	}

	m.Block.Draw(buf)

	words := splitBySpace(m.Text)
	var lines []string
	currentLine := ""
	maxWidth := m.Inner.Dx() - 2

	for _, word := range words {
		if word == "\n" {
			lines = append(lines, currentLine)
			currentLine = ""
			continue
		}

		candidate := currentLine
		if len(candidate) > 0 {
			candidate += " "
		}
		candidate += word

		if rw.StringWidth(candidate) > maxWidth {
			if len(currentLine) > 0 {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			currentLine = candidate
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	textHeight := len(lines)
	startY := m.Inner.Min.Y + (m.Inner.Dy()-textHeight-4)/2
	if startY < m.Inner.Min.Y {
		startY = m.Inner.Min.Y
	}

	for i, line := range lines {
		w := rw.StringWidth(line)
		x := m.Inner.Min.X + (m.Inner.Dx()-w)/2
		buf.SetString(line, m.TextStyle, image.Pt(x, startY+i))
	}

	m.layoutButtons()
	for _, b := range m.Buttons {
		b.Draw(buf)
	}
}

func splitBySpace(s string) []string {
	var words []string
	var currentWord []rune
	for _, r := range s {
		switch r {
		case '\n':
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
			words = append(words, "\n")
		case ' ':
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
		default:
			currentWord = append(currentWord, r)
		}
	}
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}
	return words
}
