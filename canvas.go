package vt100

import (
	"fmt"
	"strconv"
)

type Char struct {
	fg     string
	bright bool
	s      string
}

type Canvas [80 * 25]Char

func NewCanvas() *Canvas {
	var c Canvas
	return &c
}

func setxy(x, y int) {
	Set("Cursor Home", map[string]string{"{ROW}": strconv.Itoa(y), "{COLUMN}": strconv.Itoa(x)})
}

func down(n int) {
	Set("Cursor Down", map[string]string{"{COUNT}": strconv.Itoa(n)})
}

func (c *Canvas) Draw() {
	// TODO: Only draw the characters that have changed since last draw
	for y := 0; y < 25; y++ {
		for x := 0; x < 80; x++ {
			ch := (*c)[y*80+x]
			if ch.s != "" {
				setxy(x, y)
				if ch.bright {
					fmt.Print(AttributeAndColor("Bright", ch.fg) + ch.s + NoColor())
				} else {
					fmt.Print(AttributeOrColor(ch.fg) + ch.s + NoColor())
				}
			}
		}
	}
	setxy(79, 24)
}

func (c *Canvas) Plot(x, y int, s string) {
	(*c)[y*80+x].s = s
}

func (c *Canvas) PlotC(x, y int, fg, s string) {
	(*c)[y*80+x].s = s
	(*c)[y*80+x].fg = fg
}
