package main

import (
	"github.com/xyproto/vt100"
	"time"
)

func main() {
	c := vt100.NewCanvas()
	c.Reset()
	c.Clear()
	c.SetCursor(false)
	c.SetLineWrap(false)
	c.Plot(10, 10, "!")
	c.PlotC(20, 20, "Blue", "?")
	c.Draw()
	time.Sleep(time.Second * 2)
	c.SetLineWrap(true)
	c.SetCursor(true)
}
