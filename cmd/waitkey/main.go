package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println("Waiting for either one of these: Return, Escape, Space or q ...")
	vt100.WaitForKey()
	fmt.Println("There you go!")
}
