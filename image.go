package vt100

import (
	"github.com/xyproto/burnfont"
	"image"
	"image/color"
	"image/draw"
)

func (c *Canvas) ToImage() (image.Image, error) {
	charWidth, charHeight := 8, 14
	width, height := int(c.w)*charWidth, int(c.h)*charHeight

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for y := uint(0); y < c.h; y++ {
		for x := uint(0); x < c.w; x++ {
			cr := c.chars[y*c.w+x]
			fgColor := ansiCodeToColor(cr.fg)
			bgColor := ansiCodeToColor(cr.bg)

			charRect := image.Rect(int(x)*charWidth, int(y)*charHeight, (int(x)+1)*charWidth, (int(y)+1)*charHeight)
			draw.Draw(img, charRect, &image.Uniform{bgColor}, image.Point{}, draw.Src)

			if cr.r != rune(0) {
				burnfont.DrawString(img, int(x)*charWidth, int(y)*charHeight, string(cr.r), fgColor)
			}
		}
	}

	return img, nil
}

func ansiCodeToColor(ac AttributeColor) color.NRGBA {
	if len(ac) == 0 {
		return color.NRGBA{0, 0, 0, 255} // Default color
	}

	code := ac[0] // Assuming the first byte is the ANSI color code

	switch {
	case code >= 30 && code <= 37:
		// Standard ANSI foreground colors
		return standardANSIColors[code-30]
	case code >= 40 && code <= 47:
		// Standard ANSI background colors
		return standardANSIColors[code-40]
	// Add more cases for extended color codes (0-255) if needed
	default:
		return color.NRGBA{0, 0, 0, 255} // Default to black
	}
}

var standardANSIColors = []color.NRGBA{
	color.NRGBA{0, 0, 0, 255},       // Black
	color.NRGBA{255, 0, 0, 255},     // Red
	color.NRGBA{0, 255, 0, 255},     // Green
	color.NRGBA{255, 255, 0, 255},   // Yellow
	color.NRGBA{0, 0, 255, 255},     // Blue
	color.NRGBA{255, 0, 255, 255},   // Magenta
	color.NRGBA{0, 255, 255, 255},   // Cyan
	color.NRGBA{255, 255, 255, 255}, // White
}
