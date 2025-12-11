package gotui_test

import (
	"testing"

	ui "github.com/metaspartan/gotui/v4"
)

func TestNewRGBColor(t *testing.T) {
	c := ui.NewRGBColor(255, 0, 0)
	// tcell color logic is internal, but we can verify it returns a non-zero value or check specific mapping if we knew it.
	// For now, just ensure it doesn't panic and returns a valid Type.
	if c == ui.ColorClear {
		t.Error("NewRGBColor returned ColorClear")
	}
}

func TestDefaultStyle(t *testing.T) {
	// Verify that the default style uses ColorDefault for background to support transparency
	// We can't inspect the private Screen state directly easily without mocks,
	// but we can check if the StyleClear global is set correctly if we want.
	if ui.StyleClear.Bg != ui.ColorClear {
		t.Errorf("StyleClear.Bg should be ColorClear (Default), got %v", ui.StyleClear.Bg)
	}
}
