package vt100

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// From http://www.termsys.demon.co.uk/vtansi.htm. Not the full spec, but a good subset.
// Update: Not even a good subset. Some of the codes are wrong!
const specVT100 = `

ANSI/VT100 Terminal Control Escape Sequences

Many computer terminals and terminal emulators support colour and cursor control through a system of escape sequences. One such standard is commonly referred to as ANSI Colour. Several terminal specifications are based on the ANSI colour standard, including VT100.

The following is a partial listing of the VT100 control set.

<ESC> represents the ASCII "escape" character, 0x1B. Bracketed tags represent modifiable decimal parameters; eg. {ROW} would be replaced by a row number.

Device Status
The following codes are used for reporting terminal/display settings, and vary depending on the implementation:

Query Device Code	<ESC>[c

    Requests a Report Device Code response from the device.

Report Device Code	<ESC>[{code}0c

    Generated by the device in response to Query Device Code request.

Query Device Status	<ESC>[5n

    Requests a Report Device Status response from the device.

Report Device OK	<ESC>[0n

    Generated by the device in response to a Query Device Status request; indicates that device is functioning correctly.

Report Device Failure	<ESC>[3n

    Generated by the device in response to a Query Device Status request; indicates that device is functioning improperly.

Query Cursor Position	<ESC>[6n

    Requests a Report Cursor Position response from the device.

Report Cursor Position	<ESC>[{ROW};{COLUMN}R

    Generated by the device in response to a Query Cursor Position request; reports current cursor position.

Terminal Setup
The h and l codes are used for setting terminal/display mode, and vary depending on the implementation. Line Wrap is one of the few setup codes that tend to be used consistently:

Reset Device		<ESC>c

    Reset all terminal settings to default.

Enable Line Wrap	<ESC>[?7h

    Text wraps to next line if longer than the length of the display area.

Disable Line Wrap	<ESC>[?7l

    Disables line wrapping.

Fonts
Some terminals support multiple fonts: normal/bold, swiss/italic, etc. There are a variety of special codes for certain terminals; the following are fairly standard:

Font Set G0		<ESC>(

    Set default font.

Font Set G1		<ESC>)

    Set alternate font.

Cursor Control

Cursor Home 		<ESC>[{ROW};{COLUMN}H

    Sets the cursor position where subsequent text will begin. If no row/column parameters are provided (ie. <ESC>[H), the cursor will move to the home position, at the upper left of the screen.

Cursor Up		<ESC>[{COUNT}A

    Moves the cursor up by COUNT rows; the default count is 1.

Cursor Down		<ESC>[{COUNT}B

    Moves the cursor down by COUNT rows; the default count is 1.

Cursor Forward		<ESC>[{COUNT}C

    Moves the cursor forward by COUNT columns; the default count is 1.

Cursor Backward		<ESC>[{COUNT}D

    Moves the cursor backward by COUNT columns; the default count is 1.

Force Cursor Position	<ESC>[{ROW};{COLUMN}f

    Identical to Cursor Home.

Save Cursor		<ESC>[s

    Save current cursor position.

Unsave Cursor		<ESC>[u

    Restores cursor position after a Save Cursor.

Save Cursor & Attrs	<ESC>7

    Save current cursor position.

Restore Cursor & Attrs	<ESC>8

    Restores cursor position after a Save Cursor.

Scrolling

Scroll Screen		<ESC>[r

    Enable scrolling for entire display.

Scroll Screen		<ESC>[{start};{end}r

    Enable scrolling from row {start} to row {end}.

Scroll Down		<ESC>D

    Scroll display down one line.

Scroll Up		<ESC>M

    Scroll display up one line.

Tab Control

Set Tab 		<ESC>H

    Sets a tab at the current position.

Clear Tab 		<ESC>[g

    Clears tab at the current position.

Clear All Tabs 		<ESC>[3g

    Clears all tabs.

Erasing Text

Erase End of Line	<ESC>[K

    Erases from the current cursor position to the end of the current line.

Erase Start of Line	<ESC>[1K

    Erases from the current cursor position to the start of the current line.

Erase Line		<ESC>[2K

    Erases the entire current line.

Erase Down		<ESC>[J

    Erases the screen from the current line down to the bottom of the screen.

Erase Up		<ESC>[1J

    Erases the screen from the current line up to the top of the screen.

Erase Screen		<ESC>[2J

    Erases the screen with the background colour and moves the cursor to home.

Printing
Some terminals support local printing:

Print Screen		<ESC>[i

    Print the current screen.

Print Line		<ESC>[1i

    Print the current line.

Stop Print Log		<ESC>[4i

    Disable log.

Start Print Log		<ESC>[5i

    Start log; all received text is echoed to a printer.

Define Key

Set Key Definition	<ESC>[{key};"{string}"p

    Associates a string of text to a keyboard key. {key} indicates the key by its ASCII value in decimal.

Set Display Attributes

Set Attribute Mode	<ESC>[{attr1};...;{attrn}m

    Sets multiple display attribute settings. The following lists standard attributes:

    0	Reset all attributes
    1	Bright
    2	Dim
    4	Underscore
    5	Blink
    7	Reverse
    8	Hidden

    	Foreground Colours
    30	Black
    31	Red
    32	Green
    33	Yellow
    34	Blue
    35	Magenta
    36	Cyan
    37	White

    	Background Colours
    40	Black
    41	Red
    42	Green
    43	Yellow
    44	Blue
    45	Magenta
    46	Cyan
    47	White
`

