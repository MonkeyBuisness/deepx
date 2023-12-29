package deepx

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"math/rand"
	"sync"
)

var (
	defaultStereogramCfg = StereogramConfig{
		Palette: make([]Color, 0),
		Mu:      1 / 3.,
		DPI:     72,
		ERatio:  2.5,
	}
)

// StereogramConfig represents a stereogram image processing configuration model.
type StereogramConfig struct {

	// Contains the color of the mask image pixels that should be considered transparent.
	//
	// If not specifed, every pixel in a mask image that has a zero alpha channel
	// will be considered transparent.
	MaskTransparentColor *Color

	// Represents a list of colors used to create pixels in a stereogram image.
	//
	// It's recommended to specify at least 2 different colors so that
	// the final stereogram image is colored (if you specify one color, the final image
	// will always be monotonously filled with this color, which doesn't make sense 0_o).
	// Each color presented will be randomly selected to form a unique pixel
	// in the stereogram image.
	//
	// If the list of colors is not specified (by defaul), then a randomization algorithm
	// will be applied when forming each individual pixel color.
	Palette []Color

	// Depth of field (fraction of viewing distance).
	//
	// Equal to 1/3 by default.
	Mu float64

	// Output stereogram image DPI.
	//
	// By defualt has 72 pixels per inch.
	DPI int

	// Eye separation ratio.
	//
	// Eye separation is assumed to be 2.5 * DPI in by default.
	ERatio float64
}

// StereogramOption represents type for stereogram image processing option.
type StereogramOption func(*StereogramConfig)

// NewStereogramFromMask creates a new "Random-Dot Stereogram" image from the provided
// mask source using the algorithm of Harold W. Thimbleby, Stuart Inglis and Ian H:
// https://www2.cs.sfu.ca/CourseCentral/414/li/material/refs/SIRDS-Computer-94.pdf
//
// The mask source must contain an encoded valid png, jpeg or gif image data.
// The mask image will be interpreted as monochrome, regardless of the actual number of colors
// encoded in that image.
// All pixels in the mask image that have a zero alpha channel (transparent)
// will be ignored (by default), and the remaining pixels will be included in the final mask image.
// To explicitly specify the color that should be perceived as transparent in the mask image,
// specify a `WithMaskTransparentColor(...)` in the list of options.
//
// A list of options can be provided to specify additional stereogram processing settings.
func NewStereogramFromMask(maskSrc io.Reader, opts ...StereogramOption) (*image.RGBA, error) {
	maskImg, _, err := image.Decode(maskSrc)
	if err != nil {
		return nil, fmt.Errorf("could not decode mask image data: %v", err)
	}
	cfg := defaultStereogramCfg
	for _, opt := range opts {
		opt(&cfg)
	}
	e := math.Ceil(cfg.ERatio * float64(cfg.DPI))
	maskImgBounds := maskImg.Bounds()
	imgWidth, imgHeight := maskImgBounds.Dx(), maskImgBounds.Dy()
	stereogramImg := drawAutoStereogram(
		newDepthBufferFromImage(maskImg, cfg.MaskTransparentColor),
		imgWidth, imgHeight, cfg.Mu, e, cfg.Palette,
	)
	return stereogramImg, nil
}

// WithMaskTransparentColor sets the color that must be transparent for the mask source image.
//
// By default, every pixel in a mask image that has a zero alpha channel
// will be considered transparent.
func WithMaskTransparentColor(color Color) StereogramOption {
	return func(cfg *StereogramConfig) {
		cfg.MaskTransparentColor = &color
	}
}

// WithColorPalette sets the list of colors used to create pixels in a stereogram image.
//
// By default, this list is empty and therefore the colors of each pixel will be selected randomly.
func WithColorPalette(palette ...Color) StereogramOption {
	return func(cfg *StereogramConfig) {
		cfg.Palette = palette
	}
}

// WithOutputDPI sets the DPI of the stegeogram image.
//
// 72 by default.
func WithOutputDPI(dpi int) StereogramOption {
	return func(cfg *StereogramConfig) {
		cfg.DPI = dpi
	}
}

// WithEyeSepartionRatio sets the eye separtion ratio.
//
// 2.5 by default.
func WithEyeSepartionRatio(ratio float64) StereogramOption {
	return func(cfg *StereogramConfig) {
		cfg.ERatio = ratio
	}
}

func projSeparation(z, mu, e float64) int {
	return int(math.Ceil((1 - mu*z) * e / (2 - mu*z)))
}

func getRandomPaletteColor(palette []Color) Color {
	if len(palette) == 0 {
		return Color{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}
	}
	return palette[rand.Intn(len(palette))]
}

func isTransparentMaskPixel(pxColor color.Color, maskTransparentColor *Color) bool {
	if maskTransparentColor == nil {
		_, _, _, a := pxColor.RGBA()
		return a == 0
	}
	return ColorRGBA(pxColor).Equal(*maskTransparentColor)
}

func drawAutoStereogram(
	zBuf [][]float64,
	imgWidth, imgHeight int,
	mu, e float64,
	palette []Color,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	var wg sync.WaitGroup
	wg.Add(imgHeight)
	for y := 0; y < imgHeight; y++ {
		go func(y int) {
			defer wg.Done()
			same := make([]int, imgWidth)
			for x := 0; x < imgWidth; x++ {
				same[x] = x
			}
			for x := 0; x < imgWidth; x++ {
				s := projSeparation(zBuf[x][y], mu, e)
				left := x - (s+(s&y&1))/2
				right := left + s
				if left < 0 || right >= imgWidth {
					continue
				}
				var isVisible bool
				for t := 1; ; t++ {
					zt := zBuf[x][y] + 2*(2-mu*zBuf[x][y])*float64(t)/(mu*e)
					isVisible = zBuf[x-t][y] < zt && zBuf[x+t][y] < zt
					if !(isVisible && zt < 1) {
						break
					}
				}
				if !isVisible {
					continue
				}
				for k := same[left]; k != left && k != right; k = same[left] {
					if k < right {
						left = k
						continue
					}
					left, right = right, k
				}
				same[left] = right
			}
			pixels := make([]Color, imgWidth)
			for x := imgWidth - 1; x >= 0; x-- {
				pixels[x] = pixels[same[x]]
				if same[x] == x {
					pixels[x] = getRandomPaletteColor(palette)
				}
				img.Set(x, y, pixels[x].RGBA())
			}
		}(y)
	}
	wg.Wait()
	return img
}

func newDepthBufferFromImage(img image.Image, transparentColor *Color) [][]float64 {
	imgBounds := img.Bounds()
	sizeX, sizeY := imgBounds.Dx(), imgBounds.Dy()
	z := make([][]float64, sizeX)
	for x := 0; x < sizeX; x++ {
		z[x] = make([]float64, sizeY)
		for y := 0; y < sizeY; y++ {
			if isTransparentMaskPixel(img.At(x, y), transparentColor) {
				continue
			}
			z[x][y] = 1
		}
	}
	return z
}
