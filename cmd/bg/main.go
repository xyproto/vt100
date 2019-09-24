package main

import (
	"github.com/xyproto/vt100"
)

func main() {
	vt100.Init()

	c := vt100.NewCanvas()
	c.FillBackground(vt100.Blue)
	c.Draw()

	vt100.WaitForKey()

	vt100.Close()
}
