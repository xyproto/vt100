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
