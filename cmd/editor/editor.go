package main

import (
	"bytes"
	"errors"
	"github.com/xyproto/vt100"
	"io/ioutil"
	"strings"
)

type Editor struct {
	lines        map[int][]rune
	insertMode   bool
	changed      bool
	fg           vt100.AttributeColor
	bg           vt100.AttributeColor
	spacesPerTab int
}

func NewEditor(spacesPerTab int) *Editor {
	e := &Editor{}
	e.lines = make(map[int][]rune)
	e.insertMode = true
	e.fg = vt100.LightYellow
	e.bg = vt100.BackgroundDefault
	e.spacesPerTab = spacesPerTab
	return e
}

func (e *Editor) InsertMode() bool {
	return e.insertMode
}

func (e *Editor) ToggleInsertMode() {
	e.insertMode = !e.insertMode
}

func (e *Editor) Set(x, y int, r rune) {
	if e.lines == nil {
		e.lines = make(map[int][]rune)
	}
	_, ok := e.lines[y]
	if !ok {
		e.lines[y] = make([]rune, 0, x+1)
	}
	if x < int(len(e.lines[y])) {
		e.lines[y][x] = r
		e.changed = true
		return
	}
	// If the line is too short, fill it up with spaces
	for x >= int(len(e.lines[y])) {
		e.lines[y] = append(e.lines[y], ' ')
	}
	e.lines[y][x] = r
	e.changed = true
}

func (e *Editor) Get(x, y int) rune {
	if e.lines == nil {
		return ' '
	}
	runes, ok := e.lines[y]
	if !ok {
		return ' '
	}
	if x >= int(len(runes)) {
		return ' '
	}
	return runes[x]
}

func (e *Editor) Changed() bool {
	return e.changed
}

// Line returns the contents of line number N, counting from 0
func (e *Editor) Line(n int) string {
	line, ok := e.lines[n]
	if ok {
		var sb strings.Builder
		for _, r := range line {
			sb.WriteRune(r)
		}
		return sb.String()
	}
	return ""
}

// LastDataPosition returns the last X index for this line, for the data (does not expand tabs)
// Can be negative, if the line is empty.
func (e *Editor) LastDataPosition(n int) int {
	return len(e.Line(n)) - 1
}

// LastScreenPosition returns the last X index for this line, for the screen (expands tabs)
// Can be negative, if the line is empty.
func (e *Editor) LastScreenPosition(n, spacesPerTab int) int {
	extraSpaceBecauseOfTabs := int(e.Count(n, '\t') * (spacesPerTab - 1))
	return e.LastDataPosition(n) + extraSpaceBecauseOfTabs
}

// For a given line index, count the number of given runes
func (e *Editor) Count(n int, r rune) int {
	var counter int
	line, ok := e.lines[n]
	if ok {
		for _, l := range line {
			if l == r {
				counter++
			}
		}
	}
	return counter
}

// Len returns the number of lines
func (e *Editor) Len() int {
	maxy := 0
	for y := range e.lines {
		if y > maxy {
			maxy = y
		}
	}
	return maxy + 1
}

// String returns the contents of the editor
func (e *Editor) String() string {
	var sb strings.Builder
	for i := 0; i < e.Len(); i++ {
		sb.WriteString(e.Line(i) + "\n")
	}
	return sb.String()
}

func (e *Editor) Clear() {
	e.lines = make(map[int][]rune)
}

func (e *Editor) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	datalines := bytes.Split(data, []byte{'\n'})
	e.Clear()
	for y, dataline := range datalines {
		line := string(dataline)
		for x, letter := range line {
			e.Set(int(x), int(y), letter)
		}
	}
	return nil
}

func (e *Editor) Save(filename string) error {
	return ioutil.WriteFile(filename, []byte(e.String()), 0664)
}

// Write editor lines from "fromline" to and up to "toline" to the canvas at cx, cy
func (e *Editor) WriteLines(c *vt100.Canvas, fromline, toline, cx, cy int) error {
	w := int(c.Width())
	if fromline >= toline {
		return errors.New("fromline >= toline in WriteLines")
	}
	numlines := toline - fromline
	offset := fromline
	for y := 0; y < numlines; y++ {
		counter := 0
		for _, letter := range e.Line(y + offset) {
			if letter == '\t' {
				c.Write(uint(cx+counter), uint(cy+y), vt100.White, vt100.BackgroundBlue, "    ")
				counter += 4
			} else {
				c.WriteRune(uint(cx+counter), uint(cy+y), vt100.White, vt100.BackgroundBlue, letter)
				counter++
			}
		}
		// Fill the rest of the line on the canvas with "blanks"
		for x := counter; x < w; x++ {
			c.WriteRune(uint(cx+x), uint(cy+y), vt100.White, vt100.BackgroundBlue, ' ')
		}
	}
	return nil
}

func (e *Editor) CreateLineIfMissing(n int) {
	if e.lines == nil {
		e.lines = make(map[int][]rune)
	}
	_, ok := e.lines[n]
	if !ok {
		e.lines[n] = make([]rune, 0)
	}
}
