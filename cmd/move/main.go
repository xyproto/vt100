package main

import (
	"github.com/xyproto/vt100"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Bob struct {
	x, y  int
	color string
	state rune
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
	c.PlotColor(uint(b.x), uint(b.y), vt100.LightColorMap[b.color], b.state)
}

func (b *Bob) Right(c *vt100.Canvas) bool {
	b.x += 1
	if b.x >= int(c.W()) {
		b.x -= 1
		return false
	}
	return true
}

func (b *Bob) Left(c *vt100.Canvas) bool {
	if b.x-1 < 0 {
		return false
	}
	b.x -= 1
	return true
}

func (b *Bob) Up(c *vt100.Canvas) bool {
	if b.y-1 < 0 {
		return false
	}
	b.y -= 1
	return true
}

func (b *Bob) Down(c *vt100.Canvas) bool {
	b.y += 1
	if b.y >= int(c.H()) {
		b.y -= 1
		return false
	}
	return true
}

func main() {

	c := vt100.NewCanvas()
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}

	var bob Bob
	bob.state = 'o'
	bob.color = "Yellow"
	bob.x = 10
	bob.y = 10

	var draw sync.RWMutex
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			// Terminal was resized

			// Prepare to resize the canvas
			draw.Lock()

			// Clear the screen after the resize
			vt100.Clear()

			// Create a new canvas, with the new size
			nc := c.Resized()
			if nc != nil {
				c.Clear()
				c.Draw()
				c = nc
			}
			// Clear again, for good measure
			vt100.Clear()
			// Redraw the characters
			c.Redraw()
			// Done
			draw.Unlock()
		}
	}()

	vt100.Clear()
	vt100.ShowCursor(false)
	vt100.SetLineWrap(true)

	running := true
	for running {

		//vt100.Clear()

		// Draw elements in their new positions
		draw.Lock()
		bob.Draw(c)
		draw.Unlock()

		// Update the canvas
		draw.Lock()
		c.Draw()
		draw.Unlock()

		// Wait a bit
		time.Sleep(time.Millisecond * 10)

		// Change state
		oldx := bob.x
		oldy := bob.y
		moved := false

		// Handle events
		draw.Lock()
		switch tty.Key() {
		case 253: // Up
			moved = bob.Up(c)
		case 255: // Down
			moved = bob.Down(c)
		case 254: // Right
			moved = bob.Right(c)
		case 252: // Left
			moved = bob.Left(c)
		case 27, 113: // ESC or q
			running = false
			break
		case 32: // Space
			bob.ToggleColor()
		}
		draw.Unlock()

		if moved {
			bob.ToggleState()

			// Erase elements at their old positions
			draw.Lock()
			c.Plot(uint(oldx), uint(oldy), ' ')
			draw.Unlock()
		}
	}
	tty.Close()

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
	vt100.Home()
}
