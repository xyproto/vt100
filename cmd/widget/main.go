package main

import (
	"github.com/xyproto/vt100"
	"os"
	"os/signal"
	"syscall"
)

func draw() {
	c := vt100.NewCanvas()
	c.FillBackground(vt100.BackgroundBlue)

	box := NewBox()

	frame := box.GetFrame()
	frame.W = int(c.W())
	frame.H = int(c.H())
	box.SetFrame(frame)

	inner := box.GetInner()
	inner.X = 0
	inner.Y = 3 // This space is used by the title
	inner.W = frame.W - inner.X
	inner.H = frame.H - inner.Y
	box.SetInner(inner)

	infoBox := NewBox()
	infoBox.SetThirdSize(box)
	infoBox.FillWithPercentageMargins(box, 0.07, 0.1)

	t := NewTheme()
	infoBox.SetInner(t.DrawBox(c, infoBox, true))

	listBox := NewBox()
	choices := []string{"first", "second", "third"}
	listBox.SetInner(&Rect{0, 0, 6 + 2, len(choices)})
	listBox.Center(infoBox)
	t.DrawList(c, listBox, choices, 1)

	buttonBox1 := NewBox()
	buttonBox1.SetInner(&Rect{0, 0, 6 + 2, 1})
	buttonBox1.BottomCenterLeft(infoBox)
	t.DrawButton(c, buttonBox1, "OK", true)

	buttonBox2 := NewBox()
	buttonBox2.SetInner(&Rect{0, 0, 10 + 2, 1})
	buttonBox2.BottomCenterRight(infoBox)
	t.DrawButton(c, buttonBox2, "Cancel", false)

	c.Draw()
}

func main() {
	var (
		sigChan = make(chan os.Signal, 1)
	)

	signal.Notify(sigChan, syscall.SIGWINCH)
	go func() {
		for range sigChan {
			draw()
		}
	}()

	vt100.Init()
	defer vt100.Close()

	draw()

	vt100.WaitForKey()
}
