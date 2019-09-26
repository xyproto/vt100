package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/vt100"
	"os"
	"time"
)

func main() {

	var (
		// These are used for initializing various structs
		editorForeground = vt100.Blue
		editorBackground = vt100.BackgroundGray
		floatForeground  = vt100.Red
		floatBackground  = vt100.BackgroundGray
		statusForeground = vt100.Red
		statusBackground = vt100.BackgroundBlack

		statusDuration = 2000 * time.Millisecond
	)

	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		fmt.Fprintln(os.Stderr, "Please supply a filename.")
		os.Exit(1)
	}

	vt100.Init()
	vt100.ShowCursor(true)

	c := vt100.NewCanvas()

	// 4 spaces per tab, scroll 10 lines at a time
	e := NewEditor(4, 10, editorForeground, editorBackground)
	//c.FillBackground(vt100.BackgroundBlue)

	if filename != "" {
		e.Load(filename)
		// Draw editor lines from line 0 up to h onto the canvas at 0,0
		h := int(c.Height())
		e.WriteLines(c, 0, h, 0, 0)
	}

	redraw := false
	offset := 0

	status := NewStatusBar(statusForeground, statusBackground, e, statusDuration)
	c.Draw()

	status.SetMessage("ved 1.0.0")
	status.Show(c)

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
			e.ToggleEOLMode()
			if e.EOLMode() {
				e.fg = editorForeground
				e.bg = editorBackground
				c.FillBackground(e.bg)
				status.SetMessage("text edit mode")
				redraw = true
			} else {
				e.fg = floatForeground
				e.bg = floatBackground
				c.FillBackground(e.bg)
				status.SetMessage("ASCII graphics mode")
				redraw = true
			}
		case 17: // ctrl-q
			quit = true
		case 37: // left arrow
			atStart := 0 == dataCursor.X
			atDocumentStart := 0 == dataCursor.X && 0 == dataCursor.Y
			if !atDocumentStart {
				// Move the data cursor
				if atStart {
					dataCursor.Y--
					if e.EOLMode() {
						dataCursor.X = e.LastDataPosition(int(dataCursor.Y))
					}
				} else {
					dataCursor.X--
				}
				dataCursor.Wrap(c)
				if atStart {
					screenCursor.Y--
					if e.EOLMode() {
						screenCursor.X = e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab)) + 1
					}
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
			if atEnd && e.EOLMode() {
				// Move the data cursor
				dataCursor.X = 0

				if dataCursor.Y < e.Len() {
					dataCursor.Y++
					screenCursor.Y++
				} else {
					status.SetMessage("end of text")
					status.Show(c)
				}
				dataCursor.Wrap(c)
				// Move the screen cursor
				screenCursor.X = 0
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
			// Move the screen cursor
			if screenCursor.Y == 0 {
				// If at the top, don't move up, but scroll the contents
				status.SetMessage("top of screen")
				status.Show(c)
			} else {
				// Move the data cursor
				dataCursor.Y--
				dataCursor.Wrap(c)
				// Move the screen cursor
				screenCursor.Y--
				screenCursor.Wrap(c)
			}
			// If the cursor is after the length of the current line, move it to the end of the current line
			if e.EOLMode() {
				if dataCursor.X > e.LastDataPosition(int(dataCursor.Y)) {
					dataCursor.X = int(e.LastDataPosition(int(dataCursor.Y))) + 1
				}
				if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
					screenCursor.X = int(e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab))) + 1
				}
			}
		case 40: // down arrow
			if !e.EOLMode() || (e.EOLMode() && dataCursor.Y < e.Len()) {
				// Move the screen cursor
				if screenCursor.Y == int(c.H()-1) {
					// If at the bottom, don't move down, but scroll the contents
					status.SetMessage("bottom of screen")
					status.Show(c)
				} else {
					// Move the data cursor
					dataCursor.Y++
					dataCursor.Wrap(c)
					// Move the screen cursor
					screenCursor.Y++
					screenCursor.Wrap(c)
				}
				// If the cursor is after the length of the current line, move it to the end of the current line
				if e.EOLMode() {
					if dataCursor.X > e.LastDataPosition(int(dataCursor.Y)) {
						dataCursor.X = int(e.LastDataPosition(int(dataCursor.Y))) + 1
					}
					if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
						screenCursor.X = int(e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab))) + 1
					}
				}
			} else if e.EOLMode() {
				status.SetMessage("end of text")
				status.Show(c)
			}
			// If the cursor is after the length of the current line, move it to the end of the current line
			if e.EOLMode() {
				if dataCursor.X > e.LastDataPosition(int(dataCursor.Y)) {
					dataCursor.X = int(e.LastDataPosition(int(dataCursor.Y))) + 1
				}
				if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
					screenCursor.X = int(e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab))) + 1
				}
			}
		case 14: // ctrl-n, scroll down
			h := int(c.H())
			if offset >= e.Len()-h {
				// Status message
				status.SetMessage("EOF")
				status.Show(c)
				c.Draw()
			} else {
				// Find out if we can scroll e.scrollSpeed, or less
				canScroll := e.scrollSpeed
				if (offset + canScroll) >= (e.Len() - h) {
					// Almost at the bottom, we can scroll the remaining lines
					canScroll = (e.Len() - h) - offset
				}
				// Only move the data cursor down one or more lines, do not move the screen cursor
				dataCursor.Y += canScroll
				dataCursor.Wrap(c)
				// Move the scroll offset
				offset += canScroll
				// Prepare to redraw
				redraw = true
			}
		case 16: // ctrl-p, scroll up
			if offset == 0 {
				// Can't scroll further up
				// Status message
				status.SetMessage("at top")
				status.Show(c)
				c.Draw()
			} else {
				// Find out if we can scroll e.scrollSpeed, or less
				canScroll := e.scrollSpeed
				if offset-canScroll < 0 {
					// Almost at the top, we can scroll the remaining lines
					canScroll = offset
				}
				// Only move the data cursor up one or more lines, do not move the screen cursor
				dataCursor.Y -= canScroll
				dataCursor.Wrap(c)
				// Move the scroll offset
				offset -= canScroll
				// Prepare to redraw
				redraw = true
			}
		default:
			if key == 32 { // space
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, ' ')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, ' ')
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
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, rune(key))
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
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, ' ')
			} else if key == 9 { // tab
				// Data cursor
				e.Set(int(dataCursor.X), int(dataCursor.Y), '\t')
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.Write(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, "    ")
				screenCursor.X += e.spacesPerTab
				screenCursor.Wrap(c)
			} else if key == 1 { // ctrl-a, home
				dataCursor.X = 0
				screenCursor.X = 0
			} else if key == 5 { // ctrl-e, end
				dataCursor.X = int(e.LastDataPosition(int(dataCursor.Y))) + 1
				screenCursor.X = int(e.LastScreenPosition(int(dataCursor.Y), int(e.spacesPerTab))) + 1
			} else if key == 19 { // ctrl-s, save
				err := e.Save(filename, true)
				if err != nil {
					tty.Close()
					vt100.Close()
					fmt.Fprintln(os.Stderr, vt100.Red.Get(err.Error()))
					os.Exit(1)
				}
				// Status message
				status.SetMessage("wrote " + filename)
				status.Show(c)
				c.Draw()
			} else if key == 12 { // ctrl-l, redraw
				status.SetMessage("redraw")
				status.Show(c)
				c.Draw()
				redraw = true
			} else if key == 11 { // ctrl-k, delete to end of line
				e.DeleteRestOfLine(dataCursor.X, dataCursor.Y)
				vt100.Do("Erase End of Line")
				redraw = true
			} else if key != 0 { // any other key
				// Data cursor
				e.Set(screenCursor.X, screenCursor.Y, rune(key))
				dataCursor.X++
				dataCursor.Wrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, rune(key))
				screenCursor.X++
				screenCursor.Wrap(c)
			}
		}
		if redraw {
			// redraw all characters
			h := int(c.Height())
			e.WriteLines(c, 0+offset, h+offset, 0, 0)
			c.Draw()
			status.Show(c)
			redraw = false
		} else if e.Changed() {
			c.Draw()
		}
		vt100.SetXY(uint(screenCursor.X), uint(screenCursor.Y))
	}
	tty.Close()
	vt100.Close()
}
