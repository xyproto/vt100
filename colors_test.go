package vt100

import (
	"fmt"
	"testing"
)

func TestBackground(t *testing.T) {
	if BackgroundBlue.String() == Blue.Background().String() {
		//fmt.Println("BLUE IS BLUE")
	} else {
		fmt.Println("BLUE BG IS NOT BLUE BG")
		fmt.Println(BackgroundBlue.String() + "FIRST" + Stop())
		fmt.Println(Blue.Background().String() + "SECOND" + Stop())
		t.Fail()
	}
}

func TestInts(t *testing.T) {
	ai := BackgroundBlue.Ints()
	bi := Blue.Background().Ints()
	if len(ai) != len(bi) {
		fmt.Println("A", ai)
		fmt.Println("B", bi)
		fmt.Println("length mismatch")
		t.Fail()
	}
	for i := 0; i < len(ai); i++ {
		if ai[i] != bi[i] {
			fmt.Println("NO")
			t.Fail()
		}
	}
}
