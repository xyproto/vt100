package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	fmt.Println(vt100.BrightColor("hi", "Green"))
	fmt.Println(vt100.BrightColor("done", "Blue"))

	fmt.Print(vt100.LightGreen.Get("process: "))
	vt100.LightRed.Output("ERROR")

	vt100.LightYellow.Output("jk")

	blue := vt100.BackgroundBlue.Get
	green := vt100.LightGreen.Get

	fmt.Printf("%s: %s\n", blue("status"), green("good"))
}
