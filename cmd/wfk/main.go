package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println("Waiting for either one of these: Space, Return or Escape...")
	vt100.WaitForKey()
	fmt.Println("There you go!")
}
