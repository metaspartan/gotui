package widgets

import (
	"fmt"
	"image"

	ui "github.com/metaspartan/gotui/v4"
)

// Standard spinner styles
var (
	SpinnerLine           = []string{"|", "/", "-", "\\"}
	SpinnerDots           = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SpinnerMiniDots       = []string{"⠋", "⠙", "⠚", "⠞", "⠖", "⠦", "⠴", "⠲", "⠳", "⠓"}
	SpinnerPulse          = []string{"█", "▓", "▒", "░"}
	SpinnerPoints         = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}
	SpinnerGlobe          = []string{"🌍", "🌎", "🌏"}
	SpinnerMoon           = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	SpinnerClock          = []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
	SpinnerMonkey         = []string{"🙈", "🙉", "🙊"}
	SpinnerStar           = []string{"✶", "✸", "✹", "✺", "✹", "✸"}
	SpinnerHamburger      = []string{"☱", "☲", "☴"}
	SpinnerGrowVertical   = []string{" ", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"}
	SpinnerGrowHorizontal = []string{"▉", "▊", "▋", "▌", "▍", "▎", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}
	SpinnerArrow          = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	SpinnerTriangle       = []string{"◢", "◣", "◤", "◥"}
	SpinnerCircleHalves   = []string{"◐", "◓", "◑", "◒"}
	SpinnerBouncingBall   = []string{"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"}
)

type Spinner struct {
	ui.Block
	Frames       []string
	Index        int
	Label        string
	LabelOnRight bool
	FormatString string
	TextStyle    ui.Style
}

func NewSpinner() *Spinner {
	return &Spinner{
		Block:        *ui.NewBlock(),
		Frames:       SpinnerLine,
		Index:        0,
		FormatString: "%s [%s]",
		TextStyle:    ui.NewStyle(ui.ColorWhite),
	}
}

func (s *Spinner) Advance() {
	s.Index = (s.Index + 1) % len(s.Frames)
}

func (s *Spinner) Draw(buf *ui.Buffer) {
	s.Block.Draw(buf)

	if len(s.Frames) == 0 {
		return
	}

	symbol := fmt.Sprintf(s.FormatString, s.Frames[s.Index], "")
	if len(s.Label) > 0 {
		if s.LabelOnRight {
			symbol = fmt.Sprintf(s.FormatString, s.Frames[s.Index], s.Label)
		} else {
			symbol = fmt.Sprintf(s.FormatString, s.Label, s.Frames[s.Index])
		}
	}
	x := s.Inner.Min.X
	y := s.Inner.Min.Y

	buf.SetString(symbol, s.TextStyle, image.Pt(x, y))
}
