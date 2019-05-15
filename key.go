package vt100

import (
	"github.com/pkg/term"
)

type RawTerminal term.Term

// NewRawTerminal opens /dev/tty in raw mode as a term.Term
func NewRawTerminal() *RawTerminal {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	var r RawTerminal
	r = RawTerminal(*t)
	return &r
}


// Term will return the underlying term.Term
func (r *RawTerminal) Term() *term.Term {
	var t term.Term
	t = term.Term(*r)
	return &t
}

// RawMode will switch the terminal to raw mode
func (r *RawTerminal) RawMode() {
	term.RawMode(r.Term())
}

// Restore will restore the terminal
func (r *RawTerminal) Restore() {
	r.Term().Restore()
}

// Close will Restore and close the raw terminal
func (r *RawTerminal) Close() {
	t := r.Term()
	t.Restore()
	t.Close()
}

// Thanks https://stackoverflow.com/a/32018700/131264
// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func asciiAndKeyCode(r *RawTerminal) (ascii, keyCode int, err error) {
	bytes := make([]byte, 3)
	var numRead int
	r.RawMode()
	numRead, err = r.Term().Read(bytes)
	r.Restore()
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	return
}

// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func asciiAndKeyCodeOnce() (ascii, keyCode int, err error) {
	t := NewRawTerminal()
	a, kc, err := asciiAndKeyCode(t)
	t.Close()
	return a, kc, err
}

func (r *RawTerminal) ASCII() int {
	ascii, _, err := asciiAndKeyCode(r)
	if err != nil {
		return 0
	}
	return ascii
}

func ASCIIOnce() int {
	ascii, _, err := asciiAndKeyCodeOnce()
	if err != nil {
		return 0
	}
	return ascii
}

func (r *RawTerminal) KeyCode() int {
	_, keyCode, err := asciiAndKeyCode(r)
	if err != nil {
		return 0
	}
	return keyCode
}

func KeyCodeOnce() int {
	_, keyCode, err := asciiAndKeyCodeOnce()
	if err != nil {
		return 0
	}
	return keyCode
}

func (r *RawTerminal) Key() int {
	ascii, keyCode, err := asciiAndKeyCode(r)
	if err != nil {
		return 0
	}
	if keyCode != 0 {
		return keyCode
	}
	return ascii
}

func KeyOnce() int {
	ascii, keyCode, err := asciiAndKeyCodeOnce()
	if err != nil {
		return 0
	}
	if keyCode != 0 {
		return keyCode
	}
	return ascii
}
