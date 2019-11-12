package main

import (
	"fmt"

	"github.com/xyproto/vt100"
)

func main() {
	escCount := 0
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	for {
		key := tty.String()
		if key != "" {
			fmt.Println(key)
		}
		if key == "c:27" {
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
