package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println("Available VT100 commands:")
	for _, command := range vt100.Commands() {
		fmt.Println("\t" + command)
	}
	fmt.Println()
	fmt.Println("Available VT100 colors:")
	for _, color := range vt100.Colors() {
		fmt.Println("\t" + color)
	}
}
