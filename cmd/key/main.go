package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	escCount := 0
	r, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	for {
		key := r.Key()
		fmt.Println(key)
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
	r.Close()
}
