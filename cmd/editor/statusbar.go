package main

import (
	"github.com/xyproto/vt100"
	"strings"
	"time"
)

type StatusBar struct {
	msg    string               // status message
	fg     vt100.AttributeColor // draw foreground color
	bg     vt100.AttributeColor // draw background color
	editor *Editor              // an editor struct (for getting the colors when clearing the status)
	show   time.Duration        // show the message for how long before clearing
}

// Takes a foreground color, background color, foreground color for clearing,
// background color for clearing and a duration for how long to display status
// messages.
func NewStatusBar(fg, bg vt100.AttributeColor, editor *Editor, show time.Duration) *StatusBar {
	return &StatusBar{"", fg, bg, editor, show}
}

func (sb *StatusBar) Draw(c *vt100.Canvas) {
	w := int(c.W())
	c.Write(uint((w-len(sb.msg))/2), c.H()-1, sb.fg, sb.bg, sb.msg)
}

func (sb *StatusBar) SetMessage(msg string) {
	sb.msg = "       " + msg + "       "
}

func (sb *StatusBar) Clear(c *vt100.Canvas) {
	sb.msg = strings.Repeat(" ", len(sb.msg)+1)
	w := int(c.W())
	c.Write(uint((w-len(sb.msg))/2), c.H()-1, sb.editor.fg, sb.editor.bg, sb.msg)
	sb.msg = ""
}

// Draw a status message, then clear it after a configurable delay
func (sb *StatusBar) Show(c *vt100.Canvas) {
	if sb.msg == "" {
		return
	}
	sb.Draw(c)
	go func() {
		time.Sleep(sb.show)
		sb.Clear(c)
	}()
}
