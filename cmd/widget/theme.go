package main

import (
	"github.com/xyproto/vt100"
)

// Structs and themes that can be used when drawing widgets

type Theme struct {
	Text, Background, Title,
	BoxLight, BoxDark, BoxBackground,
	ButtonFocus, ButtonText,
	ListFocus, ListText, ListBackground vt100.AttributeColor
	TL, TR, BL, BR, VL, VR, HT, HB rune
}

func NewTheme() *Theme {
	return &Theme{
		Text:           vt100.Black,
		Background:     vt100.BackgroundBlue,
		Title:          vt100.LightCyan,
		BoxLight:       vt100.White,
		BoxDark:        vt100.Black,
		BoxBackground:  vt100.BackgroundGray,
		ButtonFocus:    vt100.LightYellow,
		ButtonText:     vt100.White,
		ListFocus:      vt100.Red,
		ListText:       vt100.Black,
		ListBackground: vt100.BackgroundGray,
		TL:             '╭', // top left
		TR:             '╮', // top right
		BL:             '╰', // bottom left
		BR:             '╯', // bottom right
		VL:             '│', // vertical line, left side
		VR:             '│', // vertical line, right side
		HT:             '─', // horizontal line
		HB:             '─', // horizontal bottom line
	}
}

// Output text at the given coordinates, with the configured theme
func (t *Theme) Say(c *vt100.Canvas, x, y int, text string) {
	c.Write(uint(x), uint(y), t.Text, t.Background, text)
}

// Set the text color
func (t *Theme) SetTextColor(c vt100.AttributeColor) {
	t.Text = c
}

// Set the background color
func (t *Theme) SetBackgroundColor(c vt100.AttributeColor) {
	t.Background = c.Background()
}
