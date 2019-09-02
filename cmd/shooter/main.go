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
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}

	// Mutex used when the terminal is resized
	resizeMut := &sync.RWMutex{}

	var (
		bob     = NewBob()
		sigChan = make(chan os.Signal, 1)
		bullets = make([]*Bullet, 0)
	)

	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			resizeMut.Lock()
			// Create a new canvas, with the new size
			nc := c.Resized()
			if nc != nil {
				c.Clear()
				c.Draw()
				c = nc
			}

			// Inform all elements that the terminal was resized
			// TODO: Use a slice of interfaces that can contain all elements
			for _, bullet := range bullets {
				bullet.Resize()
			}
			bob.Resize()
			resizeMut.Unlock()
		}
	}()

	vt100.Clear()
	vt100.ShowCursor(false)
	vt100.SetLineWrap(false)

	// The loop time that is aimed for
	loopDuration := time.Millisecond * 20
	start := time.Now()

	running := true

	for running {

		// Draw elements in their new positions
		vt100.Clear()

		resizeMut.RLock()
		for _, bullet := range bullets {
			bullet.Draw(c)
		}
		bob.Draw(c)
		resizeMut.RUnlock()

		// Update the canvas
		c.Draw()

		// Don't output keypress terminal codes on the screen
		tty.NoBlock()

		// Wait a bit
		end := time.Now()
		passed := end.Sub(start)
		if passed < loopDuration {
			remaining := loopDuration - passed
			time.Sleep(remaining)
		}
		start = time.Now()

		// Has the player moved?
		moved := false

		// Handle events
		switch tty.Key() {
		case 38: // Up
			resizeMut.Lock()
			moved = bob.Up(c)
			resizeMut.Unlock()
		case 40: // Down
			resizeMut.Lock()
			moved = bob.Down(c)
			resizeMut.Unlock()
		case 39: // Right
			resizeMut.Lock()
			moved = bob.Right(c)
			resizeMut.Unlock()
		case 37: // Left
			resizeMut.Lock()
			moved = bob.Left(c)
			resizeMut.Unlock()
		case 27, 113: // ESC or q
			running = false
			break
		case 32: // Space
			resizeMut.Lock()
			bob.ToggleColor()
			resizeMut.Unlock()
			// Check if the place to the right is available
			r, err := c.At(uint(bob.x+1), uint(bob.y))
			if err != nil {
				// No free place to the right
				break
			}
			if r == rune(0) {
				// Fire a new bullet
				bullets = append(bullets, NewBullet(bob.x+1, bob.y, 1, 0))
			}
		case 97: // a
			// Write the canvas characters to file
			resizeMut.RLock()
			b := []byte(c.String())
			resizeMut.RUnlock()
			err := ioutil.WriteFile("canvas.txt", b, 0644)
			if err != nil {
				panic(err)
				running = false
				break
			}
		}

		// Change state
		resizeMut.Lock()
		for _, bullet := range bullets {
			bullet.Next(c)
		}
		if moved {
			bob.ToggleState()
		}
		resizeMut.Unlock()

		// Erase all previous positions
		c.Plot(uint(bob.oldx), uint(bob.oldy), rune(0))
		for _, bullet := range bullets {
			c.Plot(uint(bullet.oldx), uint(bullet.oldy), rune(0))
		}
	}

	tty.Close()

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
	vt100.Home()
}
