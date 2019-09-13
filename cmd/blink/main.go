package main

import (
	"fmt"
	"github.com/xyproto/vt100"
)

func main() {
	// Define four functions that takes text and returns text together with VT100 codes
	df := vt100.DarkGray.Start
	rf := vt100.Blink.Combine(vt100.Red).StartStop
	yf := vt100.Blink.Combine(vt100.Yellow).StartStop
	gf := vt100.Blink.Combine(vt100.Green).StartStop

	// Define a colored traffic light, with blinking lights
	trafficLight := df(`
	.-----.
	|  `) + rf("O") + df(`  |
	|     |
	|  `) + yf("O") + df(`  |
	|     |
	|  `) + gf("O") + df(`  |
	'-----'
	  | |
	  | |`+vt100.Stop())

	// Output the amazing artwork
	fmt.Printf("%s\n\n", trafficLight)
}
