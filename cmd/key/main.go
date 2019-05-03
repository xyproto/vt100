package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	for {
		fmt.Println(vt100.Key())
	}
}
