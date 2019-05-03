package vt100

import (
	"fmt"
	"strconv"
)

type Char struct {
	fg     string // Foreground color
	bright bool   // Bright color, or not
	s      rune   // The character to draw
	drawn  bool   // Has been drawn to screen yet?
}

type Canvas struct {
	w     uint
	h     uint
	chars []Char
}

func NewCanvas() *Canvas {
	var err error
	c := &Canvas{}
	c.w, c.h, err = TermSize()
	// TermSize is 1 too small for the buffer
	//c.w++
	//c.h++
	if err != nil {
		c.w = 80
		c.h = 25
	}
	c.chars = make([]Char, c.w*c.h)
	return c
}

// Return the size of the current canvas
func (c *Canvas) Size() (uint, uint) {
	return c.w, c.h
}

func umin(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// Move cursor to the given position (from 0 and up, the terminal code is from 1 and up)
func SetXY(x, y uint) {
	Set("Cursor Home", map[string]string{"{ROW}": strconv.Itoa(int(y + 1)), "{COLUMN}": strconv.Itoa(int(x + 1))})
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

func Home() {
	Set("Cursor Home", map[string]string{"{ROW};{COLUMN}": ""})
}

func Reset() {
	Do("Reset Device")
}

// Clear screen
func Clear() {
	Do("Erase Screen")
}

// Clear canvas
func (c *Canvas) Clear() {
	for _, ch := range c.chars {
		ch.s = rune(0)
		ch.drawn = false
	}
}

func SetLineWrap(enable bool) {
	if enable {
		Do("Enable Line Wrap")
	} else {
		Do("Disable Line Wrap")
	}
}

func ShowCursor(enable bool) {
	// Thanks https://rosettacode.org/wiki/Terminal_control/Hiding_the_cursor#Escape_code
	if enable {
		fmt.Print("\033[?25h")
	} else {
		fmt.Print("\033[?25l")
	}
}

func (c *Canvas) W() uint {
	return c.w
}

func (c *Canvas) H() uint {
	return c.h
}

func (c *Canvas) Draw() {
	// TODO: Consider using a single for-loop over index instead of 2 (x,y)
	for y := uint(0); y < c.h; y++ {
		for x := uint(0); x < c.w; x++ {
			ch := &((*c).chars[y*c.w+x])
			if !ch.drawn && ch.s != rune(0) {
				SetXY(x, y)
				if ch.bright {
					fmt.Print(AttributeAndColor("Bright", ch.fg) + string(ch.s) + NoColor())
				} else {
					fmt.Print(AttributeOrColor(ch.fg) + string(ch.s) + NoColor())
				}
				ch.drawn = true
			}
		}
	}
	SetXY(c.w-1, c.h-1)
}

func (c *Canvas) Redraw() {
	// TODO: Consider using a single for-loop instead of 1 (range) + 2 (x,y)
	for _, ch := range c.chars {
		ch.drawn = false
	}
	c.Draw()
}

func (c *Canvas) Plot(x, y uint, s rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	ch := &((*c).chars[y*c.w+x])
	ch.s = s
	ch.drawn = false
}

// Plot a bright color
func (c *Canvas) PlotC(x, y uint, fg string, s rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	ch := &((*c).chars[y*c.w+x])
	ch.s = s
	ch.fg = fg
	ch.bright = true
	ch.drawn = false
}

// Plot a dark color
func (c *Canvas) PlotDC(x, y uint, fg string, s rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	ch := &((*c).chars[y*c.w+x])
	ch.s = s
	ch.fg = fg
	ch.bright = false
	ch.drawn = false
}

func (c *Canvas) Resize() {
	w, h, err := TermSize()
	if err != nil {
		return
	}
	//w++
	//h++
	if (w != c.w) || (h != c.h) {
		// Resize to the new size
		c.w = w
		c.h = h
		c.chars = make([]Char, w*h)
	}
}

// Check if the canvas was resized, and adjust values accordingly.
// Returns a new canvas, or nil.
func (c *Canvas) Resized() *Canvas {
	w, h, err := TermSize()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//w++
	//h++
	if (w != c.w) || (h != c.h) {
		// The terminal was resized!
		oldc := c

		nc := &Canvas{}
		nc.w = w
		nc.h = h
		nc.chars = make([]Char, w*h)
	OUT:
		// Plot in the old characters
		for y := uint(0); y < umin(oldc.h, h); y++ {
			for x := uint(0); x < umin(oldc.w, w); x++ {
				oldIndex := y*oldc.w + x
				index := y*nc.w + x
				if oldIndex > index {
					break OUT
				}
				// Copy over old characters, and mark them as not drawn
				ch := oldc.chars[oldIndex]
				ch.drawn = false
				nc.chars[index] = ch
			}
		}
		// Return the new canvas
		return nc
	}
	return nil
}
