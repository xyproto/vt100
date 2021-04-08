package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var (
		choices  = []string{"[0] Set up keyboard", "[1] Configure timezone", "[2] Set up disk layout", "[3] Install packages", "[4] Reboot"}
		selected = Menu("Installation Example", "yellow", choices, 20*time.Millisecond, "darkgray", "blue", "cyan", "red")
	)
	if selected < 0 {
		fmt.Println("No selection.")
		os.Exit(1)
	}
	// Output the selected item text
	fmt.Println(choices[selected])
}
