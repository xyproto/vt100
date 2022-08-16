package main

import (
	"math"
	"math/rand"

	"github.com/xyproto/vt100"
)

const enemyEraseChar = ' ' // for erasing when moving

type Enemy struct {
	x, y       int                  // current position
	oldx, oldy int                  // previous position
	state      rune                 // looks
	color      vt100.AttributeColor // foreground color
}

func NewEnemies(n int) []*Enemy {
	enemies := make([]*Enemy, n)
	for i := range enemies {
		enemies[i] = NewEnemy()
	}
	return enemies
}

func NewEnemy() *Enemy {
	return &Enemy{
		x:     10,
		y:     10,
		oldx:  10,
		oldy:  10,
		state: 'x',
		color: vt100.LightCyan,
	}
}

func (b *Enemy) Draw(c *vt100.Canvas) {
	c.PlotColor(uint(b.x), uint(b.y), b.color, b.state)
}

func (b *Enemy) Right(c *vt100.Canvas) bool {
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

func (b *Enemy) Left(c *vt100.Canvas) bool {
	oldx := b.x
	if b.x-1 < 0 {
		return false
	}
	b.x -= 1
	b.oldx = oldx
	b.oldy = b.y
	return true
}

func (b *Enemy) Up(c *vt100.Canvas) bool {
	oldy := b.y
	if b.y-1 < 0 {
		return false
	}
	b.y -= 1
	b.oldx = b.x
	b.oldy = oldy
	return true
}

func (b *Enemy) Down(c *vt100.Canvas) bool {
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
func (b *Enemy) Resize() {
	b.color = vt100.LightCyan
}

// Next moves the object to the next position, and returns true if it moved
func (e *Enemy) Next(c *vt100.Canvas, bob *Bob) bool {
	e.oldx = e.x
	e.oldy = e.y

	// Now try to move the enemy intelligently, given the position of bob

	distance := math.Sqrt(float64(bob.x)*float64(bob.x) + float64(bob.y)*float64(bob.y) - float64(e.x)*float64(e.x) + float64(e.y)*float64(e.y))

	if distance > 10 {
		e.x += rand.Intn(3) - 1 // -1 0 or 1
		e.y += rand.Intn(3) - 1 // -1 0 or 1
	} else {
		e.x += rand.Intn(5) - 2 // -2 -1 0 1 or 2
		e.y += rand.Intn(5) - 2 // -2 -1 0 1 or 2
	}

	if e.HitSomething(c) {
		// enemy did hit something, move back one step
		e.x = e.oldx
		e.y = e.oldy
		return false
	}
	return true
}

func (e *Enemy) HitSomething(c *vt100.Canvas) bool {
	r, err := c.At(uint(e.x), uint(e.y))
	if err != nil {
		return false
	}
	if r != rune(0) && r != bulletEraseChar && r != bobEraseChar {
		// Hit something
		return true
	}
	return false
}
