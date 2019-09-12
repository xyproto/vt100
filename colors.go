package vt100

import (
	"fmt"
	"image/color"
	"strings"
)

// Color aliases, for ease of use, not for performance

type AttributeColor struct {
	attributes []string
}

var (
	// Dark foreground colors (+ light gray)
	Black     = NewAttributeColor("Black")
	Red       = NewAttributeColor("Red")
	Green     = NewAttributeColor("Green")
	Yellow    = NewAttributeColor("Yellow")
	Blue      = NewAttributeColor("Blue")
	Magenta   = NewAttributeColor("Magenta")
	Cyan      = NewAttributeColor("Cyan")
	LightGray = NewAttributeColor("White")

	// Light foreground colors (+ dark gray)
	DarkGray     = NewAttributeColor("Bright", "Black")
	LightRed     = NewAttributeColor("Bright", "Red")
	LightGreen   = NewAttributeColor("Bright", "Green")
	LightYellow  = NewAttributeColor("Bright", "Yellow")
	LightBlue    = NewAttributeColor("Bright", "Blue")
	LightMagenta = NewAttributeColor("Bright", "Magenta")
	LightCyan    = NewAttributeColor("Bright", "Cyan")
	White        = NewAttributeColor("Bright", "White")

	// Aliases
	Pink = LightMagenta
	Gray = LightGray

	// Dark background colors (+ light gray)
	BackgroundBlack     = NewAttributeColor("40")
	BackgroundRed       = NewAttributeColor("41")
	BackgroundGreen     = NewAttributeColor("42")
	BackgroundYellow    = NewAttributeColor("43")
	BackgroundBlue      = NewAttributeColor("44")
	BackgroundMagenta   = NewAttributeColor("45")
	BackgroundCyan      = NewAttributeColor("46")
	BackgroundLightGray = NewAttributeColor("47")

	// Light background colors (+ dark gray), not supported by all terminal emulators
	BackgroundDarkGray     = NewAttributeColor("Bright", "40")
	BackgroundLightRed     = NewAttributeColor("Bright", "41")
	BackgroundLightGreen   = NewAttributeColor("Bright", "42")
	BackgroundLightYellow  = NewAttributeColor("Bright", "43")
	BackgroundLightBlue    = NewAttributeColor("Bright", "44")
	BackgroundLightMagenta = NewAttributeColor("Bright", "45")
	BackgroundLightCyan    = NewAttributeColor("Bright", "46")
	BackgroundWhite        = NewAttributeColor("Bright", "47")

	// Aliases
	BackgroundPink = BackgroundLightMagenta
	BackgroundGray = BackgroundLightGray
)

func NewAttributeColor(attributes ...string) *AttributeColor {
	return &AttributeColor{attributes}
}

// Get the terminal codes for setting the attributes for colors, background colors, brightness etc
func (ac *AttributeColor) GetStart() string {
	attributeString := strings.Join(mapS(ac.attributes, AttributeNumber), ";")
	// Replace '{attr1};...;{attrn}' with the generated attribute string and return
	return get(specVT100, "Set Attribute Mode", map[string]string{"{attr1};...;{attrn}": attributeString}, false)
}

// Get the full string needed for outputting colored text + stopping the color attribute
func (ac *AttributeColor) Get(text string) string {
	return ac.GetStart() + text + NoColor()
}

// Use this color to output the given text
func (ac *AttributeColor) Output(text string) {
	fmt.Println(ac.Get(text))
}

func (ac *AttributeColor) Combine(other *AttributeColor) *AttributeColor {
	// Set an initial size of the map, where keys are attributes and values are bool
	amap := make(map[string]bool, len(ac.attributes)+len(other.attributes))
	for _, attr := range ac.attributes {
		amap[attr] = true
	}
	for _, attr := range other.attributes {
		amap[attr] = true
	}
	newAttributes := make([]string, len(amap))
	index := 0
	for attr, _ := range amap {
		newAttributes[index] = attr
		index++
	}
	return &AttributeColor{newAttributes}
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
