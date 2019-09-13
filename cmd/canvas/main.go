package main

import (
	"github.com/xyproto/vt100"
	//"time"
)

func main() {
	vt100.Reset()
	vt100.Clear()
	vt100.ShowCursor(false)
	vt100.SetLineWrap(false)

	c := vt100.NewCanvas()
	c.Plot(10, 10, '!')
	c.Write(12, 12, vt100.LightGreen, vt100.Default, "hi")
	c.Write(15, 15, vt100.White, vt100.Magenta, "floating")
	c.PlotC(12, 17, "Red", '*')
	c.PlotC(10, 20, "Blue", '?')

	c.Draw()
	vt100.WaitForKey()
	//time.Sleep(time.Second * 2)

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
}
