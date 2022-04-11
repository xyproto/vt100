package main

import (
	"fmt"
	"os"

	"github.com/xyproto/vt100"
)

func main() {
	tty, err := vt100.NewTTY()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer tty.Close()

	// Try to get the current background color
	r, g, b, err := vt100.GetBackgroundColor(tty)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("GOT: %.2f %.2f %.2f\n", r, g, b)
}