// memoization
var (
	memo        = make(map[string]string)
	memoMut     = &sync.RWMutex{}
	colorLookup = make(map[string]string)
	colorMut    = &sync.RWMutex{}
)

func flatten(m map[string]string) string {
	// TODO: Sort the keys first
	s := ""
	for k, v := range m {
		s += k + v + ";"
	}
	return s
}

// Given a terminal specification, a command and a map to replace strings with,
// return the terminal codes. If dummy is true <ESC> is returned instead of \033.
func get(specVT100, command string, replacemap map[string]string) string {
	if command == "" {
		return ""
	}
	combined := "get:" + command + flatten(replacemap)
	memoMut.RLock()
	if val, ok := memo[combined]; ok {
		memoMut.RUnlock()
		return val
	}
	memoMut.RUnlock()
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, command) {
			termCommand := strings.TrimSpace(trimmed[len(command):])
			for k, v := range replacemap {
				termCommand = strings.Replace(termCommand, k, v, 1)
			}
			// Return the terminal command
			termCommand = strings.Replace(termCommand, "<ESC>", "\033", -1)
			memoMut.Lock()
			memo[combined] = termCommand
			memoMut.Unlock()
			return termCommand
		}
	}
	return ""
}

// Return the terminal command, given a map to replace the values mentioned in the spec.
func Get(command string, replacemap map[string]string) string {
	return get(specVT100, command, replacemap)
}

// Do the terminal command, given a map to replace the values mentioned in the spec.
func Set(command string, replacemap map[string]string) {
	fmt.Print(get(specVT100, command, replacemap))
}

// Do the given command, with no parameters
func Do(command string) {
	fmt.Print(get(specVT100, command, map[string]string{}))
}

// Get the terminal command for setting a given color number
func ColorNum(colorNum int) string {
	return get(specVT100, "Set Attribute Mode", map[string]string{"{attr1};...;{attrn}": strconv.Itoa(colorNum)})
}

// Execute the terminal command for setting a given color number
func SetColorNum(colorNum int) {
	fmt.Print(ColorNum(colorNum))
}

// Returns the number (as a string) for a given attribute name.
// Returns the given string if the attribute was not found in the spec.
func AttributeNumber(name string) string {
	colorMut.RLock()
	if val, ok := colorLookup[name]; ok {
		colorMut.RUnlock()
		return val
	}
	colorMut.RUnlock()
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, name) {
			colorName := strings.TrimSpace(trimmed[:len(trimmed)-len(name)])
			colorMut.Lock()
			colorLookup[name] = colorName
			colorMut.Unlock()
			return colorName
		}
	}
	return name
}

