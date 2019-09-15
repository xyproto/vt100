package main

import (
	"github.com/xyproto/vt100"
)

const bobEraseChar = ' ' // for erasing when moving

type Bob struct {
	x, y       int    // current position
	oldx, oldy int    // previous position
	state      rune   // looks
	color      string // foreground color
}

func NewBob() *Bob {
	return &Bob{
		x:     10,
		y:     10,
		oldx:  10,
		oldy:  10,
		state: 'o',
		color: "Yellow",
	}
}

func (b *Bob) ToggleColor() {
	const c1 = "Red"
	const c2 = "Yellow"
	if b.color == c1 {
		b.color = c2
	} else {
		b.color = c1
	}
}

func (b *Bob) ToggleState() {
	const up = 'O'
	const down = 'o'
	if b.state == up {
		b.state = down
	} else {
		b.state = up
	}
}

func (b *Bob) Draw(c *vt100.Canvas) {
	c.PlotC(uint(b.x), uint(b.y), b.color, b.state)
}

func (b *Bob) Right(c *vt100.Canvas) bool {
	oldx := b.x
	b.x += 1
	if b.x >= int(c.W()) {
		b.x -= 1
		return false
	}
	b.oldx = oldx
	b.oldy = b.y
	return true
}

func (b *Bob) Left(c *vt100.Canvas) bool {
	oldx := b.x
	if b.x-1 < 0 {
		return false
	}
	b.x -= 1
	b.oldx = oldx
	b.oldy = b.y
	return true
}

func (b *Bob) Up(c *vt100.Canvas) bool {
	oldy := b.y
	if b.y-1 < 0 {
		return false
	}
	b.y -= 1
	b.oldx = b.x
	b.oldy = oldy
	return true
}

func (b *Bob) Down(c *vt100.Canvas) bool {
	oldy := b.y
	b.y += 1
	if b.y >= int(c.H()) {
		b.y -= 1
		return false
	}
	b.oldx = b.x
	b.oldy = oldy
	return true
}

// Terminal was resized
func (b *Bob) Resize() {
	b.color = "Magenta"
}
