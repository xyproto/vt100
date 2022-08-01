package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/xyproto/vt100"
)

func main() {
	//rand.Seed(time.Now().UnixNano())

	c := vt100.NewCanvas()
	//c.FillBackground(vt100.Blue)

	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	defer tty.Close()

	vt100.EchoOff()

	// Mutex used when the terminal is resized
	resizeMut := &sync.RWMutex{}

	var (
		bob         = NewBob()
		sigChan     = make(chan os.Signal, 1)
		evilGobbler = NewEvilGobbler()
		gobblers    = NewGobblers(10)
		bullets     = make([]*Bullet, 0)
		enemies     = NewEnemies(7)
		score       = uint(0)
		highScore   = uint(0)
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
			for _, enemy := range enemies {
				enemy.Resize()
			}
			for _, gobbler := range gobblers {
				gobbler.Resize()
			}
			bob.Resize()
			evilGobbler.Resize()
			resizeMut.Unlock()
		}
	}()

	vt100.Init()
	defer vt100.Close()

	// The loop time that is aimed for
	loopDuration := time.Millisecond * 10
	start := time.Now()

	running := true
	paused := false
	var statusText string

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
		for _, enemy := range enemies {
			enemy.Draw(c)
		}
		evilGobbler.Draw(c)
		for _, gobbler := range gobblers {
			gobbler.Draw(c)
		}
		bob.Draw(c)
		c.Write(5, 1, vt100.LightRed, vt100.BackgroundDefault, statusText)
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
		case 253, 119: // Up or w
			resizeMut.Lock()
			moved = bob.Up(c)
			resizeMut.Unlock()
		case 255, 115: // Down or s
			resizeMut.Lock()
			moved = bob.Down(c)
			resizeMut.Unlock()
		case 254, 100: // Right or d
			resizeMut.Lock()
			moved = bob.Right(c)
			resizeMut.Unlock()
		case 252, 97: // Left or a
			resizeMut.Lock()
			moved = bob.Left(c)
			resizeMut.Unlock()
		case 27, 113: // ESC or q
			running = false
		case 32: // Space
			// Check if the place to the right is available
			r, err := c.At(uint(bob.x+1), uint(bob.y))
			if err != nil {
				// No free place to the right
				break
			}
			if r == rune(0) || r == bobEraseChar || r == bulletEraseChar || r == enemyEraseChar {
				// Fire a new bullet
				bullets = append(bullets, NewBullet(bob.x, bob.y, bob.x-bob.oldx, bob.y-bob.oldy))
			}
		case 112: // p
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

		if !paused {
			// Change state
			resizeMut.Lock()
			for _, bullet := range bullets {
				bullet.Next(c)
			}
			for _, enemy := range enemies {
				enemy.Next(c, bob)
			}
			for _, gobbler := range gobblers {
				gobbler.Next(c, bullets, bob)
			}
			evilGobbler.Next(c, gobblers, bob)
			if moved {
				bob.ToggleState()
			}
			resizeMut.Unlock()
		}

		// Erase all previous positions not occupied by current items
		c.Plot(uint(bob.oldx), uint(bob.oldy), bobEraseChar)
		c.Plot(uint(evilGobbler.oldx), uint(evilGobbler.oldy), evilGobblerEraseChar)
		for _, bullet := range bullets {
			c.Plot(uint(bullet.oldx), uint(bullet.oldy), bulletEraseChar)
		}
		for _, enemy := range enemies {
			c.Plot(uint(enemy.oldx), uint(enemy.oldy), enemyEraseChar)
		}
		for _, gobbler := range gobblers {
			c.Plot(uint(gobbler.oldx), uint(gobbler.oldy), gobblerEraseChar)
		}

		if !paused {

			// Clean up removed bullets
			filteredBullets := make([]*Bullet, 0, len(bullets))
			for _, bullet := range bullets {
				if !bullet.removed {
					filteredBullets = append(filteredBullets, bullet)
				} else {
					c.Plot(uint(bullet.x), uint(bullet.y), bulletEraseChar)
				}
			}
			bullets = filteredBullets

			gobblersAlive := false
			for _, gobbler := range gobblers {
				score += gobbler.counter
				(*gobbler).counter = 0
				if !gobbler.dead {
					gobblersAlive = true
				}
			}
			// evilGobbler.counter
			if gobblersAlive {
				statusText = fmt.Sprintf("Score: %d", score)
			} else {
				paused = true
				statusText = "Game over"

				// The player can still move around bob
				bob.state = '@'
				bob.color = vt100.White

				if score > highScore {
					statusText = fmt.Sprintf("Game over! New highscore: %d", score)
				} else if score > 0 {
					statusText = fmt.Sprintf("Game over! Score: %d", score)
				}
			}
		}

	}
}
