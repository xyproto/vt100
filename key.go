package vt100

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/term"
)

var (
	defaultTimeout = 2 * time.Millisecond
	lastKey        int
)

// Lookup tables for faster key recognition

// Key codes for 3-byte sequences (arrows, Home, End)
var keyCodeLookup = map[[3]byte]int{
	{27, 91, 65}:  38, // Up Arrow
	{27, 91, 66}:  40, // Down Arrow
	{27, 91, 67}:  39, // Right Arrow
	{27, 91, 68}:  37, // Left Arrow
	{27, 91, 'H'}: 36, // Home
	{27, 91, 'F'}: 35, // End
}

// Key codes for 4-byte sequences (Page Up, Page Down, Home, End)
var pageNavLookup = map[[4]byte]int{
	{27, 91, 49, 126}: 36, // Home (ESC [1~)
	{27, 91, 52, 126}: 35, // End (ESC [4~)
	{27, 91, 53, 126}: 33, // Page Up
	{27, 91, 54, 126}: 34, // Page Down
}

// Key codes for 6-byte sequences (Ctrl-Insert)
var ctrlInsertLookup = map[[6]byte]int{
	{27, 91, 50, 59, 53, 126}: 258, // Ctrl-Insert (ESC [2;5~)
}

// String representations for 3-byte sequences
var keyStringLookup = map[[3]byte]string{
	{27, 91, 65}:  "↑", // Up Arrow
	{27, 91, 66}:  "↓", // Down Arrow
	{27, 91, 67}:  "→", // Right Arrow
	{27, 91, 68}:  "←", // Left Arrow
	{27, 91, 'H'}: "⇱", // Home
	{27, 91, 'F'}: "⇲", // End
}

// String representations for 4-byte sequences
var pageStringLookup = map[[4]byte]string{
	{27, 91, 49, 126}: "⇱", // Home
	{27, 91, 52, 126}: "⇲", // End
	{27, 91, 53, 126}: "⇞", // Page Up
	{27, 91, 54, 126}: "⇟", // Page Down
}

// String representations for 6-byte sequences (Ctrl-Insert)
var ctrlInsertStringLookup = map[[6]byte]string{
	{27, 91, 50, 59, 53, 126}: "⎘", // Ctrl-Insert (Copy)
}

type TTY struct {
	t       *term.Term
	timeout time.Duration
}

// NewTTY opens /dev/tty in raw and cbreak mode as a term.Term
func NewTTY() (*TTY, error) {
	t, err := term.Open("/dev/tty", term.RawMode, term.CBreakMode, term.ReadTimeout(defaultTimeout))
	if err != nil {
		return nil, err
	}
	return &TTY{t, defaultTimeout}, nil
}

// SetTimeout sets a timeout for reading a key
func (tty *TTY) SetTimeout(d time.Duration) {
	tty.timeout = d
	tty.t.SetReadTimeout(tty.timeout)
}

// Close will restore and close the raw terminal
func (tty *TTY) Close() {
	tty.t.Restore()
	tty.t.Close()
}

// asciiAndKeyCode processes input into an ASCII code or key code, handling multi-byte sequences like Ctrl-Insert
func asciiAndKeyCode(tty *TTY) (ascii, keyCode int, err error) {
	bytes := make([]byte, 6) // Use 6 bytes to cover longer sequences like Ctrl-Insert
	var numRead int

	// Set the terminal into raw mode and non-blocking mode with a timeout
	tty.RawMode()
	tty.NoBlock()
	tty.SetTimeout(tty.timeout)
	// Read bytes from the terminal
	numRead, err = tty.t.Read(bytes)
	// Restore the terminal settings
	tty.Restore()
	tty.t.Flush()

	if err != nil {
		return
	}

	switch numRead {
	case 1:
		ascii = int(bytes[0])
	case 3:
		// Check the lookup table for 3-byte sequences
		if code, found := keyCodeLookup[[3]byte{bytes[0], bytes[1], bytes[2]}]; found {
			keyCode = code
			return
		} else {
			// If not found, set ascii to the first byte
			ascii = int(bytes[0])
			return
		}
	case 4:
		// Check the lookup table for 4-byte sequences
		if code, found := pageNavLookup[[4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}]; found {
			keyCode = code
			return
		} else {
			ascii = int(bytes[0])
			return
		}
	case 6:
		// Check the lookup table for 6-byte sequences (Ctrl-Insert)
		if code, found := ctrlInsertLookup[[6]byte{bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5]}]; found {
			keyCode = code
			return
		} else {
			ascii = int(bytes[0])
			return
		}
	default:
		// For other cases, set ascii to the first byte
		ascii = int(bytes[0])
	}

	return
}

