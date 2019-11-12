package main

import (
	"fmt"
	"unicode"

	"github.com/xyproto/vt100"
)

func main() {
	escCount := 0
	tty, err := vt100.NewTTY()
	if err != nil {
		panic(err)
	}
	for {
		key := tty.Rune()
		if key != rune(0) {
			if unicode.IsPrint(key) {
				fmt.Println(string(key))
			} else {
				fmt.Printf("%U\n", key)
			}
		}
		if key == rune(27) {
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
