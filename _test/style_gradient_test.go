package gotui_test

import (
	"testing"

	ui "github.com/metaspartan/gotui/v5"
)

func TestGenerateGradient(t *testing.T) {
	start := ui.ColorWhite
	end := ui.ColorBlack

	// Test positive length
	g := ui.GenerateGradient(start, end, 5)
	if len(g) != 5 {
		t.Errorf("Expected length 5, got %d", len(g))
	}

	// Test 0 length
	g = ui.GenerateGradient(start, end, 0)
	if len(g) != 0 {
		t.Errorf("Expected length 0 for 0 input, got %d", len(g))
	}

	// Test negative length
	g = ui.GenerateGradient(start, end, -5)
	if len(g) != 0 {
		t.Errorf("Expected length 0 for negative input, got %d", len(g))
	}

	// Test 1 length
	g = ui.GenerateGradient(start, end, 1)
	if len(g) != 1 {
		t.Errorf("Expected length 1, got %d", len(g))
	}
}
