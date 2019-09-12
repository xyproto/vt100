package vt100

import (
	"fmt"
	"image/color"
	"strings"
)

// Color aliases, for ease of use, not for performance

type AttributeColor []string

var (
	// Non-color attributes
	ResetAll   = NewAttributeColor("Reset all attributes")
	Bright     = NewAttributeColor("Bright")
	Dim        = NewAttributeColor("Dim")
	Underscore = NewAttributeColor("Underscore")
	Blink      = NewAttributeColor("Blink")
	Reverse    = NewAttributeColor("Reverse")
	Hidden     = NewAttributeColor("Hidden")

	// There is also: reset, dim, underscore, reverse and hidden

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

	// Aliases
	BackgroundGray = BackgroundLightGray

	// Light background colors (+ dark gray), not used because "Bright" would apply to the foreground
	//BackgroundDarkGray     = NewAttributeColor("Bright", "40")
	//BackgroundLightRed     = NewAttributeColor("Bright", "41")
	//BackgroundLightGreen   = NewAttributeColor("Bright", "42")
	//BackgroundLightYellow  = NewAttributeColor("Bright", "43")
	//BackgroundLightBlue    = NewAttributeColor("Bright", "44")
	//BackgroundLightMagenta = NewAttributeColor("Bright", "45")
	//BackgroundLightCyan    = NewAttributeColor("Bright", "46")
	//BackgroundWhite        = NewAttributeColor("Bright", "47")

	DarkColorMap = map[string]AttributeColor{
		"black":   Black,
		"Black":   Black,
		"red":     Red,
		"Red":     Red,
		"green":   Green,
		"Green":   Green,
		"yellow":  Yellow,
		"Yellow":  Yellow,
		"blue":    Blue,
		"Blue":    Blue,
		"magenta": Magenta,
		"Magenta": Magenta,
		"cyan":    Cyan,
		"Cyan":    Cyan,
		"gray":    DarkGray,
		"Gray":    DarkGray,
		"white":   White,
		"White":   White,
	}

	LightColorMap = map[string]AttributeColor{
		"black":   DarkGray,
		"Black":   DarkGray,
		"red":     LightRed,
		"Red":     LightRed,
		"green":   LightGreen,
		"Green":   LightGreen,
		"yellow":  LightYellow,
		"Yellow":  LightYellow,
		"blue":    LightBlue,
		"Blue":    LightBlue,
		"magenta": LightMagenta,
		"Magenta": LightMagenta,
		"cyan":    LightCyan,
		"Cyan":    LightCyan,
		"gray":    LightGray,
		"Gray":    LightGray,
		"white":   White,
		"White":   White,
	}
)

func NewAttributeColor(attributes ...string) AttributeColor {
	return AttributeColor(attributes)
}

// For each element in a slice, apply the function f
func mapS(sl []string, f func(string) string) []string {
	result := make([]string, len(sl))
	for i, s := range sl {
		result[i] = f(s)
	}
	return result
}

// Get the terminal codes for setting the attributes for colors, background colors, brightness etc
func (ac AttributeColor) GetStart() string {
	attributeString := strings.Join(mapS(ac, AttributeNumber), ";")
	// Replace '{attr1};...;{attrn}' with the generated attribute string and return
	return get(specVT100, "Set Attribute Mode", map[string]string{"{attr1};...;{attrn}": attributeString}, false)
}

// Get the full string needed for outputting colored text + stopping the color attribute
func (ac AttributeColor) Get(text string) string {
	return ac.GetStart() + text + NoColor()
}

// Use this color to output the given text
func (ac AttributeColor) Output(text string) {
	fmt.Println(ac.Get(text))
}

func (ac AttributeColor) Combine(other AttributeColor) AttributeColor {
	// Set an initial size of the map, where keys are attributes and values are bool
	amap := make(map[string]bool, len(ac)+len(other))
	for _, attr := range ac {
		amap[attr] = true
	}
	for _, attr := range other {
		amap[attr] = true
	}
	newAttributes := make([]string, len(amap))
	index := 0
	for attr, _ := range amap {
		newAttributes[index] = attr
		index++
	}
	return AttributeColor(newAttributes)
}

// Return a new AttributeColor that has "Bright" added to the list of attributes
func (ac AttributeColor) Bright() AttributeColor {
	//lenAttr := len(ac.attributes)
	//newAttributes := make([]string, lenAttr + 1)
	//newAttributes[lenAttr] = "Bright"
	//return &AttributeColor{newAttributes}
	return AttributeColor(append(ac, "Bright"))
}

// Output a string at x, y with the given colors
func Write(x, y int, text string, fg, bg AttributeColor) {
	SetXY(uint(x), uint(y))
	fmt.Print(fg.Combine(bg).Get(text))
}

// Output a rune at x, y with the given colors
func WriteRune(x, y int, r rune, fg, bg AttributeColor) {
	Write(x, y, string(r), fg, bg)
}

// Easteregg for displaying 24-bit true color on some terminals.
// This is not part of the VT100 spec.
// Example use:
// fmt.Println("not VT100, but " + vt100.TrueColor(color.RGBA{0xa0, 0xe0, 0xff, 0xff}, "TrueColor"))
func TrueColor(fg color.Color, text string) string {
	c := color.NRGBAModel.Convert(fg).(color.NRGBA)
	//fmt.Printf("(%d,%d,%d)\n", c.R, c.G, c.B)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", c.R, c.G, c.B, text)
}
