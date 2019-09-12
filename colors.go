package vt100

import "fmt"

// Color aliases, for ease of use

type AttributeColor struct {
	attribute string
	color     string
}

var (
	Black     = &AttributeColor{"Dark", "Black"}
	Red       = &AttributeColor{"Dark", "Red"}
	Green     = &AttributeColor{"Dark", "Green"}
	Yellow    = &AttributeColor{"Dark", "Yellow"}
	Blue      = &AttributeColor{"Dark", "Blue"}
	Magenta   = &AttributeColor{"Dark", "Magenta"}
	Cyan      = &AttributeColor{"Dark", "Cyan"}
	LightGray = &AttributeColor{"Dark", "White"}

	DarkGray     = &AttributeColor{"Bright", "Black"}
	LightRed     = &AttributeColor{"Bright", "Red"}
	LightGreen   = &AttributeColor{"Bright", "Green"}
	LightYellow  = &AttributeColor{"Bright", "Yellow"}
	LightBlue    = &AttributeColor{"Bright", "Blue"}
	LightMagenta = &AttributeColor{"Bright", "Magenta"}
	LightCyan    = &AttributeColor{"Bright", "Cyan"}
	White        = &AttributeColor{"Bright", "White"}

	Pink = LightMagenta
	Gray = DarkGray
)

func (ac *AttributeColor) Get(text string) string {
	if ac.attribute == "Dark" || ac.attribute == "" {
		return AttributeOrColor(ac.color) + text + NoColor()
	}
	return AttributeAndColor(ac.attribute, ac.color) + text + NoColor()
}

func (ac *AttributeColor) Output(text string) {
	fmt.Println(ac.Get(text))
}
