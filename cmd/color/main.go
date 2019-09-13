package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

// This program demonstrates several different ways of outputting colored text
// Use `./color | cat -v` to see the color codes that are used.

func main() {
	vt100.Blue.Output("This is in blue")

	fmt.Println(vt100.BrightColor("hi", "Green"))
	fmt.Println(vt100.BrightColor("done", "Blue"))

	fmt.Println(vt100.Words("process: ERROR", "green", "red"))

	vt100.LightYellow.Output("jk")

	blue := vt100.BackgroundBlue.Get
	green := vt100.LightGreen.Get

	fmt.Printf("%s: %s\n", blue("status"), green("good"))

	combined := vt100.Blue.Background().Combine(vt100.Yellow).Combine(vt100.Reverse)
	combined.Output("DONE")
}
