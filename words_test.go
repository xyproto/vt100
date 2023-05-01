package vt100

import (
	"testing"
)

func TestWords(_ *testing.T) {
	s := "words can be colored, this is all gray"
	Words(s, "lightblue", "red", "lightgreen", "yellow", "darkgray")
}
