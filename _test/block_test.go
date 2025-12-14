package gotui_test

import (
	"testing"

	ui "github.com/metaspartan/gotui/v4"
)

func TestBlockSetRect(t *testing.T) {
	b := ui.NewBlock()
	b.SetRect(0, 0, 10, 10)

	if b.Inner.Min.X != 1 || b.Inner.Min.Y != 1 {
		t.Errorf("Expected Inner.Min (1,1), got %v", b.Inner.Min)
	}
	if b.Inner.Max.X != 9 || b.Inner.Max.Y != 9 {
		t.Errorf("Expected Inner.Max (9,9), got %v", b.Inner.Max)
	}
}

func TestBlockPaddingOverflow(t *testing.T) {
	b := ui.NewBlock()
	b.PaddingLeft = 10
	b.PaddingRight = 10
	b.PaddingTop = 10
	b.PaddingBottom = 10

	b.SetRect(0, 0, 10, 10)

	if b.Inner.Min.X != 5 || b.Inner.Min.Y != 5 {
		t.Errorf("Expected Inner.Min clamped to (5,5), got %v", b.Inner.Min)
	}
	if b.Inner.Max.X != 5 || b.Inner.Max.Y != 5 {
		t.Errorf("Expected Inner.Max clamped to (5,5), got %v", b.Inner.Max)
	}
}

func TestBlockZeroSize(t *testing.T) {
	b := ui.NewBlock()
	b.Border = true
	b.SetRect(0, 0, 0, 0)
	if b.Inner.Min.X != 0 || b.Inner.Min.Y != 0 {
		t.Errorf("Expected Inner.Min (0,0) for 0x0 block, got %v", b.Inner.Min)
	}
	if b.Inner.Max.X != 0 || b.Inner.Max.Y != 0 {
		t.Errorf("Expected Inner.Max (0,0) for 0x0 block, got %v", b.Inner.Max)
	}

	b.SetRect(0, 0, 1, 1)
	if b.Inner.Min.X != 0 || b.Inner.Min.Y != 0 {
		t.Errorf("Expected Inner.Min (0,0) for 1x1 block, got %v", b.Inner.Min)
	}
}