// Execute the terminal command for setting a given display attribute name, like "Bright" or "Blink"
func AttributeOrColor(name string) string {
	combined := "DA:" + name
	memoMut.RLock()
	if val, ok := memo[combined]; ok {
		memoMut.RUnlock()
		return val
	}
	memoMut.RUnlock()
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, name) {
			numString := strings.TrimSpace(trimmed[:len(trimmed)-len(name)])
			num, err := strconv.Atoi(numString)
			if err != nil {
				return ""
			}
			termCommand := ColorNum(num)
			memoMut.Lock()
			memo[combined] = termCommand
			memoMut.Unlock()
			return termCommand
		}
	}
	return ""
}

// Execute the terminal command for setting a given display attribute name, like "Bright" or "Blink"
func SetAttribute(name string) {
	fmt.Print(AttributeOrColor(name))
}

// Get the terminal command for setting no colors or other display attributes
func NoColor() string {
	return get(specVT100, "Set Attribute Mode", map[string]string{"{attr1};...;{attrn}": "0"})
}

// Execute the terminal command for setting no colors or other display attributes
func SetNoColor() {
	fmt.Print(NoColor())
}

// Get the terminal command for setting a terminal attribute and a color
func AttributeAndColor(attr, name string) string {
	combined := "AAC:" + attr + name
	memoMut.RLock()
	if val, ok := memo[combined]; ok {
		memoMut.RUnlock()
		return val
	}
	memoMut.RUnlock()
	attribute := ""
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, attr) {
			attribute = strings.TrimSpace(trimmed[:len(trimmed)-len(attr)]) + ";"
			break
		}
	}
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, name) {
			numString := strings.TrimSpace(trimmed[:len(trimmed)-len(name)])
			termCommand := get(specVT100, "Set Attribute Mode", map[string]string{"{attr1};...;{attrn}": attribute + numString})
			memoMut.Lock()
			memo[combined] = termCommand
			memoMut.Unlock()
			return termCommand
		}
	}
	return ""
}

// Execute the terminal command for setting a terminal attribute and a color
func SetAttributeAndColor(attr, name string) {
	fmt.Print(AttributeAndColor(attr, name))
}

// Return all available commands
func Commands() []string {
	var commands []string
	for _, line := range strings.Split(specVT100, "\n") {
		if strings.Contains(line, "<ESC>") {
			elements := strings.SplitN(line, "<ESC>", 2)
			if len(elements) > 0 {
				command := strings.TrimSpace(elements[0])
				if command == "" || strings.Contains(command, ".") {
					continue
				}
				commands = append(commands, command)
			}
		}
	}
	return commands
}

// Check if a given string slice has the given string
func has(sl []string, s string) bool {
	for _, e := range sl {
		if s == e {
			return true
		}
	}
	return false
}

// Return all available colors
func Colors() []string {
	var colors []string
	colormode := false
	for _, line := range strings.Split(specVT100, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "Colours") || strings.Contains(trimmed, "Colors") {
			colormode = true
			continue
		} else if trimmed == "" {
			colormode = false
			continue
		}
		if colormode {
			fields := strings.Fields(trimmed)
			if len(fields) > 1 {
				color := strings.TrimSpace(fields[1])
				if !has(colors, color) {
					colors = append(colors, fields[1])
				}
			}
		}
	}
	return colors
}

// Return text with a bright color applied
func BrightColor(text, color string) string {
	return AttributeAndColor("Bright", color) + text + NoColor()
}

// Return text with a dark color applied
func DarkColor(text, color string) string {
	return AttributeOrColor(color) + text + NoColor()
}

func Init() {
	Reset()
	Clear()
	ShowCursor(false)
	SetLineWrap(false)
	EchoOff()
}

func Close() {
	SetLineWrap(true)
	ShowCursor(true)
	Home()
}

func EchoOff() {
	fmt.Print("\033[12h")
}
