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

	oColor := "Yellow"
	oRune := 'o'

	running := true
	for running {

		// Draw elements in their new positions
		draw.Lock()
		c.PlotC(x, y, oColor, oRune)
		draw.Unlock()

		// Update the canvas
		draw.Lock()
		c.Draw()
		draw.Unlock()

		// Wait a bit
		time.Sleep(time.Millisecond * 15)

		// Change state
		oldx := x
		oldy := y
		arrow := false

		// Handle events
		draw.Lock()
		switch vt100.Key() {
		case 38: // Up
			y -= 1
			arrow = true
		case 40: // Down
			y += 1
			arrow = true
		case 39: // Right
			x += 1
			arrow = true
		case 37: // Left
			x -= 1
			arrow = true
		case 27, 113: // ESC or q
			running = false
			break
		case 32: // Space
			if oColor == "Yellow" {
				oColor = "Red"
			} else {
				oColor = "Yellow"
			}
		}
		draw.Unlock()

		if arrow {
			if oRune == 'o' {
				oRune = 'O'
			} else {
				oRune = 'o'
			}
		}

		// Erase elements at their old positions
		draw.Lock()
		c.Plot(oldx, oldy, ' ')
		draw.Unlock()
	}

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
	vt100.Home()
}
