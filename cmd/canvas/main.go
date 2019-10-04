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
	c.Plot(10, 10, '!')
	c.Write(12, 12, vt100.LightGreen, vt100.BackgroundDefault, "hi")
	c.Write(15, 15, vt100.White, vt100.BackgroundMagenta, "floating")
	c.PlotColor(12, 17, vt100.LightRed, '*')
	c.PlotColor(10, 20, vt100.LightBlue, 'ø')
	c.PlotColor(11, 20, vt100.LightBlue, 'l')

	c.WriteString(10, 21, vt100.White, vt100.BackgroundRed, "øl")

	// Draw the contents of the canvas
	c.Draw()

	// Wait for a keypress
	vt100.WaitForKey()

	// Reset the vt100 terminal settings
	vt100.Close()
}
