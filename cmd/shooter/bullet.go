package main

import (
	"github.com/xyproto/vt100"
)

type Bullet struct {
	x, y       int    // current position
	oldx, oldy int    // previous position
	vx, vy     int    // velocity
	state      rune   // looks
	color      string // foreground color
}

func NewBullet(x, y, vx, vy int) *Bullet {
	return &Bullet{
		x:     x,
		y:     y,
		oldx:  x,
		oldy:  y,
		vx:    vx,
		vy:    vy,
		state: '×',
		color: "Blue",
	}
}

func (b *Bullet) ToggleColor() {
	const c1 = "Red"
	const c2 = "Blue"
	if b.color == c1 {
		b.color = c2
	} else {
		b.color = c1
	}
}

func (b *Bullet) ToggleState() {
	const up = '×'
	const down = '-'
	if b.state == up {
		b.state = down
	} else {
		b.state = up
	}
}

func (b *Bullet) Draw(c *vt100.Canvas) {
	c.PlotC(uint(b.x), uint(b.y), b.color, b.state)
}

// Next moves the object to the next position, and returns true if it moved
func (b *Bullet) Next(c *vt100.Canvas) bool {
	if b.x-b.vx < 0 {
		return false
	}
	if b.y-b.vy < 0 {
		return false
	}
	b.oldx = b.x
	b.x += b.vx
	b.oldy = b.y
	b.y += b.vy
	if b.x >= int(c.W()) {
		b.x -= b.vx
		return false
	}
	if b.y >= int(c.H()) {
		b.y -= b.vy
		return false
	}
	return true
}

func (b *Bullet) HitSomething(c *vt100.Canvas) bool {
	r := c.At(uint(b.x), uint(b.y))
	if r != rune(0) {
		// Hit something
		b.ToggleColor()
		return true
	}
	return false
}
