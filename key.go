package vt100

import (
	"github.com/pkg/term"
)

// Thanks https://stackoverflow.com/a/32018700/131264

// Returns either an ascii code, or (if input is an arrow) a Javascript key code.
func asciiAndKeyCode() (ascii, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
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
	t.Restore()
	t.Close()
	return
}

func ASCII() int {
	ascii, _, err := asciiAndKeyCode()
	if err != nil {
		return 0
	}
	return ascii
}

func KeyCode() int {
	_, keyCode, err := asciiAndKeyCode()
	if err != nil {
		return 0
	}
	return keyCode
}

func Key() int {
	ascii, keyCode, err := asciiAndKeyCode()
	if err != nil {
		return 0
	}
	if keyCode != 0 {
		return keyCode
	}
	return ascii
}
