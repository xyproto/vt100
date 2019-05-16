package vt100

import (
	"github.com/pkg/term"
	"time"
)

type TTY term.Term

// NewTTY opens /dev/tty in raw and cbreak mode as a term.Term
func NewTTY() (*TTY, error) {
	t, err := term.Open("/dev/tty", term.RawMode, term.CBreakMode, term.ReadTimeout(20*time.Millisecond))
	if err != nil {
		return nil, err
	}
	tty := TTY(*t)
	return &tty, nil
}

// Term will return the underlying term.Term
func (tty *TTY) Term() *term.Term {
	var t term.Term
	t = term.Term(*tty)
	return &t
}

// RawMode will switch the terminal to raw mode
func (tty *TTY) RawMode() {
	term.RawMode(tty.Term())
}

// NoBlock leaves "cooked" mode and enters "cbreak" mode
func (tty *TTY) NoBlock() {
	tty.Term().SetCbreak()
}

// Timeout sets a timeout for reading a key
func (tty *TTY) Timeout(d time.Duration) {
	tty.Term().SetReadTimeout(d)
}

// Restore will restore the terminal
func (tty *TTY) Restore() {
	tty.Term().Restore()
}

// Close will Restore and close the raw terminal
func (tty *TTY) Close() {
	t := tty.Term()
	t.Restore()
	t.Close()
}

// Thanks https://stackoverflow.com/a/32018700/131264
// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func asciiAndKeyCode(tty *TTY) (ascii, keyCode int, err error) {
	takes := 20 * time.Millisecond
	bytes := make([]byte, 3)
	var numRead int
	tty.RawMode()
	tty.NoBlock()
	tty.Timeout(takes)
	numRead, err = tty.Term().Read(bytes)
	tty.Restore()
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
	t, err := NewTTY()
	if err != nil {
		return 0, 0, err
	}
	a, kc, err := asciiAndKeyCode(t)
	t.Close()
	return a, kc, err
}

func (tty *TTY) ASCII() int {
	ascii, _, err := asciiAndKeyCode(tty)
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

func (tty *TTY) KeyCode() int {
	_, keyCode, err := asciiAndKeyCode(tty)
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

func (tty *TTY) Key() int {
	ascii, keyCode, err := asciiAndKeyCode(tty)
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
