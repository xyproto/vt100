package main

import (
	"math/rand"

	"github.com/xyproto/vt100"
)

const evilGobblerEraseChar = ' ' // for erasing when moving

type EvilGobbler struct {
	hunting         *Gobbler             // current gobbler to hunt
	color           vt100.AttributeColor // foreground color
	x, y            int                  // current position
	oldx, oldy      int                  // previous position
	huntingDistance float64              // how far to closest gobbler
	counter         uint
	state           rune // looks
}

func NewEvilGobbler() *EvilGobbler {
	return &EvilGobbler{
		x:       10,
		y:       10,
		oldx:    10,
		oldy:    10,
		state:   'e',
		color:   vt100.LightYellow,
		counter: 0,
	}
}

func (b *EvilGobbler) Draw(c *vt100.Canvas) {
	c.PlotColor(uint(b.x), uint(b.y), b.color, b.state)
}

func (g *EvilGobbler) Next(c *vt100.Canvas, gobblers []*Gobbler, bob *Bob) bool {
	g.oldx = g.x
	g.oldy = g.y

	var hunting *Gobbler = nil
	var huntingDistance float64 = 99999.9

	for _, b := range gobblers {
		if d := distance(b.x, g.x, b.y, g.y); !b.dead && d <= huntingDistance {
			hunting = b
			huntingDistance = d
		}
	}

	if hunting == nil {

		g.x += rand.Intn(3) - 1
		g.y += rand.Intn(3) - 1

	} else {

		xspeed := 1
		yspeed := 1

		if g.x < hunting.x {
			g.x += xspeed
		} else if g.x > hunting.x {
			g.x -= xspeed
		}
		if g.y < hunting.y {
			g.y += yspeed
		} else if g.y > hunting.y {
			g.y -= yspeed
		}

		if distance(bob.x, g.x, bob.y, g.y) < 5 {
			g.x = g.oldx + (rand.Intn(3) - 1)
			g.y = g.oldy + (rand.Intn(3) - 1)
		}

		if !hunting.dead && huntingDistance < 2 || (hunting.x == g.x && hunting.y == g.y) {
			hunting.dead = true
			g.counter++
			hunting = nil
		}
	}

	if g.x > int(c.W()) {
		g.x = g.oldx
	} else if g.x < 0 {
		g.x = g.oldx
	}

	if g.y > int(c.H()) {
		g.y = g.oldy
	} else if g.y < 0 {
		g.y = g.oldy
	}

	return (g.x != g.oldx || g.y != g.oldy)
}

// Terminal was resized
func (b *EvilGobbler) Resize() {
	b.color = vt100.White
}
