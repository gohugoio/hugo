package images

import (
	"fmt"
	"image"
	"image/draw"
	"image/color"

	"github.com/disintegration/gift"
	"github.com/disintegration/imaging"
)

var _ gift.Filter = (*overlayFilter)(nil)

// maskFilter applies a mask image to a base image.
type maskFilter struct {
    mask ImageSource
}

// Draw applies the mask to the base image.
func (f maskFilter) Draw(dst draw.Image, baseImage image.Image, options *gift.Options) {
	maskImage, err := f.mask.DecodeImage()
	if err != nil {
		panic(fmt.Sprintf("failed to decode image: %s", err))
	}

	// Ensure the mask is the same size as the base image
    baseBounds := baseImage.Bounds()
    maskBounds := maskImage.Bounds()
	

	// Resize mask to match base image size if necessary
    if maskBounds.Dx() != baseBounds.Dx() || maskBounds.Dy() != baseBounds.Dy() {
        maskImage = imaging.Resize(maskImage, baseBounds.Dx(), baseBounds.Dy(), imaging.Lanczos)
    }

	alphaMask := image.NewAlpha(baseBounds)
	for y := baseBounds.Min.Y; y < baseBounds.Max.Y; y++ {
		for x := baseBounds.Min.X; x < baseBounds.Max.X; x++ {
			r, g, b, _ := maskImage.At(x, y).RGBA()
			brightness := (r + g + b) / 3 // Average RGB to get brightness
			alphaMask.SetAlpha(x, y, color.Alpha{A: uint8(brightness >> 8)})
		}
	}

	// Create an RGBA output image
	outputImage := image.NewRGBA(baseBounds)

	// Apply the mask using draw.DrawMask
	draw.DrawMask(outputImage, baseBounds, baseImage, image.Point{}, alphaMask, image.Point{}, draw.Over)

    // Copy the result to the destination
	//draw.Draw(dst, dst.Bounds(), outputImage, image.Point{}, draw.Src)
    gift.New().Draw(dst, outputImage)
}

// Bounds returns the bounds of the resulting image.
func (f maskFilter) Bounds(imgBounds image.Rectangle) image.Rectangle {
    return image.Rect(0, 0, imgBounds.Dx(), imgBounds.Dy())
}
