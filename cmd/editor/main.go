package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/vt100"
	"os"
	"time"
)

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

	cursor := &Cursor{}

	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	tty.SetTimeout(10 * time.Millisecond)
	quit := false
	for !quit {
		key := tty.Key()
		switch key {
		case 17:
			quit = true
		case 37:
			if cursor.X >= 4 && e.Get(uint(cursor.X-4), uint(cursor.Y)) == '\t' {
				cursor.X -= 4
			} else {
				cursor.X--
			}
			cursor.Wrap(c)
		case 38:
			cursor.Y--
			cursor.Wrap(c)
		case 39:
			if e.Get(uint(cursor.X), uint(cursor.Y)) == '\t' {
				cursor.X += 4
			} else {
				cursor.X++
			}
			cursor.Wrap(c)
		case 40:
			cursor.Y++
			cursor.Wrap(c)
		default:
			if key == 32 {
				// space
				c.WriteRune(uint(cursor.X), uint(cursor.Y), vt100.White, vt100.BackgroundBlue, ' ')
				e.Set(uint(cursor.X), uint(cursor.Y), ' ')
				cursor.X++
				cursor.Wrap(c)
			} else if key == 13 {
				// return
				cursor.Y++
				cursor.X = 0
				e.CreateLineIfMissing(uint(cursor.Y))
				cursor.Wrap(c)
			} else if (key >= 'a' && key <= 'z') || (key >= 'A' && key <= 'Z') {
				// letter
				c.WriteRune(uint(cursor.X), uint(cursor.Y), vt100.White, vt100.BackgroundBlue, rune(key))
				e.Set(uint(cursor.X), uint(cursor.Y), rune(key))
				cursor.X++
				cursor.Wrap(c)
			} else if key == 127 {
				// backspace
				if cursor.X == 0 {
					cursor.Y--
				} else {
					if e.Get(uint(cursor.X), uint(cursor.Y)) == '\t' {
						cursor.X -= 4
					} else {
						cursor.X--
					}
				}
				cursor.Wrap(c)
				c.WriteRune(uint(cursor.X), uint(cursor.Y), vt100.White, vt100.BackgroundBlue, ' ')
				e.Set(uint(cursor.X), uint(cursor.Y), ' ')
			} else if key == 9 {
				// tab
				c.Write(uint(cursor.X), uint(cursor.Y), vt100.White, vt100.BackgroundBlue, "    ")
				e.Set(uint(cursor.X), uint(cursor.Y), '\t')
				cursor.X += 4
				cursor.Wrap(c)
			} else if key == 19 {
				// ctrl-s, save
				err := e.Save(filename)
				if err != nil {
					tty.Close()
					vt100.Close()
					fmt.Fprintln(os.Stderr, vt100.Red.Get(err.Error()))
					os.Exit(1)
				}
			} else if key != 0 {
				// any other key
				c.WriteRune(uint(cursor.X), uint(cursor.Y), vt100.White, vt100.BackgroundBlue, rune(key))
				e.Set(uint(cursor.X), uint(cursor.Y), rune(key))
				cursor.X++
				cursor.Wrap(c)
			}
		}
		if e.Changed() {
			c.Draw()
		}
		vt100.SetXY(uint(cursor.X), uint(cursor.Y))
	}
	tty.Close()
	vt100.Close()
	vt100.BackgroundDefault.Combine(vt100.LightBlue).Output("bye!")
}
