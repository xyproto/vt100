package vt100

import (
	"image/color"
	"testing"
)

func TestToImage(t *testing.T) {
	// Character dimensions in pixels
	charWidth, charHeight := 8, 8

	// Initialize terminal settings
	Init()
	defer Close()

	// Create a new canvas
	canvas := NewCanvas()
	textWidth, textHeight := canvas.Size()

	// Calculate expected image dimensions
	expectedImgWidth, expectedImgHeight := int(textWidth)*charWidth, int(textHeight)*charHeight

	// Draw something on the canvas
	x, y := 10, 5 // Example coordinates
	canvas.chars[uint(y)*canvas.w+uint(x)].r = 'X'
	canvas.chars[uint(y)*canvas.w+uint(x)].fg = Red // Red foreground

	// Generate the image
	img, err := canvas.ToImage()
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	// Check the dimensions of the generated image
	if img.Bounds().Dx() != expectedImgWidth || img.Bounds().Dy() != expectedImgHeight {
		t.Errorf("Image dimensions are incorrect. Got %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), expectedImgWidth, expectedImgHeight)
	}

	// Check if the pixel at the specified character position has the expected color
	// Adjusting pixel position to account for character dimensions
	pixelX, pixelY := x*charWidth, y*charHeight
	expectedColor := ansiCodeToColor(Red, true) // Red
	actualColor := img.At(pixelX, pixelY)
	if !colorsAreEqual(expectedColor, actualColor) {
		t.Errorf("Pixel color at (%d, %d) is incorrect. Got %v, want %v", pixelX, pixelY, actualColor, expectedColor)
	}
}

func colorsAreEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
