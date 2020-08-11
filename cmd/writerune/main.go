package main

import (
	"github.com/xyproto/vt100"
)

func main() {
	// Initialize vt100 terminal settings
	vt100.Init()

	// Prepare a canvas
	c := vt100.NewCanvas()

	// Draw things on the canvas
	c.Plot(0, 0, '!')
	c.Plot(0, 1, '!')
	c.Plot(1, 0, '!')
	c.Plot(1, 1, '!')

	c.WriteString(9, 20, vt100.Red, vt100.BackgroundBlue, "# Thank you code_nomad: http://9m.no/ꪯ鵞")

	bg := vt100.BackgroundBlue
	fg := vt100.LightYellow

	// Draw the contents of the canvas
	c.Draw()

	// Wait for a keypress
	vt100.WaitForKey()

	c.WriteRuneB(0, 0, fg, bg, 'A')
	c.WriteRuneB(0, 1, fg, bg, 'B')
	c.WriteRuneB(1, 0, fg, bg, 'C')
	c.WriteRuneB(1, 1, fg, bg, 'D')

	// Draw the contents of the canvas
	c.Draw()

	// Wait for a keypress
	vt100.WaitForKey()

	// Reset the vt100 terminal settings
	vt100.Close()
}
