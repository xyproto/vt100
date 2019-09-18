package main

import (
	"bytes"
	"github.com/xyproto/vt100"
)

const bulletEraseChar = ' ' // for erasing when moving

type Bullet struct {
	x, y       int                  // current position
	oldx, oldy int                  // previous position
	vx, vy     int                  // velocity
	state      rune                 // looks
	color      vt100.AttributeColor // foreground color
	stopped    bool                 // is the movement stopped?
}

func NewBullet(x, y, vx, vy int) *Bullet {
	return &Bullet{
		x:       x,
		y:       y,
		oldx:    x,
		oldy:    y,
		vx:      vx,
		vy:      vy,
		state:   '×',
		color:   vt100.LightBlue,
		stopped: false,
	}
}

func (b *Bullet) ToggleColor() {
	c1 := vt100.LightGreen
	c2 := vt100.LightBlue
	if bytes.Equal(b.color, c1) {
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
	c.PlotColor(uint(b.x), uint(b.y), b.color, b.state)
}

// Next moves the object to the next position, and returns true if it moved
func (b *Bullet) Next(c *vt100.Canvas) bool {
	if b.stopped {
		b.ToggleColor()
		return false
	}
	if b.x-b.vx < 0 {
		b.ToggleColor()
		return false
	}
	if b.y-b.vy < 0 {
		b.ToggleColor()
		return false
	}
	b.oldx = b.x
	b.oldy = b.y
	b.x += b.vx
	b.y += b.vy
	if b.HitSomething(c) {
		b.x = b.oldx
		b.y = b.oldy
		return false
	}
	if b.x >= int(c.W()) {
		b.x -= b.vx
		b.ToggleColor()
		return false
	}
	if b.y >= int(c.H()) {
		b.y -= b.vy
		b.ToggleColor()
		return false
	}
	return true
}

func (b *Bullet) Stop() {
	b.vx = 0
	b.vy = 0
	b.stopped = true
}

func (b *Bullet) HitSomething(c *vt100.Canvas) bool {
	r, err := c.At(uint(b.x), uint(b.y))
	if err != nil {
		return false
	}
	if r != rune(0) && r != bulletEraseChar && r != bobEraseChar {
		// Hit something. Check the next-next position too
		r2, err := c.At(uint(b.x+b.vx), uint(b.y+b.vy))
		if err != nil {
			return false
		}
		if r2 != rune(0) && r2 != bulletEraseChar && r2 != bobEraseChar {
			b.Stop()
		}
		return true
	}
	return false
}

// Terminal was resized
func (b *Bullet) Resize() {
	b.stopped = false
}
