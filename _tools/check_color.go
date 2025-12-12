package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func main() {
	c := tcell.ColorRed
	r, g, b := c.RGB()
	fmt.Printf("Red: %d,%d,%d\n", r, g, b)

	c2 := tcell.Color(100)
	r2, g2, b2 := c2.RGB()
	fmt.Printf("Color(100): %d,%d,%d\n", r2, g2, b2)
}
