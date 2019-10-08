package vt100

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
)

type Char struct {
	fg AttributeColor // Foreground color
	bg AttributeColor // Background color
	s  rune           // The character to draw
	// Not having a background color, and storing the foreground color as a string is a design choice
}

type Canvas struct {
	w             uint
	h             uint
	chars         []Char
	mut           *sync.RWMutex
	cursorVisible bool
	linesChanged  map[uint]bool
}

func NewCanvas() *Canvas {
	var err error
	c := &Canvas{}
	c.w, c.h, err = TermSize()
	if err != nil {
		// Use 80x25 if the size can't be detected
		c.w = 80
		c.h = 25
	}
	c.chars = make([]Char, c.w*c.h)
	c.mut = &sync.RWMutex{}
	c.cursorVisible = true
	c.linesChanged = make(map[uint]bool)
	return c
}

func (c *Canvas) markLineAsChangedByIndex(i uint) {
	c.linesChanged[uint(math.Floor(float64(i/c.w)))] = true
}

func (c *Canvas) markAllLinesAsChanged() {
	for y := uint(0); y < c.h; y++ {
		c.linesChanged[y] = true
	}
}

func (c *Canvas) markAllLinesAsNotChanged() {
	for y := uint(0); y < c.h; y++ {
		c.linesChanged[y] = false
	}
}

// Change the background color for each character
func (c *Canvas) FillBackground(bg AttributeColor) {
	c.mut.Lock()
	converted := bg.Background()
	for i := range c.chars {
		c.chars[i].bg = converted
	}
	c.markAllLinesAsChanged()
	c.mut.Unlock()
}

// Change the foreground color for each character
func (c *Canvas) Fill(fg AttributeColor) {
	c.mut.Lock()
	for i := range c.chars {
		c.chars[i].fg = fg
	}
	c.markAllLinesAsChanged()
	c.mut.Unlock()
}

// Bytes returns only the characters, as a long string with a newline after each row
func (c *Canvas) String() string {
	var sb strings.Builder
	for y := uint(0); y < c.h; y++ {
		c.mut.RLock()
		for x := uint(0); x < c.w; x++ {
			ch := &((*c).chars[y*c.w+x])
			if ch.s == rune(0) {
				sb.WriteRune(' ')
			} else {
				sb.WriteRune(ch.s)
			}
		}
		c.mut.RUnlock()
		sb.WriteRune('\n')
	}
	return sb.String()
}

// Return the size of the current canvas
func (c *Canvas) Size() (uint, uint) {
	return c.w, c.h
}

func (c *Canvas) Width() uint {
	return c.w
}

