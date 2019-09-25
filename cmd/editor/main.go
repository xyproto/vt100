package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/vt100"
	"os"
	"time"
)

const spacesPerTab = 4

type Cursor struct {
	X, Y int
}

func (cur *Cursor) Wrap(c *vt100.Canvas) {
	w, h := c.Size()
	if cur.X < 0 {
		cur.X = int(w)
		cur.Y--
	}
	if cur.Y < 0 {
		cur.Y = 0
		c.Write(20, 20, vt100.LightGreen, vt100.Default, "SCROLL UP   ")
	}
	if cur.X >= int(w) {
		cur.X = 0
		cur.Y++
	}
	if cur.Y >= int(h) {
		cur.Y = int(h)
		c.Write(20, 20, vt100.LightGreen, vt100.Default, "SCROLL DOWN")
	}
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Fprintln(os.Stderr, vt100.Red.Get("a filename must be given"))
		os.Exit(1)
	}

	vt100.Init()
	vt100.ShowCursor(true)

	c := vt100.NewCanvas()
	c.FillBackground(vt100.Blue)

	e := NewEditor()
	if filename != "" {
		e.Load(filename)
		// Draw editor lines from line 0 up to h onto the canvas at 0,0
		_, h := c.Size()
		e.WriteLines(c, 0, h, 0, 0)
	}

	c.Draw()

	screenCursor := &Cursor{}
	dataCursor := &Cursor{}

	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	tty.SetTimeout(5 * time.Millisecond)
	quit := false
	for !quit {
		key := tty.Key()
		switch key {
		case 27: // esc
			e.ToggleInsertMode()
		case 17: // ctrl-q
			quit = true
		case 37: // left arrow
			atTab := '\t' == e.Get(uint(dataCursor.X), uint(dataCursor.Y))
			// Move the data cursor
			dataCursor.X--
			dataCursor.Wrap(c)
			// Move the screen cursor
			if atTab && screenCursor.X >= spacesPerTab {
				screenCursor.X -= spacesPerTab
			} else {
				screenCursor.X--
			}
			screenCursor.Wrap(c)
		case 38: // up arrow
			// Move the data cursor
			dataCursor.Y--
			dataCursor.Wrap(c)
			// Move the screen cursor
			screenCursor.Y--
			screenCursor.Wrap(c)
		case 39: // right arrow
			atTab := '\t' == e.Get(uint(dataCursor.X), uint(dataCursor.Y))
			atEnd := dataCursor.X >= e.LastPosition(uint(dataCursor.Y))
			if atEnd {
				// Move the data cursor
				dataCursor.X = 0
				dataCursor.Y++
				dataCursor.Wrap(c)
				// Move the screen cursor
				screenCursor.X = 0
				screenCursor.Y++
				screenCursor.Wrap(c)
			} else {
				// Move the data cursor
				dataCursor.X++
				dataCursor.Wrap(c)
				// Move the screen cursor
				if atTab && uint(screenCursor.X) < (c.Width()-spacesPerTab) {
					screenCursor.X += spacesPerTab
				} else {
					screenCursor.X++
				}
				screenCursor.Wrap(c)
			}
		case 40: // down arrow
			// Move the data cursor
			dataCursor.Y++
			dataCursor.Wrap(c)
			// Move the screen cursor
			screenCursor.Y++
			screenCursor.Wrap(c)
		default:
			if key == 32 { // space
				// Data cursor
				e.Set(uint(dataCursor.X), uint(dataCursor.Y), ' ')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.White, vt100.BackgroundBlue, ' ')
				screenCursor.X++
				screenCursor.Wrap(c)
			} else if key == 13 { // return
				// Data cursor
				dataCursor.Y++
				dataCursor.X = 0
				e.CreateLineIfMissing(uint(dataCursor.Y))
				dataCursor.Wrap(c)
				// Screen cursor
				screenCursor.Y++
				screenCursor.X = 0
				screenCursor.Wrap(c)
			} else if (key >= 'a' && key <= 'z') || (key >= 'A' && key <= 'Z') { // letter
				// Data cursor
				e.Set(uint(dataCursor.X), uint(dataCursor.Y), rune(key))
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.White, vt100.BackgroundBlue, rune(key))
				screenCursor.X++
				screenCursor.Wrap(c)
			} else if key == 127 { // backspace
				atTab := '\t' == e.Get(uint(dataCursor.X), uint(dataCursor.Y))
				// Data cursor
				if dataCursor.X == 0 {
					dataCursor.Y--
				} else {
					dataCursor.X--
				}
				dataCursor.Wrap(c)
				e.Set(uint(dataCursor.X), uint(dataCursor.Y), ' ')
				// Screen cursor
				if screenCursor.X == 0 {
					screenCursor.Y--
				} else {
					if atTab {
						screenCursor.X -= spacesPerTab
					} else {
						screenCursor.X--
					}
				}
				screenCursor.Wrap(c)
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.White, vt100.BackgroundBlue, ' ')
			} else if key == 9 { // tab
				// Data cursor
				e.Set(uint(dataCursor.X), uint(dataCursor.Y), '\t')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.Write(uint(screenCursor.X), uint(screenCursor.Y), vt100.White, vt100.BackgroundBlue, "    ")
				screenCursor.X += spacesPerTab
				screenCursor.Wrap(c)
			} else if key == 1 { // ctrl-a, home
				dataCursor.X = 0
				screenCursor.X = 0
			} else if key == 5 { // ctrl-e, end
				lastPosition := e.LastPosition(uint(dataCursor.Y))
				dataCursor.X = int(lastPosition)
				tabCount := e.Count(uint(dataCursor.Y), '\t')
				screenCursor.X = int(tabCount)*int(spacesPerTab) + (lastPosition - int(tabCount))
			} else if key == 19 { // ctrl-s, save
				err := e.Save(filename)
				if err != nil {
					tty.Close()
					vt100.Close()
					fmt.Fprintln(os.Stderr, vt100.Red.Get(err.Error()))
					os.Exit(1)
				}
			} else if key != 0 { // any other key
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.White, vt100.BackgroundBlue, rune(key))
				e.Set(uint(screenCursor.X), uint(screenCursor.Y), rune(key))
				screenCursor.X++
				screenCursor.Wrap(c)
			}
		}
		if e.Changed() {
			c.Draw()
		}
		vt100.SetXY(uint(screenCursor.X), uint(screenCursor.Y))
	}
	tty.Close()
	vt100.Close()
	vt100.BackgroundDefault.Combine(vt100.LightBlue).Output("bye!")
}
