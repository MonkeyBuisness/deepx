# DeepX

![logo](./examples/logo/deepx.png)

Add a new dimension to your pixels!

## About

As a child, I really loved looking at stereograms from various magazines and newspapers.
I eagerly searched for them in every publication. As an adult, I realized that now I can create them myself.
I hope you too are amazed at how easily the human eye can be fooled.

The implementation of the algorithm is based on a [scientific publication](https://www2.cs.sfu.ca/CourseCentral/414/li/material/refs/SIRDS-Computer-94.pdf) of `Harold W. Thimbleby` (University of Stirling), `Stuart Inglis` and `Ian H. Witten` (University of Waikato).

## Quick Start

Make sure you have the latest version of the [golang compiler](https://go.dev/doc/install) installed.

> Currently, the creation of stereograms is supported only based on other images (i.e. masks)

Prepare your **mask image**. For example:

![mask](./examples/gopher/mask.png)

The mask image will be interpreted as monochrome, regardless of the actual number of colors encoded in that image.
All pixels in the mask image that have a zero alpha channel (transparent) will be ignored (by default), and the remaining pixels will be included in the final mask image.

> You can specify your own color, which should be considered transparent for each specific mask.

Create a new `.go` file and import **deepX** package. Write a piece of code that will load the mask image and convert it into a stereogram:

```go
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
			deepx.MustColorFromHex("#6ad6e3"),
			deepx.MustColorFromHex("#000000"),
		),

		// INFO: increase the default DPI to avoid stereogram artifacts.
		deepx.WithOutputDPI(144),

		// INFO: the mask has a #6ad6e3 color as a background, select it as transparent.
		deepx.WithMaskTransparentColor(deepx.MustColorFromHex("#ffffff")),
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

```

> The library code is documented, so refer to it for additional functionality.

Run your program:

```console
$ go run main.go
```

Result:

![result.png](./examples/gopher/result.png)

You can find more examples [here](./examples/).

## General Tips

To achieve better generation result:

- Try not to use very large or very small images as masks
- If possible, convert the image to `PNG` format to explicitly indicate transparency
- Avoid images with a lot of detail or color
- If, after generating a stereogram, there are artifacts in the image, experiment with the values for the `DPI` parameter. Make it larger or smaller and see what results it produces
- Use contrasting colors (such as white and black) for the final palette, and don't use colors with low intensity
