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

func (c *Canvas) Draw() {
	for y := 0; y < 25; y++ {
		Set("Cursor Home", map[string]string{"{ROW}": "0", "{COLUMN}": strconv.Itoa(y)})
		for x := 0; x < 80; x++ {
			ch := (*c)[y*80+x]
			if ch.bright {
				fmt.Print(AttributeAndColor("Bright", ch.fg) + ch.s + NoColor())
			} else {
				fmt.Print(AttributeOrColor(ch.fg) + ch.s + NoColor())
			}
		}
	}
}

func (c *Canvas) Plot(x, y int, s string) {
	(*c)[y*80+x].s = s
}

func (c *Canvas) PlotC(x, y int, fg, s string) {
	(*c)[y*80+x].s = s
	(*c)[y*80+x].fg = fg
}
