package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/vt100"
	"os"
	"time"
)

const versionString = "rED 1.0.0"

func main() {
	var (
		// These are used for initializing various structs
		defaultEditorForeground       = vt100.LightGreen
		defaultEditorBackground       = vt100.BackgroundBlack
		defaultEditorStatusForeground = vt100.Black
		defaultEditorStatusBackground = vt100.BackgroundGray

		defaultASCIIGraphicsForeground       = vt100.LightYellow
		defaultASCIIGraphicsBackground       = vt100.BackgroundBlue
		defaultASCIIGraphicsStatusForeground = vt100.White
		defaultASCIIGraphicsStatusBackground = vt100.BackgroundMagenta

		statusDuration = 3000 * time.Millisecond

		offset = 0
		redraw = false
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
	e := NewEditor(4, 10, defaultEditorForeground, defaultEditorBackground)

	status := NewStatusBar(defaultEditorStatusForeground, defaultEditorStatusBackground, e, statusDuration)

	// Try to load the filename, ignore errors since giving a new filename is also okay
	// TODO: Check if the file exists and add proper error reporting
	err := e.Load(filename)
	loaded := err == nil

	// Draw editor lines from line 0 up to h onto the canvas at 0,0
	h := int(c.Height())
	e.WriteLines(c, 0, h, 0, 0)

	// Friendly status message
	if loaded {
		status.SetMessage("Loaded " + filename)
	} else {
		status.SetMessage(versionString)
	}
	status.Show(c, offset)
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
			e.ToggleEOLMode()
			if e.EOLMode() {
				e.SetColors(defaultEditorForeground, defaultEditorBackground)
				status.SetColors(defaultEditorStatusForeground, defaultEditorStatusBackground)
				c.FillBackground(e.bg)
				status.SetMessage("Text edit mode")
				redraw = true
			} else {
				e.SetColors(defaultASCIIGraphicsForeground, defaultASCIIGraphicsBackground)
				status.SetColors(defaultASCIIGraphicsStatusForeground, defaultASCIIGraphicsStatusBackground)
				c.FillBackground(e.bg)
				status.SetMessage("ASCII graphics mode")
				redraw = true
			}
		case 17: // ctrl-q, quit
			quit = true
		case 7: // ctrl-g, status information
			currentRune := e.Get(dataCursor.X, dataCursor.Y)
			if currentRune > 32 {
				status.SetMessage(fmt.Sprintf("%d,%d (data %d,%d) %c (%U) wordcount: %d", screenCursor.X, screenCursor.Y, dataCursor.X, dataCursor.Y, currentRune, currentRune, e.WordCount()))
			} else {
				status.SetMessage(fmt.Sprintf("%d,%d (data %d,%d) %U wordcount: %d", screenCursor.X, screenCursor.Y, dataCursor.X, dataCursor.Y, currentRune, e.WordCount()))
			}
			status.Show(c, offset)
		case 37: // left arrow
			atStart := 0 == dataCursor.X
			atDocumentStart := 0 == dataCursor.X && 0 == dataCursor.Y
			if !atDocumentStart {
				// Move the data cursor
				if atStart {
					dataCursor.Y--
					if e.EOLMode() {
						dataCursor.X = e.LastDataPosition(dataCursor.Y)
					}
				} else {
					dataCursor.X--
				}
				dataCursor.ScreenWrap(c)
				if atStart {
					screenCursor.Y--
					if e.EOLMode() {
						screenCursor.X = e.LastScreenPosition(dataCursor.Y, int(e.spacesPerTab)) + 1
					}
				} else {
					// Check if we hit a tab character
					atTab := '\t' == e.Get(dataCursor.X, dataCursor.Y)
					// Move the screen cursor
					if atTab && screenCursor.X >= e.spacesPerTab {
						screenCursor.X -= e.spacesPerTab
					} else {
						screenCursor.X--
					}
				}
				screenCursor.ScreenWrap(c)
			}
		case 39: // right arrow
			atTab := '\t' == e.Get(dataCursor.X, dataCursor.Y)
			atEnd := dataCursor.X >= e.LastDataPosition(dataCursor.Y)
			if atEnd && e.EOLMode() {
				// Move the data cursor
				dataCursor.X = 0

				if dataCursor.Y < e.Len() {
					dataCursor.Y++
					screenCursor.Y++
				} else {
					status.SetMessage("End of text")
					status.Show(c, offset)
				}
				dataCursor.ScreenWrap(c)
				// Move the screen cursor
				screenCursor.X = 0
				screenCursor.ScreenWrap(c)
			} else {
				// Move the data cursor
				dataCursor.X++
				dataCursor.ScreenWrap(c)
				// Move the screen cursor
				if atTab && screenCursor.X < (int(c.Width())-e.spacesPerTab) {
					screenCursor.X += e.spacesPerTab
				} else {
					screenCursor.X++
				}
				screenCursor.ScreenWrap(c)
			}
		case 38: // up arrow
			// Move the screen cursor
			if screenCursor.Y == 0 {
				// If at the top, don't move up, but scroll the contents
				//redraw, offset = scrollUp(c, offset, status, e, dataCursor, 1)
				// Output a helpful message
				if dataCursor.Y == 0 {
					status.SetMessage("Start of text")
				} else {
					status.SetMessage("Top of screen, scroll with ctrl-p")
				}
				status.Show(c, offset)
			} else {
				// Move the data cursor
				dataCursor.Y--
				dataCursor.ScreenWrap(c)
				// Move the screen cursor
				screenCursor.Y--
				screenCursor.ScreenWrap(c)
			}
			// If the cursor is after the length of the current line, move it to the end of the current line
			if e.EOLMode() {
				if dataCursor.X > e.LastDataPosition(dataCursor.Y) {
					dataCursor.X = int(e.LastDataPosition(dataCursor.Y)) + 1
				}
				if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
					screenCursor.X = int(e.LastScreenPosition(dataCursor.Y, int(e.spacesPerTab))) + 1
				}
			}
		case 40: // down arrow
			if !e.EOLMode() || (e.EOLMode() && dataCursor.Y < e.Len()) {
				// Move the screen cursor
				if screenCursor.Y >= int(c.H()-1) {
					// If at the bottom, don't move down, but scroll the contents
					// redraw, offset = scrollDown(c, offset, status, e, dataCursor, 1)
					// Output a helpful message
					if dataCursor.Y == (e.Len() - 1) {
						status.SetMessage("End of text")
					} else {
						status.SetMessage("Bottom of screen, scroll with ctrl-n")
					}
					status.Show(c, offset)
				} else {
					// Move the data cursor
					dataCursor.Y++
					//dataCursor.ScreenWrap(c)
					// Move the screen cursor
					screenCursor.Y++
					//screenCursor.ScreenWrap(c)
				}
				// If the cursor is after the length of the current line, move it to the end of the current line
				if e.EOLMode() {
					if dataCursor.X > e.LastDataPosition(dataCursor.Y) {
						dataCursor.X = int(e.LastDataPosition(dataCursor.Y)) + 1
					}
					if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
						screenCursor.X = int(e.LastScreenPosition(dataCursor.Y, int(e.spacesPerTab))) + 1
					}
				}
			} else if e.EOLMode() {
				status.SetMessage("End of text")
				status.Show(c, offset)
			}
			// If the cursor is after the length of the current line, move it to the end of the current line
			if e.EOLMode() {
				if dataCursor.X > e.LastDataPosition(dataCursor.Y) {
					dataCursor.X = int(e.LastDataPosition(dataCursor.Y)) + 1
				}
				if screenCursor.X > e.LastScreenPosition(int(screenCursor.Y), int(e.spacesPerTab)) {
					screenCursor.X = int(e.LastScreenPosition(dataCursor.Y, int(e.spacesPerTab))) + 1
				}
			}
		case 14: // ctrl-n, scroll down
			redraw, offset = scrollDown(c, offset, status, e, dataCursor, e.scrollSpeed)
		case 16: // ctrl-p, scroll up
			redraw, offset = scrollUp(c, offset, status, e, dataCursor, e.scrollSpeed)
		default:
			if key == 32 { // space
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, ' ')
				dataCursor.X++
				dataCursor.ScreenWrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, ' ')
				screenCursor.X++
				screenCursor.ScreenWrap(c)
			} else if key == 13 { // return
				// Data cursor
				dataCursor.Y++
				dataCursor.X = 0
				e.CreateLineIfMissing(dataCursor.Y)
				dataCursor.ScreenWrap(c)
				// Screen cursor
				screenCursor.Y++
				screenCursor.X = 0
				screenCursor.ScreenWrap(c)
			} else if (key >= 'a' && key <= 'z') || (key >= 'A' && key <= 'Z') { // letter
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, rune(key))
				dataCursor.X++
				dataCursor.ScreenWrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, rune(key))
				screenCursor.X++
				screenCursor.ScreenWrap(c)
			} else if key == 127 { // backspace
				atTab := '\t' == e.Get(dataCursor.X, dataCursor.Y)
				// Data cursor
				if dataCursor.X > 0 {
					dataCursor.X--
					dataCursor.ScreenWrap(c)
					e.Set(dataCursor.X, dataCursor.Y, ' ')
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
					screenCursor.ScreenWrap(c)
					c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, ' ')
				}
			} else if key == 9 { // tab
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, '\t')
				dataCursor.X++
				dataCursor.ScreenWrap(c)
				// Screen cursor
				c.Write(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, "    ")
				screenCursor.X += e.spacesPerTab
				screenCursor.ScreenWrap(c)
			} else if key == 1 { // ctrl-a, home
				dataCursor.X = 0
				screenCursor.X = 0
			} else if key == 5 { // ctrl-e, end
				dataCursor.X = int(e.LastDataPosition(dataCursor.Y)) + 1
				screenCursor.X = int(e.LastScreenPosition(dataCursor.Y, int(e.spacesPerTab))) + 1
			} else if key == 19 { // ctrl-s, save
				err := e.Save(filename, true)
				if err != nil {
					tty.Close()
					vt100.Close()
					fmt.Fprintln(os.Stderr, vt100.Red.Get(err.Error()))
					os.Exit(1)
				}
				// Status message
				status.SetMessage("Saved " + filename)
				status.Show(c, offset)
				c.Draw()
				// Redraw after save, for syntax highlighting
				redraw = true
			} else if key == 12 { // ctrl-l, redraw
				redraw = true
			} else if key == 11 { // ctrl-k, delete to end of line
				e.DeleteRestOfLine(dataCursor.X, dataCursor.Y)
				vt100.Do("Erase End of Line")
				redraw = true
			} else if key != 0 { // any other key
				// Data cursor
				e.Set(dataCursor.X, dataCursor.Y, rune(key))
				dataCursor.X++
				dataCursor.ScreenWrap(c)
				// Screen cursor
				c.WriteRune(uint(screenCursor.X), uint(screenCursor.Y), e.fg, e.bg, rune(key))
				screenCursor.X++
				screenCursor.ScreenWrap(c)
			}
		}
		if redraw {
			// redraw all characters
			h := int(c.Height())
			e.WriteLines(c, 0+offset, h+offset, 0, 0)
			c.Draw()
			status.Show(c, offset)
			redraw = false
		} else if e.Changed() {
			c.Draw()
		}
		vt100.SetXY(uint(screenCursor.X), uint(screenCursor.Y))
	}
	tty.Close()
	vt100.Close()
}
