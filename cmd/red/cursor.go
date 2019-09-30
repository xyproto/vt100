package main

import (
	"github.com/xyproto/vt100"
)

type Cursor struct {
	X, Y int
}

func (cur *Cursor) ScreenWrap(c *vt100.Canvas) {
	w, h := c.Size()
	if cur.X < 0 {
		cur.X = int(w)
		cur.Y--
	}
	if cur.Y < 0 {
		cur.Y = 0
	}
	if cur.X >= int(w) {
		cur.X = 0
		cur.Y++
	}
	if cur.Y >= int(h) {
		cur.Y = int(h)
	}
}

func scrollDown(c *vt100.Canvas, offset int, status *StatusBar, e *Editor, dataCursor *Cursor, scrollSpeed int) (bool, int) {
	redraw := false
	h := int(c.H())
	if offset == e.Len()-h {
		// Status message
		status.SetMessage("End of text")
		status.Show(c, offset)
		c.Draw()
	} else {
		status.Clear(c)
		// Find out if we can scroll scrollSpeed, or less
		canScroll := scrollSpeed
		if (offset + canScroll) >= (e.Len() - h) {
			// Almost at the bottom, we can scroll the remaining lines
			canScroll = (e.Len() - h) - offset
		}
		// Only move the data cursor down one or more lines, do not move the screen cursor
		dataCursor.Y += canScroll
		dataCursor.ScreenWrap(c)
		// Move the scroll offset
		offset += canScroll
		// Prepare to redraw
		//vt100.Clear()
		redraw = true
	}
	return redraw, offset
}

func scrollUp(c *vt100.Canvas, offset int, status *StatusBar, e *Editor, dataCursor *Cursor, scrollSpeed int) (bool, int) {
	redraw := false
	if offset == 0 {
		// Can't scroll further up
		// Status message
		status.SetMessage("Start of text")
		status.Show(c, offset)
		c.Draw()
	} else {
		status.Clear(c)
		// Find out if we can scroll scrollSpeed, or less
		canScroll := scrollSpeed
		if offset-canScroll < 0 {
			// Almost at the top, we can scroll the remaining lines
			canScroll = offset
		}
		// Only move the data cursor up one or more lines, do not move the screen cursor
		dataCursor.Y -= canScroll
		dataCursor.ScreenWrap(c)
		// Move the scroll offset
		offset -= canScroll
		// Prepare to redraw
		redraw = true
	}
	return redraw, offset
}

type Position struct {
	sx     int // the position of the cursor in the current scrollview
	sy     int // the position of the cursor in the current scrollview
	scroll int // how far one has scrolled
}

func position2datacursor(p *Position, e *Editor) *Cursor {
	var dataX int
	// the y position in the data is the lines scrolled + current screen cursor Y position
	dataY := p.scroll + p.sy
	// get the current line of text
	line := e.Line(dataY)
	screenCounter := 0 // counter the characters on the screen
	// loop, while also keeping track of tab expansion
	for i, r := range line {
		if r == '\t' {
			screenCounter += e.spacesPerTab
		} else {
			screenCounter += 1
		}
		// When we reached the correct screen position, use i as the data position
		if screenCounter == p.sx {
			dataX = i
			break
		}
	}
	// Return the data cursor
	return &Cursor{dataX, dataY}
}

func (p *Position) ScreenX() int {
	return p.sx
}

func (p *Position) ScreenY() int {
	return p.sy
}

func (p *Position) ScrollOffset() int {
	return p.scroll
}

func (p *Position) DataCursor(e *Editor) *Cursor {
	return position2datacursor(p, e)
}

func (p *Position) ScreenCursor() *Cursor {
	return &Cursor{p.sx, p.sy}
}

func (p *Position) SetScreenScursor(c *Cursor) {
	p.sx = c.X
	p.sy = c.Y
}

func (p *Position) SetScreenX(x int) {
	p.sx = x
}

func (p *Position) SetScreenY(y int) {
	p.sy = y
}

func (p *Position) SetOffset(offset int) {
	p.scroll = offset
}
