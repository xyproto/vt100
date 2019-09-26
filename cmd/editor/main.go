package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/vt100"
	"os"
	"strings"
	"time"
)

type StatusBar struct {
	msg string
	fg  vt100.AttributeColor
	bg  vt100.AttributeColor
}

func NewStatusBar() *StatusBar {
	return &StatusBar{"", vt100.LightGreen, vt100.BackgroundDefault}
}

func (sb *StatusBar) Draw(c *vt100.Canvas) {
	w := int(c.W())
	c.Write(uint((w-len(sb.msg))/2), c.H()-1, sb.fg, sb.bg, sb.msg)
}

func (sb *StatusBar) SetMessage(msg string) {
	sb.msg = msg
}

func (sb *StatusBar) Clear() {
	sb.msg = strings.Repeat(" ", len(sb.msg))
}

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
	}
	if cur.X >= int(w) {
		cur.X = 0
		cur.Y++
	}
	if cur.Y >= int(h) {
		cur.Y = int(h)
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

	e := NewEditor(4) // 4 spaces per tab
	if filename != "" {
		e.Load(filename)
		// Draw editor lines from line 0 up to h onto the canvas at 0,0
		h := int(c.Height())
		e.WriteLines(c, 0, h, 0, 0)
	}

	scrolled := false
	offset := 0

	status := &StatusBar{}
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
			atStart := 0 == dataCursor.X
			atDocumentStart := 0 == dataCursor.X && 0 == dataCursor.Y
			if !atDocumentStart {
				// Move the data cursor
				if atStart {
					dataCursor.Y--
					dataCursor.X = e.LastDataPosition(int(dataCursor.Y))
				} else {
					dataCursor.X--
				}
				dataCursor.Wrap(c)
				if atStart {
					screenCursor.Y--
					screenCursor.X = e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab))
				} else {
					// Check if we hit a tab character
					atTab := '\t' == e.Get(int(dataCursor.X), int(dataCursor.Y))
					// Move the screen cursor
					if atTab && screenCursor.X >= e.spacesPerTab {
						screenCursor.X -= e.spacesPerTab
					} else {
						screenCursor.X--
					}
				}
				screenCursor.Wrap(c)
			}
		case 39: // right arrow
			atTab := '\t' == e.Get(int(dataCursor.X), int(dataCursor.Y))
			atEnd := dataCursor.X >= e.LastDataPosition(int(dataCursor.Y))
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
				if atTab && screenCursor.X < (int(c.Width())-e.spacesPerTab) {
					screenCursor.X += e.spacesPerTab
				} else {
					screenCursor.X++
				}
				screenCursor.Wrap(c)
			}
		case 38: // up arrow
			// Move the data cursor
			dataCursor.Y--
			dataCursor.Wrap(c)
			// Move the screen cursor
			if screenCursor.Y == 0 {
				// If at the top, don't move up, but scroll the contents
				status.Clear()
				status.Draw(c)
				status.SetMessage("reached top of screen")
				status.Draw(c)
			} else {
				screenCursor.Y--
				screenCursor.Wrap(c)
			}
		case 40: // down arrow
			// Move the data cursor
			dataCursor.Y++
			dataCursor.Wrap(c)
			// Move the screen cursor
			if screenCursor.Y == int(c.H()-1) {
				// If at the bottom, don't move down, but scroll the contents
				status.Clear()
				status.Draw(c)
				status.SetMessage("reached bottom of screen")
				status.Draw(c)
			} else {
				screenCursor.Y++
				screenCursor.Wrap(c)
			}
		case 14: // ctrl-n, scroll down
			if dataCursor.Y >= e.Len()-int(c.H()) {
				// Status message
				status.Clear()
				status.Draw(c)
				status.SetMessage("EOF")
				status.Draw(c)
				c.Draw()
			} else {
				// Only move the data cursor down one line, not the screen cursor
				dataCursor.Y++
				dataCursor.Wrap(c)
				// Move the scroll offset
				offset++
				// Prepare to redraw
				scrolled = true
			}
		case 16: // ctrl-p, scroll up
			if dataCursor.Y == 0 {
				// Can't scroll further up
				// Status message
				status.Clear()
				status.Draw(c)
				status.SetMessage("at top")
				status.Draw(c)
				c.Draw()
			} else {
				// Only move the data cursor up one line, not the screen cursor
				dataCursor.Y--
				dataCursor.Wrap(c)
				// Move the scroll offset
				offset--
				// Prepare to redraw
				scrolled = true
			}
		default:
			if key == 32 { // space
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, ' ')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.LightYellow, vt100.BackgroundBlue, ' ')
				screenCursor.X++
				screenCursor.Wrap(c)
			} else if key == 13 { // return
				// Data cursor
				dataCursor.Y++
				dataCursor.X = 0
				e.CreateLineIfMissing(dataCursor.Y)
				dataCursor.Wrap(c)
				// Screen cursor
				screenCursor.Y++
				screenCursor.X = 0
				screenCursor.Wrap(c)
			} else if (key >= 'a' && key <= 'z') || (key >= 'A' && key <= 'Z') { // letter
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, rune(key))
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.LightYellow, vt100.BackgroundBlue, rune(key))
				screenCursor.X++
				screenCursor.Wrap(c)
			} else if key == 127 { // backspace
				atTab := '\t' == e.Get(int(dataCursor.X), int(dataCursor.Y))
				// Data cursor
				if dataCursor.X == 0 {
					dataCursor.Y--
				} else {
					dataCursor.X--
				}
				dataCursor.Wrap(c)
				e.Set(int(dataCursor.X), int(dataCursor.Y), ' ')
				// Screen cursor
				if screenCursor.X == 0 {
					screenCursor.Y--
				} else {
					if atTab {
						screenCursor.X -= e.spacesPerTab
					} else {
						screenCursor.X--
					}
				}
				screenCursor.Wrap(c)
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.LightYellow, vt100.BackgroundBlue, ' ')
			} else if key == 9 { // tab
				// Data cursor
				e.Set(int(dataCursor.X), int(dataCursor.Y), '\t')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.Write(uint(screenCursor.X), uint(screenCursor.Y), vt100.LightYellow, vt100.BackgroundBlue, "    ")
				screenCursor.X += e.spacesPerTab
				screenCursor.Wrap(c)
			} else if key == 1 { // ctrl-a, home
				dataCursor.X = 0
				screenCursor.X = 0
			} else if key == 5 { // ctrl-e, end
				dataCursor.X = int(e.LastDataPosition(int(dataCursor.Y)))
				screenCursor.X = int(e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab)))
			} else if key == 19 { // ctrl-s, save
				err := e.Save(filename)
				if err != nil {
					tty.Close()
					vt100.Close()
					fmt.Fprintln(os.Stderr, vt100.Red.Get(err.Error()))
					os.Exit(1)
				}
			} else if key != 0 { // any other key
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), vt100.LightYellow, vt100.BackgroundBlue, rune(key))
				e.Set(screenCursor.X, screenCursor.Y, rune(key))
				screenCursor.X++
				screenCursor.Wrap(c)
			}
		}
		if scrolled {
			// redraw all characters
			h := int(c.Height())
			e.WriteLines(c, 0+offset, h+offset, 0, 0)
			c.Draw()
		} else if e.Changed() {
			c.Draw()
		}
		vt100.SetXY(uint(screenCursor.X), uint(screenCursor.Y))
	}
	tty.Close()
	vt100.Close()
	vt100.Clear()
	vt100.LightBlue.Output("bye!")
}
