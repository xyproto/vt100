package main

import (
	"github.com/xyproto/vt100"
)

// Clear the canvas, set a background color and draw the canvas.
func (t *Theme) DrawBackground() {
	c := vt100.NewCanvas()
	c.FillBackground(t.Background)
	c.Clear()
	c.Draw()
}

// Draw a box using ASCII graphics.
// The given Box struct defines the size and placement.
// If extrude is True, the box looks a bit more like it's sticking out.
func (t *Theme) DrawBox(c *vt100.Canvas, r *Box, extrude bool) *Rect {
	x := uint(r.frame.X)
	y := uint(r.frame.Y)
	width := uint(r.frame.W)
	height := uint(r.frame.H)
	FG1 := t.BoxLight
	FG2 := t.BoxDark
	if !extrude {
		FG1 = t.BoxDark
		FG2 = t.BoxLight
	}
	c.WriteRune(x, y, FG1, t.BoxBackground, t.TL)
	//c.Write(x+1, y, FG1, t.BoxBackground, RepeatRune(t.HT, width-2))
	for i := x + 1; i < x+(width-1); i++ {
		c.WriteRune(i, y, FG1, t.BoxBackground, t.HT)
	}
	c.WriteRune(x+width-1, y, FG1, t.BoxBackground, t.TR)
	for i := y + 1; i < y+height; i++ {
		c.WriteRune(x, i, FG1, t.BoxBackground, t.VL)
		c.Write(x+1, i, FG1, t.BoxBackground, RepeatRune(' ', width-2))
		c.WriteRune(x+width-1, i, FG2, t.BoxBackground, t.VR)
	}
	c.WriteRune(x, y+height-1, FG1, t.BoxBackground, t.BL)
	for i := x + 1; i < x+(width-1); i++ {
		c.WriteRune(i, y+height-1, FG2, t.BoxBackground, t.HB)
	}
	//c.Write(x+1, y+height-1, FG2, t.BoxBackground, RepeatRune(t.HB, width-2))
	c.WriteRune(x+width-1, y+height-1, FG2, t.BoxBackground, t.BR)
	return &Rect{int(x), int(y), int(width), int(height)}
}

// Draw a list widget. Takes a Box struct for the size and position.
// Takes a list of strings to be listed and an int that represents
// which item is currently selected. Does not scroll or wrap.
func (t *Theme) DrawList(c *vt100.Canvas, r *Box, items []string, selected int) {
	for i, s := range items {
		color := t.ListText
		if i == selected {
			color = t.ListFocus
		}
		c.Write(uint(r.frame.X), uint(r.frame.Y+i), color, t.ListBackground, s)
	}
}

// Draws a button widget at the given placement,
// with the given text. If active is False,
// it will look more "grayed out".
func (t *Theme) DrawButton(c *vt100.Canvas, r *Box, text string, active bool) {
	color := t.ButtonText
	if active {
		color = t.ButtonFocus
	}
	x := r.frame.X
	y := r.frame.Y
	c.Write(uint(x), uint(y), color, t.Background, "<  ")
	c.Write(uint(x+3), uint(y), color, t.Background, text)
	c.Write(uint(x+3+len(text)), uint(y), color, t.Background, "  >")
}

// Outputs a multiline string at the given coordinates.
// Uses the box background color.
// Returns the final y coordinate after drawing.
func (t *Theme) DrawAsciiArt(c *vt100.Canvas, x, y int, text string) int {
	var i int
	for i, line := range SplitTrim(text) {
		c.Write(uint(x), uint(y+i), t.Text, t.BoxBackground, line)
	}
	return y + i
}

// Outputs a multiline string at the given coordinates.
// Uses the default background color.
// Returns the final y coordinate after drawing.
func (t *Theme) DrawRaw(c *vt100.Canvas, x, y int, text string) int {
	var i int
	for i, line := range SplitTrim(text) {
		c.Write(uint(x), uint(y+i), t.Text, t.Background, line)
	}
	return y + i
}