// Key reads the keycode or ASCII code and avoids repeated keys
func (tty *TTY) Key() int {
	ascii, keyCode, err := asciiAndKeyCode(tty)
	if err != nil {
		lastKey = 0
		return 0
	}
	var key int
	if keyCode != 0 {
		key = keyCode
	} else {
		key = ascii
	}
	if key == lastKey {
		lastKey = 0
		return 0
	}
	lastKey = key
	return key
}

// String reads a string, handling key sequences and printable characters
func (tty *TTY) String() string {
	bytes := make([]byte, 6)
	var numRead int

	// Set the terminal into raw mode with a timeout
	tty.RawMode()
	tty.SetTimeout(0)
	// Read bytes from the terminal
	numRead, err := tty.t.Read(bytes)
	// Restore the terminal settings
	tty.Restore()
	tty.t.Flush()

	if err != nil || numRead == 0 {
		return ""
	}

	switch numRead {
	case 1:
		r := rune(bytes[0])
		if unicode.IsPrint(r) {
			return string(r)
		}
		return "c:" + strconv.Itoa(int(r))
	case 3:
		if str, found := keyStringLookup[[3]byte{bytes[0], bytes[1], bytes[2]}]; found {
			return str
		} else {
			// If not found, attempt to interpret as UTF-8 string
			return string(bytes[:numRead])
		}
	case 4:
		if str, found := pageStringLookup[[4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}]; found {
			return str
		} else {
			return string(bytes[:numRead])
		}
	case 6:
		if str, found := ctrlInsertStringLookup[[6]byte{bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5]}]; found {
			return str
		} else {
			return string(bytes[:numRead])
		}
	default:
		// For other cases, attempt to interpret as UTF-8 string
		return string(bytes[:numRead])
	}
}

// Rune reads a rune, handling special sequences for arrows, Home, End, etc.
func (tty *TTY) Rune() rune {
	bytes := make([]byte, 6)
	var numRead int

	// Set the terminal into raw mode with a timeout
	tty.RawMode()
	tty.SetTimeout(0)
	// Read bytes from the terminal
	numRead, err := tty.t.Read(bytes)
	// Restore the terminal settings
	tty.Restore()
	tty.t.Flush()

	if err != nil || numRead == 0 {
		return rune(0)
	}

	switch numRead {
	case 1:
		return rune(bytes[0])
	case 3:
		if str, found := keyStringLookup[[3]byte{bytes[0], bytes[1], bytes[2]}]; found {
			return []rune(str)[0]
		} else {
			// If not found, attempt to interpret as UTF-8 rune
			r, _ := utf8.DecodeRune(bytes[:numRead])
			return r
		}
	case 4:
		if str, found := pageStringLookup[[4]byte{bytes[0], bytes[1], bytes[2], bytes[3]}]; found {
			return []rune(str)[0]
		} else {
			r, _ := utf8.DecodeRune(bytes[:numRead])
			return r
		}
	case 6:
		if str, found := ctrlInsertStringLookup[[6]byte{bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5]}]; found {
			return []rune(str)[0]
		} else {
			r, _ := utf8.DecodeRune(bytes[:numRead])
			return r
		}
	default:
		// For other cases, attempt to interpret as UTF-8 rune
		r, _ := utf8.DecodeRune(bytes[:numRead])
		return r
	}
}

// RawMode switches the terminal to raw mode
func (tty *TTY) RawMode() {
	term.RawMode(tty.t)
}

// NoBlock sets the terminal to cbreak mode (non-blocking)
func (tty *TTY) NoBlock() {
	tty.t.SetCbreak()
}

// Restore the terminal to its original state
func (tty *TTY) Restore() {
	tty.t.Restore()
}

// Flush flushes the terminal output
func (tty *TTY) Flush() {
	tty.t.Flush()
}

// WriteString writes a string to the terminal
func (tty *TTY) WriteString(s string) error {
	if n, err := tty.t.Write([]byte(s)); err != nil || n == 0 {
		return errors.New("no bytes written to the TTY")
	}
	return nil
}

// ReadString reads a string from the TTY
func (tty *TTY) ReadString() (string, error) {
	b, err := io.ReadAll(tty.t)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// PrintRawBytes for debugging raw byte sequences
func (tty *TTY) PrintRawBytes() {
	bytes := make([]byte, 6)
	var numRead int

	// Set the terminal into raw mode with a timeout
	tty.RawMode()
	tty.SetTimeout(0)
	// Read bytes from the terminal
	numRead, err := tty.t.Read(bytes)
	// Restore the terminal settings
	tty.Restore()
	tty.t.Flush()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Raw bytes: %v\n", bytes[:numRead])
}
