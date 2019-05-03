package main

import (
	"github.com/xyproto/vt100"
	"time"
)

func main() {
	vt100.Reset()
	vt100.Clear()
	vt100.ShowCursor(false)
	vt100.SetLineWrap(false)

	c := vt100.NewCanvas()
	c.Plot(10, 10, '!')
	c.PlotC(20, 20, "Blue", '?')
	c.Draw()

	time.Sleep(time.Second * 2)

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
}
