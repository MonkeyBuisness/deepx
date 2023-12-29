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

	// INFO: generate stereogram with default settings.
	stereogramImg, err := deepx.NewStereogramFromMask(mask)
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
