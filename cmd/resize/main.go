package main

import (
	"fmt"
	"github.com/xyproto/vt100"
	"time"
)

func main() {
	fmt.Println("Try resizing the terminal")
	for {
		w, h, err := vt100.TermSize()
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			fmt.Printf("%dx%d\n", w, h)
		}
		time.Sleep(time.Millisecond * 500)
	}
}
