package main

import (
	"github.com/xyproto/vt100"
)

// RuneBox draws a blue box with yellow runes with upper left at (10,10)
// r is the rune to draw
// w and h is the width and height of the box
func RuneBox(c *vt100.Canvas, r rune, w, h int) {
	for y := uint(10); y < uint(10+h); y++ {
		c.WriteRunesB(10, y, vt100.LightYellow, vt100.BackgroundBlue, r, uint(w))
	}
}

func main() {
	// Initialize vt100 terminal settings
	vt100.Init()

	// Prepare a canvas
	c := vt100.NewCanvas()

	// Draw a box of exclamation marks on the canvas
	RuneBox(c, '!', 20, 10)

	// Draw the contents of the canvas
	c.Draw()

	// Wait for a keypress
	vt100.WaitForKey()

	// Reset the vt100 terminal settings
	vt100.Close()
}
