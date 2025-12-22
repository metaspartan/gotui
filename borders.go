package gotui

const (
	BoxDrawingsLightHorizontal            = '─'
	BoxDrawingsLightVertical              = '│'
	BoxDrawingsLightDownAndRight          = '┌'
	BoxDrawingsLightDownAndLeft           = '┐'
	BoxDrawingsLightUpAndRight            = '└'
	BoxDrawingsLightUpAndLeft             = '┘'
	BoxDrawingsLightDownAndHorizontal     = '┬'
	BoxDrawingsLightUpAndHorizontal       = '┴'
	BoxDrawingsLightVerticalAndRight      = '├'
	BoxDrawingsLightVerticalAndLeft       = '┤'
	BoxDrawingsLightVerticalAndHorizontal = '┼'

	BoxDrawingsLightArcDownAndRight = '╭'
	BoxDrawingsLightArcDownAndLeft  = '╮'
	BoxDrawingsLightArcUpAndRight   = '╰'
	BoxDrawingsLightArcUpAndLeft    = '╯'

	BoxDrawingsHeavyHorizontal        = '━'
	BoxDrawingsHeavyVertical          = '┃'
	BoxDrawingsHeavyDownAndRight      = '┏'
	BoxDrawingsHeavyDownAndLeft       = '┓'
	BoxDrawingsHeavyUpAndRight        = '┗'
	BoxDrawingsHeavyUpAndLeft         = '┛'
	BoxDrawingsHeavyDownAndHorizontal = '┳'
	BoxDrawingsHeavyUpAndHorizontal   = '┻'
	BoxDrawingsHeavyVerticalAndRight  = '┣'
	BoxDrawingsHeavyVerticalAndLeft   = '┫'

	BoxDrawingsDoubleHorizontal        = '═'
	BoxDrawingsDoubleVertical          = '║'
	BoxDrawingsDoubleDownAndRight      = '╔'
	BoxDrawingsDoubleDownAndLeft       = '╗'
	BoxDrawingsDoubleUpAndRight        = '╚'
	BoxDrawingsDoubleUpAndLeft         = '╝'
	BoxDrawingsDoubleDownAndHorizontal = '╦'
	BoxDrawingsDoubleUpAndHorizontal   = '╩'
	BoxDrawingsDoubleVerticalAndRight  = '╠'
	BoxDrawingsDoubleVerticalAndLeft   = '╣'
)

// BorderSet defines various borders used when primitives are drawn.
type BorderSet struct {
	Top         rune
	Bottom      rune
	Left        rune
	Right       rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	TopT        rune
	BottomT     rune
	LeftT       rune
	RightT      rune
}

func BorderSetHidden() BorderSet {
	return BorderSet{
		Top:         ' ',
		Bottom:      ' ',
		Left:        ' ',
		Right:       ' ',
		TopLeft:     ' ',
		TopRight:    ' ',
		BottomLeft:  ' ',
		BottomRight: ' ',
		TopT:        ' ',
		BottomT:     ' ',
		LeftT:       ' ',
		RightT:      ' ',
	}
}

func BorderSetPlain() BorderSet {
	return BorderSet{
		Top:         BoxDrawingsLightHorizontal,
		Bottom:      BoxDrawingsLightHorizontal,
		Left:        BoxDrawingsLightVertical,
		Right:       BoxDrawingsLightVertical,
		TopLeft:     BoxDrawingsLightDownAndRight,
		TopRight:    BoxDrawingsLightDownAndLeft,
		BottomLeft:  BoxDrawingsLightUpAndRight,
		BottomRight: BoxDrawingsLightUpAndLeft,
		TopT:        BoxDrawingsLightDownAndHorizontal,
		BottomT:     BoxDrawingsLightUpAndHorizontal,
		LeftT:       BoxDrawingsLightVerticalAndRight,
		RightT:      BoxDrawingsLightVerticalAndLeft,
	}
}

func BorderSetRound() BorderSet {
	return BorderSet{
		Top:         BoxDrawingsLightHorizontal,
		Bottom:      BoxDrawingsLightHorizontal,
		Left:        BoxDrawingsLightVertical,
		Right:       BoxDrawingsLightVertical,
		TopLeft:     BoxDrawingsLightArcDownAndRight,
		TopRight:    BoxDrawingsLightArcDownAndLeft,
		BottomLeft:  BoxDrawingsLightArcUpAndRight,
		BottomRight: BoxDrawingsLightArcUpAndLeft,
		TopT:        BoxDrawingsLightDownAndHorizontal,
		BottomT:     BoxDrawingsLightUpAndHorizontal,
		LeftT:       BoxDrawingsLightVerticalAndRight,
		RightT:      BoxDrawingsLightVerticalAndLeft,
	}
}

func BorderSetThick() BorderSet {
	return BorderSet{
		Top:         BoxDrawingsHeavyHorizontal,
		Bottom:      BoxDrawingsHeavyHorizontal,
		Left:        BoxDrawingsHeavyVertical,
		Right:       BoxDrawingsHeavyVertical,
		TopLeft:     BoxDrawingsHeavyDownAndRight,
		TopRight:    BoxDrawingsHeavyDownAndLeft,
		BottomLeft:  BoxDrawingsHeavyUpAndRight,
		BottomRight: BoxDrawingsHeavyUpAndLeft,
		TopT:        BoxDrawingsHeavyDownAndHorizontal,
		BottomT:     BoxDrawingsHeavyUpAndHorizontal,
		LeftT:       BoxDrawingsHeavyVerticalAndRight,
		RightT:      BoxDrawingsHeavyVerticalAndLeft,
	}
}

func BorderSetDouble() BorderSet {
	return BorderSet{
		Top:         BoxDrawingsDoubleHorizontal,
		Bottom:      BoxDrawingsDoubleHorizontal,
		Left:        BoxDrawingsDoubleVertical,
		Right:       BoxDrawingsDoubleVertical,
		TopLeft:     BoxDrawingsDoubleDownAndRight,
		TopRight:    BoxDrawingsDoubleDownAndLeft,
		BottomLeft:  BoxDrawingsDoubleUpAndRight,
		BottomRight: BoxDrawingsDoubleUpAndLeft,
		TopT:        BoxDrawingsDoubleDownAndHorizontal,
		BottomT:     BoxDrawingsDoubleUpAndHorizontal,
		LeftT:       BoxDrawingsDoubleVerticalAndRight,
		RightT:      BoxDrawingsDoubleVerticalAndLeft,
	}
}

type Borders uint

const (
	BordersTop Borders = 1 << iota
	BordersBottom
	BordersLeft
	BordersRight

	BordersNone Borders = 0
	BordersAll  Borders = BordersTop | BordersBottom | BordersLeft | BordersRight
)

func (b Borders) Has(flag Borders) bool {
	return b&flag != 0
}
