package main

import (
	"fmt"
	"image/png"
	"math/rand"
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
	palette := []deepx.Color{
		deepx.MustColorFromHex("#000000"),
		deepx.MustColorFromHex("#cfcfcf"),
		deepx.MustColorFromHex("#14054c"),
		deepx.MustColorFromHex("#b61500"),
		deepx.MustColorFromHex("#ffd376"),
	}
	rand.Shuffle(len(palette), func(i, j int) {
		palette[i], palette[j] = palette[j], palette[i]
	})
	stereogramImg, err := deepx.NewStereogramFromMask(mask,
		deepx.WithColorPalette(palette...),
		deepx.WithOutputDPI(300),
		deepx.WithMaskTransparentColor(deepx.MustColorFromHex("#00718d")),
	)
	if err != nil {
		panic(fmt.Errorf("could not generate stereogram image: %v", err))
	}

	// INFO: save the image with gererated stereogram.
	out, err := os.OpenFile("result.png", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("could not open file to write stereogram: %v", err))
	}
	defer out.Close()
	if err := png.Encode(out, stereogramImg); err != nil {
		panic(fmt.Errorf("could not encode stereogram into png image: %v", err))
	}
}
