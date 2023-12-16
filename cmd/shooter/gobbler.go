package main

import (
	"math/rand"

	"github.com/xyproto/vt100"
)

const gobblerEraseChar = ' ' // for erasing when moving

type Gobbler struct {
	hunting         *Bullet              // current bullet to hunt
	color           vt100.AttributeColor // foreground color
	x, y            int                  // current position
	oldx, oldy      int                  // previous position
	huntingDistance float64              // how far to closest bullet
	counter         uint
	state           rune // looks
	dead            bool
}

func NewGobbler() *Gobbler {
	return &Gobbler{
		x:               10,
		y:               10,
		oldx:            10,
		oldy:            10,
		state:           'g',
		color:           vt100.LightGreen,
		hunting:         nil,
		huntingDistance: 0,
		counter:         0,
		dead:            false,
	}
}

func NewGobblers(n int) []*Gobbler {
	gobblers := make([]*Gobbler, n)
	for i := range gobblers {
		gobblers[i] = NewGobbler()
	}
	return gobblers
}

func (g *Gobbler) Draw(c *vt100.Canvas) {
	c.PlotColor(uint(g.x), uint(g.y), g.color, g.state)
}

func (g *Gobbler) Next(c *vt100.Canvas, bullets []*Bullet, bob *Bob) bool {
	if g.dead {
		g.state = 'T'
		g.color = vt100.LightCyan
		return true
	}

	g.oldx = g.x
	g.oldy = g.y

	// Move to the nearest bullet and eat it
	if len(bullets) == 0 {

		g.x += rand.Intn(5) - 2
		g.y += rand.Intn(5) - 2

	} else {

		if g.hunting == nil || g.hunting.removed == true {
			var minDistance float64 = 99999.9
			var closestBullet *Bullet = nil
			for _, b := range bullets {
				if d := distance(b.x, g.x, b.y, g.y); !b.removed && d <= minDistance {
					closestBullet = b
					minDistance = d
				}
			}
			if closestBullet != nil {
				g.hunting = closestBullet
				g.huntingDistance = minDistance
			}
		} else {
			g.huntingDistance = distance(g.hunting.x, g.x, g.hunting.y, g.y)
		}

		if g.hunting == nil {

			g.x += rand.Intn(5) - 2
			g.y += rand.Intn(5) - 2

		} else {

			xspeed := 1
			yspeed := 1

			if abs(g.hunting.x-g.x) >= abs(g.hunting.y-g.y) {
				// Longer away along x than along y
				if g.huntingDistance > 20 {
					xspeed = 3
					yspeed = 2
				} else if g.huntingDistance > 10 {
					xspeed = 2 + rand.Intn(2)
					yspeed = 2
				}
			} else {
				// Longer away along x than along y
				if g.huntingDistance > 20 {
					xspeed = 2
					yspeed = 3
				} else if g.huntingDistance > 10 {
					xspeed = 2
					yspeed = 2 + rand.Intn(2)
				}
			}

			if g.x < g.hunting.x {
				g.x += xspeed
			} else if g.x > g.hunting.x {
				g.x -= xspeed
			}
			if g.y < g.hunting.y {
				g.y += yspeed
			} else if g.y > g.hunting.y {
				g.y -= yspeed
			}

			if distance(bob.x, g.x, bob.y, g.y) < 15 {
				g.x = g.oldx + (rand.Intn(3) - 1)
				g.y = g.oldy + (rand.Intn(3) - 1)
			}

			if !g.hunting.removed && g.huntingDistance < 2 || (g.hunting.x == g.x && g.hunting.y == g.y) {
				g.hunting.removed = true
				g.counter++
				g.hunting = nil
			}
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
func (b *Gobbler) Resize() {
	b.color = vt100.White
}
