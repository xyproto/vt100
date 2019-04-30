package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println("Cursor position:")
	x, y, err := vt100.CursorPos()
	fmt.Println(x, y, err)
}
