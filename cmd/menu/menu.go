package main

import (
	"github.com/xyproto/vt100"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Menu starts a loop where keypresses are handled. When a choice is made, a number is returned.
// -1 is "no choice", 0 and up is which choice were selected.
func Menu(title, titleColor string, choices []string, selectionDelay time.Duration, fg, hi, active string) int {
	c := vt100.NewCanvas()
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}

	// Mutex used when the terminal is resized
	resizeMut := &sync.RWMutex{}

	var (
		menu    = NewMenuWidget(title, titleColor, choices, fg, hi, active, c.W(), c.H())
		sigChan = make(chan os.Signal, 1)
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
			menu.Resize()
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
		menu.Draw(c)
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

		// Handle events
		key := tty.Key()
		switch key {
		case 38: // Up
			resizeMut.Lock()
			menu.Up(c)
			resizeMut.Unlock()
		case 40: // Down
			resizeMut.Lock()
			menu.Down(c)
			resizeMut.Unlock()
		case 27, 113: // ESC or q
			running = false
			break
		case 32, 13, 39: // Space, Return or Right
			resizeMut.Lock()
			menu.Select()
			resizeMut.Unlock()
			running = false
			break
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57: // 0 .. 9
			number := uint(key - 48)
			resizeMut.Lock()
			menu.SelectIndex(number)
			resizeMut.Unlock()
		}
	}

	if menu.Selected() >= 0 {
		// Draw the selected item in a different color for a short while
		resizeMut.Lock()
		menu.SelectDraw(c)
		resizeMut.Unlock()
		c.Draw()
		time.Sleep(selectionDelay)
	}

	tty.Close()

	vt100.SetLineWrap(true)
	vt100.ShowCursor(true)
	vt100.Home()

	return menu.Selected()
}
