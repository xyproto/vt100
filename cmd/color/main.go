package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println(vt100.BrightColor("hi", "Green"))
	fmt.Println(vt100.BrightColor("done", "Blue"))
}
