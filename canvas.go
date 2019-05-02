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

type Canvas struct {
	w     uint
	h     uint
	chars []Char
}

func NewCanvas(w, h int) *Canvas {
	var err error
	c := &Canvas{}
	c.w, c.h, err = TermSize()
	if err != nil {
		c.w = 80
		c.h = 25
	}
	c.chars = make([]Char, c.w*c.h)
	return c
}

// Move cursor to the given position
func SetXY(x, y uint) {
	Set("Cursor Home", map[string]string{"{ROW}": strconv.Itoa(int(y)), "{COLUMN}": strconv.Itoa(int(x))})
}

// Move the cursor down
func Down(n uint) {
	Set("Cursor Down", map[string]string{"{COUNT}": strconv.Itoa(int(n))})
}

// Move the cursor up
func Up(n uint) {
	Set("Cursor Up", map[string]string{"{COUNT}": strconv.Itoa(int(n))})
}

// Move the cursor to the right
func Right(n uint) {
	Set("Cursor Forward", map[string]string{"{COUNT}": strconv.Itoa(int(n))})
}

// Move the cursor to the left
func Left(n uint) {
	Set("Cursor Backward", map[string]string{"{COUNT}": strconv.Itoa(int(n))})
}

func (c *Canvas) Reset() {
	Do("Reset Device")
}

func (c *Canvas) Clear() {
	Do("Erase Screen")
}

func (c *Canvas) SetLineWrap(enable bool) {
	if enable {
		Do("Enable Line Wrap")
	} else {
		Do("Disable Line Wrap")
	}
}

func (c *Canvas) SetCursor(enable bool) {
	// Thanks https://rosettacode.org/wiki/Terminal_control/Hiding_the_cursor#Escape_code
	if enable {
		fmt.Print("\033[?25h")
	} else {
		fmt.Print("\033[?25l")
	}
}

func (c *Canvas) Draw() {
	// TODO: Only draw the characters that have changed since last draw
	for y := uint(0); y < c.h; y++ {
		for x := uint(0); x < c.w; x++ {
			ch := (*c).chars[y*c.w+x]
			if ch.s != "" {
				SetXY(x, y)
				if ch.bright {
					fmt.Print(AttributeAndColor("Bright", ch.fg) + ch.s + NoColor())
				} else {
					fmt.Print(AttributeOrColor(ch.fg) + ch.s + NoColor())
				}
			}
		}
	}
	SetXY(c.w-1, c.h-1)
}

func (c *Canvas) Plot(x, y uint, s string) {
	(*c).chars[y*c.w+x].s = s
}

func (c *Canvas) PlotC(x, y uint, fg, s string) {
	(*c).chars[y*c.w+x].s = s
	(*c).chars[y*c.w+x].fg = fg
}
