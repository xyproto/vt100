package main

import (
	"github.com/xyproto/vt100"
	"io/ioutil"
	"log"
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
	defer tty.Close()

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
				vt100.Clear()
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

	vt100.Init()
	defer vt100.Close()

	// The loop time that is aimed for
	loopDuration := time.Millisecond * 10
	start := time.Now()

	running := true

	// Don't output keypress terminal codes on the screen
	tty.NoBlock()

	var key int

	for running {

		// Draw elements in their new positions
		c.Clear()
		//c.Draw()

		resizeMut.RLock()
		for _, bullet := range bullets {
			bullet.Draw(c)
		}
		bob.Draw(c)
		resizeMut.RUnlock()

		//vt100.Clear()

		// Update the canvas
		c.Draw()

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
		key = tty.Key()
		switch key {
		case 253: // Up
			resizeMut.Lock()
			moved = bob.Up(c)
			resizeMut.Unlock()
		case 255: // Down
			resizeMut.Lock()
			moved = bob.Down(c)
			resizeMut.Unlock()
		case 254: // Right
			resizeMut.Lock()
			moved = bob.Right(c)
			resizeMut.Unlock()
		case 252: // Left
			resizeMut.Lock()
			moved = bob.Left(c)
			resizeMut.Unlock()
		case 27, 113: // ESC or q
			running = false
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
			if r == rune(0) || r == bobEraseChar || r == bulletEraseChar {
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
				log.Fatalln(err)
			}
		}

		// If a key was pressed, clear the screen, just in case it shifted
		//if key != 0 {
		//	vt100.Clear()
		//}

		// Change state
		resizeMut.Lock()
		for _, bullet := range bullets {
			bullet.Next(c)
		}
		if moved {
			bob.ToggleState()
		}
		resizeMut.Unlock()

		// Erase all previous positions not occupied by current items
		c.Plot(uint(bob.oldx), uint(bob.oldy), bobEraseChar)
		for _, bullet := range bullets {
			c.Plot(uint(bullet.oldx), uint(bullet.oldy), bulletEraseChar)
		}
	}
}
