package vt100

import (
	"fmt"
	"image/color"
	"strconv"
)

// Color aliases, for ease of use, not for performance

type AttributeColor struct {
	attribute string
	color     string
}

var (
	// Dark foreground colors (+ light gray)
	Black     = &AttributeColor{"", "Black"}
	Red       = &AttributeColor{"", "Red"}
	Green     = &AttributeColor{"", "Green"}
	Yellow    = &AttributeColor{"", "Yellow"}
	Blue      = &AttributeColor{"", "Blue"}
	Magenta   = &AttributeColor{"", "Magenta"}
	Cyan      = &AttributeColor{"", "Cyan"}
	LightGray = &AttributeColor{"", "White"}

	// Light foreground colors (+ dark gray)
	DarkGray     = &AttributeColor{"Bright", "Black"}
	LightRed     = &AttributeColor{"Bright", "Red"}
	LightGreen   = &AttributeColor{"Bright", "Green"}
	LightYellow  = &AttributeColor{"Bright", "Yellow"}
	LightBlue    = &AttributeColor{"Bright", "Blue"}
	LightMagenta = &AttributeColor{"Bright", "Magenta"}
	LightCyan    = &AttributeColor{"Bright", "Cyan"}
	White        = &AttributeColor{"Bright", "White"}

	// Aliases
	Pink = LightMagenta
	Gray = LightGray

	// Dark background colors (+ light gray)
	BackgroundBlack     = &AttributeColor{"", "40"}
	BackgroundRed       = &AttributeColor{"", "41"}
	BackgroundGreen     = &AttributeColor{"", "42"}
	BackgroundYellow    = &AttributeColor{"", "43"}
	BackgroundBlue      = &AttributeColor{"", "44"}
	BackgroundMagenta   = &AttributeColor{"", "45"}
	BackgroundCyan      = &AttributeColor{"", "46"}
	BackgroundLightGray = &AttributeColor{"", "47"}

	// Light background colors (+ dark gray), not supported by all terminal emulators
	BackgroundDarkGray     = &AttributeColor{"Bright", "40"}
	BackgroundLightRed     = &AttributeColor{"Bright", "41"}
	BackgroundLightGreen   = &AttributeColor{"Bright", "42"}
	BackgroundLightYellow  = &AttributeColor{"Bright", "43"}
	BackgroundLightBlue    = &AttributeColor{"Bright", "44"}
	BackgroundLightMagenta = &AttributeColor{"Bright", "45"}
	BackgroundLightCyan    = &AttributeColor{"Bright", "46"}
	BackgroundWhite        = &AttributeColor{"Bright", "47"}

	// Aliases
	BackgroundPink = BackgroundLightMagenta
	BackgroundGray = BackgroundLightGray
)

func (ac *AttributeColor) GetWithoutNoColor(text string) string {
	if ac.attribute == "Dark" || ac.attribute == "" {
		if num, err := strconv.Atoi(ac.color); err == nil {
			return ColorNum(num) + text
		}
		return AttributeOrColor(ac.color) + text

	}
	return AttributeAndColor(ac.attribute, ac.color) + text
}

func (ac *AttributeColor) Get(text string) string {
	return ac.GetWithoutNoColor(text) + NoColor()
}

func (ac *AttributeColor) Output(text string) {
	fmt.Println(ac.Get(text))
}

// Easteregg for displaying 24-bit true color on some terminals.
// This is not part of VT100.
// Example use:
// fmt.Println("not VT100, but " + vt100.TrueColor(color.RGBA{0xa0, 0xe0, 0xff, 0xff}, "TrueColor"))
func TrueColor(fg color.Color, text string) string {
	c := color.NRGBAModel.Convert(fg).(color.NRGBA)
	//fmt.Printf("(%d,%d,%d)\n", c.R, c.G, c.B)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.R, c.G, c.B, text)
}
