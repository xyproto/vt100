package main

import (
	"github.com/xyproto/vt100"
)

const bobEraseChar = ' ' // for erasing when moving

type Bob struct {
	color      vt100.AttributeColor // foreground color
	x, y       int                  // current position
	oldx, oldy int                  // previous position
	state      rune                 // looks
}

func NewBob() *Bob {
	return &Bob{
		x:     10,
		y:     10,
		oldx:  10,
		oldy:  10,
		state: 'o',
		color: vt100.Red,
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
	c.PlotColor(uint(b.x), uint(b.y), b.color, b.state)
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
	b.color = vt100.LightMagenta
}
