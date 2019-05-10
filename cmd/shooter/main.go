package main

import (
	"github.com/xyproto/vt100"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	c := vt100.NewCanvas()
	r := vt100.NewRawTerminal()

	var (
		draw    sync.RWMutex
		bob     = NewBob()
		sigChan = make(chan os.Signal, 1)
		bullets = make([]*Bullet, 0)
	)

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
	vt100.SetLineWrap(false)

	running := true
	for running {

		// Draw elements in their new positions
		draw.Lock()
		bob.Draw(c)
		for _, bullet := range bullets {
			bullet.Draw(c)
		}
		draw.Unlock()

		// Update the canvas
		draw.Lock()
		c.Draw()
		draw.Unlock()

		// Wait a bit
		time.Sleep(time.Millisecond * 15)

		// Change state
		for _, bullet := range bullets {
			bullet.Next(c)
		}
		// Has the player moved?
		moved := false

		// Handle events
		draw.Lock()
		switch r.Key() {
		case 38: // Up
			moved = bob.Up(c)
		case 40: // Down
			moved = bob.Down(c)
		case 39: // Right
			moved = bob.Right(c)
		case 37: // Left
			moved = bob.Left(c)
		case 27, 113: // ESC or q
			running = false
			break
		case 32: // Space
			bob.ToggleColor()
			bullets = append(bullets, NewBullet(bob.x, bob.y, 1, 0))
		case 97: // a
			// Write the canvas characters to file
			b := []byte(c.String())
			err := ioutil.WriteFile("canvas.txt", b, 0644)
			if err != nil {
				panic(err)
				running = false
				break
			}
		}
		if moved {
			bob.ToggleState()
		}
		draw.Unlock()

		// Erase all previous positions
		draw.Lock()
		if moved {
			c.Plot(uint(bob.oldx), uint(bob.oldy), rune(0))
		}
		for _, bullet := range bullets {
			c.Plot(uint(bullet.oldx), uint(bullet.oldy), rune(0))
		}
		draw.Unlock()
	}

	r.Close()

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
	vt100.Home()
}
