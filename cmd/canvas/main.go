package main

import (
	"github.com/xyproto/vt100"
)

func main() {
	c := vt100.NewCanvas()
	c.Plot(10, 10, "!")
	c.PlotC(20, 20, "Blue", "?")
	c.Draw()
}
