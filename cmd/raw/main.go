package main

import (
	"github.com/xyproto/vt100"
)

func main() {
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	defer tty.Close()
	tty.PrintRawBytes()
}
