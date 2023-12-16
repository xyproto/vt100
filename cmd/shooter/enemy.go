package main

import (
	"math"
	"math/rand"

	"github.com/xyproto/vt100"
)

const enemyEraseChar = ' ' // for erasing when moving

type Enemy struct {
	color      vt100.AttributeColor // foreground color
	x, y       int                  // current position
	oldx, oldy int                  // previous position
	state      rune                 // looks
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
		state: '°',
		color: vt100.LightMagenta,
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

	d := distance(bob.x, e.x, bob.y, e.y)
	if d > 10 {
		if e.x < bob.x {
			e.x++
		} else if e.x > bob.x {
			e.x--
		}
		if e.y < bob.y {
			e.y++
		} else if e.y > bob.y {
			e.y--
		}
	} else {
		for {
			dx := e.x - e.oldx
			dy := e.y - e.oldy
			e.x += int(math.Round(float64(dx*3+rand.Intn(5)-2) / float64(4))) // -2, -1, 0, 1, 2
			e.y += int(math.Round(float64(dy*3+rand.Intn(5)-2) / float64(4)))
			if e.x != e.oldx {
				break
			}
			if e.y != e.oldy {
				break
			}
		}
	}

	if e.HitSomething(c) {
		e.x = e.oldx
		e.y = e.oldy
		return false
	}

	if e.x >= int(c.W()) {
		e.x = e.oldx
	} else if e.x <= 0 {
		e.x = e.oldx
	}
	if e.y >= int(c.H()) {
		e.y = e.oldy
	} else if e.y <= 0 {
		e.y = e.oldy
	}

	return e.x != e.oldx || e.y != e.oldy
}

func (e *Enemy) HitSomething(c *vt100.Canvas) bool {
	r, err := c.At(uint(e.x), uint(e.y))
	if err != nil {
		return false
	}
	// Hit something?
	return r != rune(0) && r != bulletEraseChar && r != bobEraseChar && r != enemyEraseChar
}
