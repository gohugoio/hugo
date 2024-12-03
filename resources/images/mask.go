package images

import (
	"fmt"
	"image"
	"image/draw"
	"image/color"

	"github.com/disintegration/gift"
)

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
		g := gift.New(gift.Resize(baseBounds.Dx(), baseBounds.Dy(), gift.LanczosResampling))
		resizedMask := image.NewRGBA(g.Bounds(maskImage.Bounds()))
		g.Draw(resizedMask, maskImage)
		maskImage = resizedMask
	}
	
	// Use gift to convert the resized mask to grayscale
	g := gift.New(gift.Grayscale())
	grayscaleMask := image.NewGray(g.Bounds(maskImage.Bounds()))
	g.Draw(grayscaleMask, maskImage)

	// Convert grayscale mask to alpha mask
	alphaMask := image.NewAlpha(baseBounds)
	for y := baseBounds.Min.Y; y < baseBounds.Max.Y; y++ {
		for x := baseBounds.Min.X; x < baseBounds.Max.X; x++ {
			grayValue := grayscaleMask.GrayAt(x, y).Y
			alphaMask.SetAlpha(x, y, color.Alpha{A: grayValue})
		}
	}

	// Create an RGBA output image
	outputImage := image.NewRGBA(baseBounds)

	// Apply the mask using draw.DrawMask
	draw.DrawMask(outputImage, baseBounds, baseImage, image.Point{}, alphaMask, image.Point{}, draw.Over)

	// Copy the result to the destination
	gift.New().Draw(dst, outputImage)
}

// Bounds returns the bounds of the resulting image.
func (f maskFilter) Bounds(imgBounds image.Rectangle) image.Rectangle {
    return image.Rect(0, 0, imgBounds.Dx(), imgBounds.Dy())
}
