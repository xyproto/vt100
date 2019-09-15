package main

import (
	"github.com/xyproto/vt100"
)

var (
	red   = vt100.Red
	green = vt100.Green
	blue  = vt100.Blue
	none  = vt100.None
)

const (
	TL = '╭' // top left
	TR = '╮' // top right
	BL = '╰' // bottom left
	BR = '╯' // bottom right
	V  = '│' // vertical line, left side
	V2 = '│' // vertical line, right side
	H  = '─' // horizontal line
	H2 = '─' // horizontal bottom line
)

func main() {
	vt100.Init()
	defer vt100.Close()

	c := vt100.NewCanvas()

	c.WriteRune(12, 14, green, none, TL)
	c.WriteRune(13, 14, green, none, H)
	c.WriteRune(14, 14, green, none, H)
	c.WriteRune(15, 14, green, none, H)
	c.WriteRune(16, 14, green, none, H)
	c.WriteRune(17, 14, green, none, TR)

	c.WriteRune(12, 15, green, none, V)
	c.Write(14, 15, green, none, "OK")
	c.WriteRune(17, 15, green, none, V2)

	c.WriteRune(12, 16, green, none, BL)
	c.WriteRune(13, 16, green, none, H2)
	c.WriteRune(14, 16, green, none, H2)
	c.WriteRune(15, 16, green, none, H2)
	c.WriteRune(16, 16, green, none, H2)
	c.WriteRune(17, 16, green, none, BR)

	c.Draw()
	vt100.WaitForKey()
}
