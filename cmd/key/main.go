package main

import (
	"fmt"
	"github.com/xyproto/vt100"
	"time"
)

func main() {
	escCount := 0
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	tty.SetTimeout(10 * time.Millisecond)
	for {
		key := tty.Key()
		if key != 0 {
			fmt.Println(key)
		}
		if key == 27 {
			if escCount == 0 {
				fmt.Println("Press ESC again to exit")
			} else {
				fmt.Println("bye!")
			}
			escCount++
		}
		if escCount > 1 {
			break
		}
	}
	tty.Close()
}
