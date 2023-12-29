package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/MonkeyBuisness/deepx"
)

func main() {
	// INFO: load mask image.
	mask, err := os.OpenFile("mask.png", os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("could not load mask: %v", err))
	}
	defer mask.Close()

	// INFO: generate stereogram.
	stereogramImg, err := deepx.NewStereogramFromMask(mask,
		// INFO: set a custom palette colors.
		deepx.WithColorPalette(
			deepx.MustColorFromHex("#00000000"),
			deepx.MustColorFromHex("#11b7d4"),
		),

		// INFO: increase the default DPI to avoid stereogram artifacts.
		deepx.WithOutputDPI(196),

		// INFO: the mask has a white color as a background, select it as transparent.
		deepx.WithMaskTransparentColor(deepx.MustColorFromHex("#ffffff")),
	)
	if err != nil {
		panic(fmt.Errorf("could not generate stereogram image: %v", err))
	}

	// INFO: save the image with gererated stereogram.
	out, err := os.OpenFile("deepx.png", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("could not open file to write stereogram: %v", err))
	}
	defer out.Close()
	if err := png.Encode(out, stereogramImg); err != nil {
		panic(fmt.Errorf("could not encode stereogram into png image: %v", err))
	}
}
