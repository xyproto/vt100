package main

import (
	"bytes"
	"errors"
	"github.com/xyproto/vt100"
	"io/ioutil"
	"strings"
)

type Editor struct {
	lines      map[uint][]rune
	insertMode bool
	changed    bool
}

func NewEditor() *Editor {
	e := &Editor{}
	e.lines = make(map[uint][]rune)
	e.insertMode = true
	return e
}

func (e *Editor) InsertMode() bool {
	return e.insertMode
}

func (e *Editor) ToggleInsertMode() {
	e.insertMode = !e.insertMode
}

func (e *Editor) Set(x, y uint, r rune) {
	if e.lines == nil {
		e.lines = make(map[uint][]rune)
	}
	_, ok := e.lines[y]
	if !ok {
		e.lines[y] = make([]rune, 0, x+1)
	}
	if x < uint(len(e.lines[y])) {
		e.lines[y][x] = r
		e.changed = true
		return
	}
	// If the line is too short, fill it up with spaces
	for x >= uint(len(e.lines[y])) {
		e.lines[y] = append(e.lines[y], ' ')
	}
	e.lines[y][x] = r
	e.changed = true
}

func (e *Editor) Get(x, y uint) rune {
	if e.lines == nil {
		return ' '
	}
	runes, ok := e.lines[y]
	if !ok {
		return ' '
	}
	if x >= uint(len(runes)) {
		return ' '
	}
	return runes[x]
}

func (e *Editor) Changed() bool {
	return e.changed
}

// Line returns the contents of line number N, counting from 0
func (e *Editor) Line(n uint) string {
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
func (e *Editor) LastDataPosition(n uint) int {
	return len(e.Line(n)) - 1
}

// LastScreenPosition returns the last X index for this line, for the screen (expands tabs)
// Can be negative, if the line is empty.
func (e *Editor) LastScreenPosition(n, spacesPerTab uint) int {
	extraSpaceBecauseOfTabs := int(e.Count(n, '\t') * (spacesPerTab - 1))
	return e.LastDataPosition(n) + extraSpaceBecauseOfTabs
}

// For a given line index, count the number of given runes
func (e *Editor) Count(n uint, r rune) uint {
	var counter uint
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
func (e *Editor) Len() uint {
	maxy := uint(0)
	for y, _ := range e.lines {
		if y > maxy {
			maxy = y
		}
	}
	return maxy + 1
}

// String returns the contents of the editor
func (e *Editor) String() string {
	var sb strings.Builder
	for i := uint(0); i < e.Len(); i++ {
		sb.WriteString(e.Line(i) + "\n")
	}
	return sb.String()
}

func (e *Editor) Clear() {
	e.lines = make(map[uint][]rune)
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
			e.Set(uint(x), uint(y), letter)
		}
	}
	return nil
}

func (e *Editor) Save(filename string) error {
	return ioutil.WriteFile(filename, []byte(e.String()), 0664)
}

// Write editor lines from "fromline" to and up to "toline" to the canvas at cx, cy
func (e *Editor) WriteLines(c *vt100.Canvas, fromline, toline, cx, cy uint) error {
	w, _ := c.Size()
	if fromline >= toline {
		return errors.New("fromline >= toline in WriteLines")
	}
	for y := fromline; y < toline; y++ {
		counter := uint(0)
		for _, letter := range e.Line(y) {
			if letter == '\t' {
				c.Write(cx+counter, cy+y, vt100.White, vt100.BackgroundBlue, "    ")
				counter += 4
			} else {
				c.WriteRune(cx+counter, cy+y, vt100.White, vt100.BackgroundBlue, letter)
				counter++
			}
		}
		// Fill the rest of the line on the canvas with "blanks"
		for x := counter; x < w; x++ {
			c.WriteRune(cx+x, cy+y, vt100.White, vt100.BackgroundBlue, ' ')
		}
	}
	return nil
}

func (e *Editor) CreateLineIfMissing(n uint) {
	if e.lines == nil {
		e.lines = make(map[uint][]rune)
	}
	_, ok := e.lines[n]
	if !ok {
		e.lines[n] = make([]rune, 0)
	}
}
