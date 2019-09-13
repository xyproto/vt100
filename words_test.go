package vt100

import (
	//"fmt"
	"testing"
)

func TestWords(t *testing.T) {
	s := "words can be colored, this is all gray"
	Words(s, "lightblue", "red", "lightgreen", "yellow", "darkgray")
	//fmt.Println(Words(s, "lightblue", "red", "lightgreen", "yellow", "darkgray"))
}
