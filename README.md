# vt100

[![Build Status](https://travis-ci.org/xyproto/vt100.svg?branch=master)](https://travis-ci.org/xyproto/vt100) [![GoDoc](https://godoc.org/github.com/xyproto/vt100?status.svg)](https://godoc.org/github.com/xyproto/vt100) [![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/xyproto/vt100/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/xyproto/vt100)](https://goreportcard.com/report/github.com/xyproto/vt100)

### VT100 Terminal Package

* Supports colors and attributes.
* Developed for Linux. May work on other systems, but there are no guarantees.
* Can detect the terminal size.
* Can get key-presses, including arrow keys.
* Has a Canvas struct, for drawing to a buffer than only updating the characters that needs to be updated.
* Uses the spec directly, but memoizes the commands sent to the terminal, for speed.

### Simple use

Output "hi" in blue:

```go
fmt.Println(vt100.BrightColor("hi", "Blue"))
```

Erase the current line:

```go
vt100.Do("Erase Line")
```

Move the cursor 3 steps up (it's a bit verbose, but it's generated directly from spec, memoized for speed and is easy to wrap in a custom function):

```go
vt100.Set("Cursor Up", map[string]string{"{COUNT}": "3"})
```

The full overview of possible commands are at the top of `vt100.go`.

### Another example

See `cmd/move` for a more advanced example, where a character can be moved around with the arrow keys.

### Features and limitations

* Can detect letters, arrow keys and space. F12 and similar keys are not supported (they are supported by vt220).
* Resizing the terminal when using the Canvas struct may cause artifacts, for a brief moment.
* There is a bug with how characters are placed and/or read from the Canvas struct.

### General info

* Version: 0.2.0
* Licence: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
