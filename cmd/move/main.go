package main

import (
	"github.com/xyproto/vt100"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	var draw sync.RWMutex

	x := uint(10)
	y := uint(10)

	c := vt100.NewCanvas()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			// terminal was resized
			draw.Lock()
			c.Resize()
			draw.Unlock()
		}
	}()

	c.Clear()
	c.ShowCursor(false)
	//c.SetLineWrap(false)

	running := true
	for running {

		draw.Lock()

		// Draw elements in their new positions
		c.PlotC(x, y, "Yellow", "o")

		// Update the canvas
		c.Draw()

		draw.Unlock()

		// Wait a bit
		time.Sleep(time.Millisecond * 20)

		// Change state
		oldx := x
		oldy := y

		// Handle events
		switch vt100.Key() {
		case 38:
			y -= 1
		case 40:
			y += 1
		case 39:
			x += 1
		case 37:
			x -= 1
		case 27, 113: // ESC or q
			running = false
			break
		}

		draw.Lock()

		// Erase elements at their old positions
		c.Plot(oldx, oldy, " ")

		draw.Unlock()
	}

	//c.SetLineWrap(true)
	c.ShowCursor(true)
	c.Home()
}
