package vt100

import (
	"fmt"
	"github.com/tudurom/ttyname"
	"os"
	"bytes"
	"io"
	//"io/ioutil"
)

// Return text with a bright color applied
func BrightColor(text, color string) string {
	return AttributeAndColor("Bright", color) + text + NoColor()
}

// Return x and y position of cursor
func CursorPos() (int, int, error) {
	name, err := ttyname.TTY()
	fmt.Println("TTY:", name)
	if err != nil {
		return 0, 0, err
	}
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return 0, 0, err
	}
	os.Stdout = w

    outC := make(chan string)
    // copy the output in a separate goroutine so printing can't block indefinitely
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()

	//f, err := os.OpenFile(name, os.O_RDONLY, 0644)
	//if err != nil {
	//	return 0, 0, err
	//}

	Do("Query Cursor Position")

	//b, err := ioutil.ReadAll(f)
	//if err != nil {
	//	return 0, 0, err
	//}
	//fmt.Printf("b %v %v\n", b, len(b))
	//f.Close()

	// back to normal state
    w.Close()
    os.Stdout = old // restoring the real stdout
    out := <-outC

	fmt.Println("--- close ---")
	fmt.Printf("READ %s %v %v\n", out, []byte(out), len([]byte(out)))
	os.Stdout = old
	//defer func() {
	//	f.Close()
	//}()
	// TODO: Return X, Y
	return 0, 0, nil
}
