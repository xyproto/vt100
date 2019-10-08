package main

import (
	"errors"
	"github.com/xyproto/vt100"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unicode"
)

// Returns the Nth letter in a given string, as lowercase. Ignores numbers, special characters, whitespace etc.
func getLetter(s string, pos int) (rune, error) {
	counter := 0
	for _, letter := range s {
		if unicode.IsLetter(letter) {
			if counter == pos {
				return unicode.ToLower(letter), nil
			}
			counter++
		}
	}
	return rune(0), errors.New("no letter")
}

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
				vt100.Clear()
				c = nc
				c.Redraw()
			}

			// Inform all elements that the terminal was resized
			menu.Resize()
			resizeMut.Unlock()
		}
	}()

	vt100.Init()
	defer vt100.Close()

	// The loop time that is aimed for
	loopDuration := time.Millisecond * 20
	start := time.Now()

	running := true

	vt100.Clear()
	c.Redraw()

	for running {

		// Draw elements in their new positions
		//vt100.Clear()

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
		case 253, 107, 16: // Up, k or ctrl-p
			resizeMut.Lock()
			menu.Up(c)
			resizeMut.Unlock()
		case 255, 106, 14: // Down, j or ctrl-n
			resizeMut.Lock()
			menu.Down(c)
			resizeMut.Unlock()
		case 1: // Top, ctrl-a
			resizeMut.Lock()
			menu.SelectFirst()
			resizeMut.Unlock()
		case 5: // Bottom, ctrl-e
			resizeMut.Lock()
			menu.SelectLast()
			resizeMut.Unlock()
		case 27, 113: // ESC or q
			running = false
			break
		case 32, 13, 254: // Space, Return or Right // 108 is l
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
		default:
			letterNumber := 0
			// Check if the key matches the first letter (a-z,A-Z) in the choices
			if 65 <= key && key <= 90 {
				letterNumber = key - 65
			} else if 97 <= key && key <= 122 {
				letterNumber = key - 97
			} else {
				break
			}
			var r rune = rune(letterNumber + 97)

			// Select the item that starts with this letter, if possible. Try the first, then the second, etc, up to 5
			keymap := make(map[rune]int)
			for index, choice := range choices {
				for pos := 0; pos < 5; pos++ {
					letter, err := getLetter(choice, pos)
					if err == nil {
						_, exists := keymap[letter]
						// If the letter is not already stored in the keymap, and it's not q, j or k
						if !exists && (letter != 113) && (letter != 106) && (letter != 107) {
							keymap[letter] = index
							// Found a letter for this choice, move on
							break
						}
					}
				}
				// Did not find a letter for this choice, move on
			}

			// Choose the index for the letter that was pressed and found in the keymap, if found
			for letter, index := range keymap {
				if letter == r {
					resizeMut.Lock()
					menu.SelectIndex(uint(index))
					resizeMut.Unlock()
				}
			}
		}

		// If a key was pressed, draw the canvas
		if key != 0 {
			c.Redraw()
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

	return menu.Selected()
}
