package gotui_test

import (
	"image"
	"testing"

	ui "github.com/metaspartan/gotui/v4"
)

func TestGridLayout(t *testing.T) {
	// 3 items 1/3 each
	g := ui.NewGrid()
	width, height := 10, 10
	g.SetRect(0, 0, width, height)

	b1 := ui.NewBlock()
	b2 := ui.NewBlock()
	b3 := ui.NewBlock()

	g.Set(
		ui.NewRow(1.0/3.0, b1),
		ui.NewRow(1.0/3.0, b2),
		ui.NewRow(1.0/3.0, b3),
	)

	buf := ui.NewBuffer(image.Rect(0, 0, width, height))
	g.Draw(buf)

	r1 := b1.GetRect()
	r2 := b2.GetRect()
	r3 := b3.GetRect()

	// Expect strict adjacency
	if r1.Max.Y != r2.Min.Y {
		t.Errorf("Gap between 1 and 2: %d != %d", r1.Max.Y, r2.Min.Y)
	}
	if r2.Max.Y != r3.Min.Y {
		t.Errorf("Gap between 2 and 3: %d != %d", r2.Max.Y, r3.Min.Y)
	}
	if r3.Max.Y != height {
		t.Errorf("Last item does not end at grid bottom: %d != %d", r3.Max.Y, height)
	}
}

func TestGridLayoutFloats(t *testing.T) {
	// Test splitting 100 into 3 parts (33.333...)
	g := ui.NewGrid()
	width, height := 100, 100
	g.SetRect(0, 0, width, height)

	b1 := ui.NewBlock()
	b2 := ui.NewBlock()
	b3 := ui.NewBlock()

	g.Set(
		ui.NewRow(1.0/3.0, b1),
		ui.NewRow(1.0/3.0, b2),
		ui.NewRow(1.0/3.0, b3),
	)

	buf := ui.NewBuffer(image.Rect(0, 0, width, height))
	g.Draw(buf)

	r1 := b1.GetRect()
	r2 := b2.GetRect()
	r3 := b3.GetRect()

	if r1.Max.Y != r2.Min.Y {
		t.Errorf("Gap 1-2: %d != %d", r1.Max.Y, r2.Min.Y)
	}
	// r1 height 33
	// r2 height 33
	// r3 height 34

	// 33+33+34 = 100.

	if r3.Max.Y != height {
		t.Errorf("Bottom alignment: %d != %d", r3.Max.Y, height)
	}
}

func TestGridNested(t *testing.T) {
	// Nested grid: Row with 2 cols, one col has 2 rows.
	g := ui.NewGrid()
	width, height := 100, 100
	g.SetRect(0, 0, width, height)

	b1 := ui.NewBlock() // Left Col (50%)
	b2 := ui.NewBlock() // Right Col Top (50% of 50% = 25% height? No, 50% width, split vertically?)
	b3 := ui.NewBlock() // Right Col Bottom

	// Row
	//   Col 1/2 (b1)
	//   Col 1/2
	//     Row 1/2 (b2)
	//     Row 1/2 (b3)

	g.Set(
		ui.NewRow(1.0,
			ui.NewCol(0.5, b1),
			ui.NewCol(0.5,
				ui.NewRow(0.5, b2),
				ui.NewRow(0.5, b3),
			),
		),
	)

	buf := ui.NewBuffer(image.Rect(0, 0, width, height))
	g.Draw(buf)

	r1 := b1.GetRect()
	r2 := b2.GetRect()
	r3 := b3.GetRect()

	// b1 should be 0,0 to 50,100
	if r1.Max.X != 50 || r1.Max.Y != 100 {
		t.Errorf("b1 rect mismatch: %v", r1)
	}

	// b2 should be 50,0 to 100, 50
	if r2.Min.X != 50 || r2.Min.Y != 0 {
		t.Errorf("b2 min mismatch: %v", r2.Min)
	}
	if r2.Max.X != 100 || r2.Max.Y != 50 {
		t.Errorf("b2 max mismatch: %v", r2.Max)
	}

	// b3 should be 50,50 to 100,100
	if r3.Min.X != 50 || r3.Min.Y != 50 {
		t.Errorf("b3 min mismatch: %v", r3.Min)
	}
	if r3.Max.X != 100 || r3.Max.Y != 100 {
		t.Errorf("b3 max mismatch: %v", r3.Max)
	}
}
