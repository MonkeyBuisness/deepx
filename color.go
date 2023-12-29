package deepx

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// Color represents color type.
type Color color.RGBA

// Hex returns color representation in hex format: #RRGGBBAA.
func (c Color) Hex() string {
	return fmt.Sprintf("#%.2X%.2X%.2X%.2X", c.R, c.G, c.B, c.A)
}

// RGBA returns color in native RGBA format.
func (c Color) RGBA() color.RGBA {
	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A,
	}
}

// Equal compares two color models.
func (c Color) Equal(other Color) bool {
	return c.A == other.A && c.B == other.B && c.R == other.R && c.G == other.G
}

// ColorFromHex converts hex color string (#RRGGBBAA) into color model.
func ColorFromHex(hex string) (*Color, error) {
	if hex = strings.TrimPrefix(hex, "#"); len(hex) < 8 {
		hex += strings.Repeat("F", 8-len(hex))
	}
	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return nil, err
	}
	return &Color{
		A: uint8(values & 0xFF),
		B: uint8((values >> 8) & 0xFF),
		G: uint8((values >> 16) & 0xFF),
		R: uint8(values >> 24),
	}, nil
}

// MustColorFromHex converts hex color string (#RRGGBBAA) into color model.
//
// It panics if the provided hes string is not a valid hexidecimal encoded color.
func MustColorFromHex(hex string) Color {
	c, err := ColorFromHex(hex)
	if err != nil {
		panic(err)
	}
	return *c
}

// ColorRGBA converts built-in color.Color model into Color representation.
func ColorRGBA(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{
		R: uint8(255 * float64(r) / math.MaxUint16),
		G: uint8(255 * float64(g) / math.MaxUint16),
		B: uint8(255 * float64(b) / math.MaxUint16),
		A: uint8(255 * float64(a) / math.MaxUint16),
	}
}
