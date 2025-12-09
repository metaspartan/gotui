package gotui_test

import (
	"testing"

	ui "github.com/metaspartan/gotui"
)

func TestBlock_BottomTitle(t *testing.T) {
	b := ui.NewBlock()
	b.Title = "Top"
	b.TitleBottom = "Bottom"

	if b.Title != "Top" {
		t.Errorf("Expected Title 'Top', got '%s'", b.Title)
	}
	if b.TitleBottom != "Bottom" {
		t.Errorf("Expected TitleBottom 'Bottom', got '%s'", b.TitleBottom)
	}
}