func (c *Canvas) Height() uint {
	return c.h
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
	c.mut.Lock()
	for _, ch := range c.chars {
		ch.s = rune(0)
	}
	c.markAllLinesAsChanged()
	c.mut.Unlock()
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

func (c *Canvas) ShowCursor() {
	if !c.cursorVisible {
		ShowCursor(true)
	}
	c.cursorVisible = true
}

func (c *Canvas) HideCursor() {
	if c.cursorVisible {
		ShowCursor(false)
	}
	c.cursorVisible = false
}

// Draw the entire canvas
func (c *Canvas) Draw() {
	c.mut.Lock()
	ShowCursor(false)
	defer func() {
		ShowCursor(c.cursorVisible)
	}()
	defer c.mut.Unlock()
	var (
		lastfg, lastbg AttributeColor
		// Build a string per line
		line strings.Builder
		ch   *Char
	)
	for y := uint(0); y < c.h; y++ {
		// if line has not changed, skip this line
		for x := uint(0); x < c.w; x++ {
			ch = &((*c).chars[y*c.w+x])
			if len(ch.bg) != 0 {
				// Write the color attributes, if they changed
				if !ch.fg.Equal(lastfg) || !ch.bg.Equal(lastbg) {
					line.WriteString(ch.fg.Combine(ch.bg).String())
				}
				lastfg = ch.fg
				lastbg = ch.bg
				if ch.s == rune(0) || len(string(ch.s)) == 0 {
					// Write a blank
					line.WriteRune(' ')
				} else {
					// Write the rune
					line.WriteRune(ch.s)
				}
			} else {
				// Write the color attributes, if they changed
				if !ch.fg.Equal(lastfg) {
					line.WriteString(ch.fg.String())
				}
				lastfg = ch.fg
				lastbg = ch.bg
				if ch.s == rune(0) || len(string(ch.s)) == 0 {
					// Write a blank
					line.WriteRune(' ')
				} else {
					// Write the rune
					line.WriteRune(ch.s)
				}
			}
		}
		c.markAllLinesAsNotChanged()
		line.WriteString(NoColor())
		SetXY(0, y)
		fmt.Print(line.String())
		line.Reset()
	}
	SetXY(c.w-1, c.h-1)
}

func (c *Canvas) Redraw() {
	c.mut.Lock()
	c.markAllLinesAsChanged()
	c.mut.Unlock()
	c.Draw()
}

// At returns the rune at the given coordinates, or an error if out of bounds
func (c *Canvas) At(x, y uint) (rune, error) {
	index := y*c.w + x
	c.mut.RLock()
	defer c.mut.RUnlock()
	chars := (*c).chars
	if index < uint(0) || index >= uint(len(chars)) {
		return rune(0), errors.New("out of bounds")
	}
	return chars[index].s, nil
}

func (c *Canvas) Plot(x, y uint, s rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	index := y*c.w + x
	c.mut.Lock()
	chars := (*c).chars
	chars[index].s = s
	c.markLineAsChangedByIndex(index)
	c.mut.Unlock()
}

func (c *Canvas) PlotColor(x, y uint, fg AttributeColor, s rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	index := y*c.w + x
	c.mut.Lock()
	chars := (*c).chars
	chars[index].s = s
	chars[index].fg = fg
	c.markLineAsChangedByIndex(index)
	c.mut.Unlock()
}

// WriteString will write a string to the canvas.
func (c *Canvas) WriteString(x, y uint, fg, bg AttributeColor, s string) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	var (
		chars   = (*c).chars
		counter uint
		index   uint
	)
	for _, r := range s {
		c.mut.Lock()
		index = y*c.w + x + counter
		chars[index].s = r
		chars[index].fg = fg
		chars[index].bg = bg.Background()
		c.markLineAsChangedByIndex(index)
		c.mut.Unlock()
		counter++
	}
}

func (c *Canvas) Write(x, y uint, fg, bg AttributeColor, s string) {
	c.WriteString(x, y, fg, bg, s)
}

// WriteRune will write a colored rune to the canvas
func (c *Canvas) WriteRune(x, y uint, fg, bg AttributeColor, r rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= c.w || y >= c.h {
		return
	}
	index := y*c.w + x
	c.mut.Lock()
	chars := (*c).chars
	chars[index].s = r
	chars[index].fg = fg
	chars[index].bg = bg.Background()
	c.markLineAsChangedByIndex(index)
	c.mut.Unlock()
}

func (c *Canvas) Resize() {
	w, h, err := TermSize()
	if err != nil {
		return
	}
	c.mut.Lock()
	if (w != c.w) || (h != c.h) {
		// Resize to the new size
		c.w = w
		c.h = h
		c.chars = make([]Char, w*h)
		c.mut = &sync.RWMutex{}
	}
	c.mut.Unlock()
}

// Check if the canvas was resized, and adjust values accordingly.
// Returns a new canvas, or nil.
func (c *Canvas) Resized() *Canvas {
	w, h, err := TermSize()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if (w != c.w) || (h != c.h) {
		// The terminal was resized!
		oldc := c

		nc := &Canvas{}
		nc.w = w
		nc.h = h
		nc.chars = make([]Char, w*h)
		nc.mut = &sync.RWMutex{}

		nc.mut.Lock()
		c.mut.Lock()
		defer c.mut.Unlock()
		defer nc.mut.Unlock()
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
				nc.chars[index] = ch
			}
		}
		c.markAllLinesAsChanged()
		// Return the new canvas
		return nc
	}
	return nil
}
