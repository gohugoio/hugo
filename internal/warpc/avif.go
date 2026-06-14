// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package warpc

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/gohugoio/hugo/common/hugio"
)

var (
	_ SourceProvider      = AvifInput{}
	_ DestinationProvider = AvifInput{}
)

type AvifInput struct {
	Source      hugio.SizeReader `json:"-"`       // Will be sent in a separate stream.
	Destination io.Writer        `json:"-"`       // Will be used to write the result to.
	Options     map[string]any   `json:"options"` // Config options.
	Params      map[string]any   `json:"params"`  // Command params (width, height, etc.).
}

func (a AvifInput) GetSource() hugio.SizeReader {
	return a.Source
}

func (a AvifInput) GetDestination() io.Writer {
	return a.Destination
}

type AvifOutput struct {
	Params CommonImageProcessingParams `json:"params"`
}

type AvifCodec struct {
	d func() (Dispatcher[AvifInput, AvifOutput], error)
}

func (d *AvifCodec) DecodeConfig(r io.Reader) (image.Config, error) {
	dd, err := d.d()
	if err != nil {
		return image.Config{}, err
	}

	rr, err := hugio.ToSizeReader(r)
	if err != nil {
		return image.Config{}, err
	}
	message := Message[AvifInput]{
		Header: Header{
			Version:       1,
			Command:       "config",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON},
		},

		Data: AvifInput{
			Source: rr,
		},
	}

	out, err := dd.Execute(context.Background(), message)
	if err != nil {
		return image.Config{}, err
	}
	return image.Config{
		Width:      out.Data.Params.Width,
		Height:     out.Data.Params.Height,
		ColorModel: color.RGBAModel,
	}, nil
}

func (d *AvifCodec) Decode(r io.Reader) (image.Image, error) {
	dd, err := d.d()
	if err != nil {
		return nil, err
	}

	source, err := hugio.ToSizeReader(r)
	if err != nil {
		return nil, err
	}

	var destination bytes.Buffer

	message := Message[AvifInput]{
		Header: Header{
			Version:       1,
			Command:       "decode",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: AvifInput{
			Source:      source,
			Destination: &destination,
			Options:     map[string]any{},
		},
	}

	out, err := dd.Execute(context.Background(), message)
	if err != nil {
		return nil, err
	}

	w, h, stride, depth := out.Data.Params.Width, out.Data.Params.Height, out.Data.Params.Stride, out.Data.Params.Depth
	if w == 0 || h == 0 || stride == 0 {
		return nil, fmt.Errorf("received invalid image dimensions: %dx%d stride %d", w, h, stride)
	}

	// For 10+ bit HDR images, the C code returns 16-bit RGBA data.
	isHDR := depth > 8

	// libavif returns 16-bit data in native (little-endian) byte order,
	// but Go's NRGBA64 expects big-endian. Swap bytes for HDR images.
	if isHDR {
		pix := destination.Bytes()
		for i := 0; i < len(pix); i += 2 {
			pix[i], pix[i+1] = pix[i+1], pix[i]
		}
	}

	if len(out.Data.Params.FrameDurations) > 0 {
		img := &AnimatedImage{
			frameDurations:          out.Data.Params.FrameDurations,
			loopCount:               avifLoopCountToGo(out.Data.Params.LoopCount),
			depth:                   depth,
			colorPrimaries:          out.Data.Params.ColorPrimaries,
			transferCharacteristics: out.Data.Params.TransferCharacteristics,
			matrixCoefficients:      out.Data.Params.MatrixCoefficients,
			maxCLL:                  out.Data.Params.MaxCLL,
			maxPALL:                 out.Data.Params.MaxPALL,
		}

		frameSize := stride * h
		pixLen := len(destination.Bytes())
		if frameSize == 0 || pixLen%frameSize != 0 {
			return nil, fmt.Errorf("decoded AVIF buffer size %d is not a multiple of frame size %d", pixLen, frameSize)
		}
		frameCount := pixLen / frameSize
		if frameCount != len(out.Data.Params.FrameDurations) {
			return nil, fmt.Errorf("decoded AVIF frame count %d does not match frame durations %d", frameCount, len(out.Data.Params.FrameDurations))
		}
		frames := make([]image.Image, frameCount)
		for i := range frames {
			frameBytes := destination.Bytes()[i*frameSize : (i+1)*frameSize]
			if isHDR {
				// NRGBA64 for HDR - libavif returns non-premultiplied alpha.
				frameImg := &image.NRGBA64{
					Pix:    frameBytes,
					Stride: stride,
					Rect:   image.Rect(0, 0, w, h),
				}
				frames[i] = frameImg
			} else {
				// NRGBA - libavif returns non-premultiplied alpha.
				frameImg := &image.NRGBA{
					Pix:    frameBytes,
					Stride: stride,
					Rect:   image.Rect(0, 0, w, h),
				}
				frames[i] = frameImg
			}
		}
		img.SetFrames(frames)
		return img, nil
	}

	if isHDR {
		// NRGBA64 for HDR - libavif returns non-premultiplied alpha.
		// Wrap in AnimatedImage to preserve color properties through the pipeline.
		baseImg := &image.NRGBA64{
			Pix:    destination.Bytes(),
			Stride: stride,
			Rect:   image.Rect(0, 0, w, h),
		}
		img := &AnimatedImage{
			Image:                   baseImg,
			frames:                  []image.Image{baseImg},
			depth:                   depth,
			colorPrimaries:          out.Data.Params.ColorPrimaries,
			transferCharacteristics: out.Data.Params.TransferCharacteristics,
			matrixCoefficients:      out.Data.Params.MatrixCoefficients,
			maxCLL:                  out.Data.Params.MaxCLL,
			maxPALL:                 out.Data.Params.MaxPALL,
		}
		return img, nil
	}

	// NRGBA - libavif returns non-premultiplied alpha.
	img := &image.NRGBA{
		Pix:    destination.Bytes(),
		Stride: stride,
		Rect:   image.Rect(0, 0, w, h),
	}

	return img, nil
}

// avifLoopCountToGo translates libavif's repetitionCount (-1 = infinite,
// 0 = play once, N = N repetitions) to the convention used by image/gif and
// the WebP codec (0 = infinite, -1 = play once, N = N repetitions).
func avifLoopCountToGo(repetitionCount int) int {
	switch repetitionCount {
	case -1, -2: // AVIF_REPETITION_COUNT_INFINITE, AVIF_REPETITION_COUNT_UNKNOWN
		return 0
	case 0:
		return -1
	default:
		return repetitionCount
	}
}

func (d *AvifCodec) Encode(w io.Writer, src image.Image, options map[string]any) error {
	dd, err := d.d()
	if err != nil {
		return err
	}

	var cmd string
	var source hugio.SizeReader
	var params map[string]any = make(map[string]any)

	switch img := src.(type) {
	case *AnimatedImage:
		// AnimatedImage wraps HDR images to preserve color properties.
		// For single-frame, extract the underlying image and add color properties.
		frames := img.GetFrames()
		if len(frames) == 0 {
			return fmt.Errorf("AnimatedImage has no frames")
		}
		// Get color properties from AnimatedImage.
		params["colorPrimaries"] = img.GetColorPrimaries()
		params["transferCharacteristics"] = img.GetTransferCharacteristics()
		params["matrixCoefficients"] = img.GetMatrixCoefficients()
		params["maxCLL"] = img.GetMaxCLL()
		params["maxPALL"] = img.GetMaxPALL()

		// Handle the first frame based on its type.
		firstFrame := frames[0]
		switch frame := firstFrame.(type) {
		case *image.NRGBA64:
			cmd = "encodeNRGBA"
			pix := make([]byte, len(frame.Pix))
			for i := 0; i < len(pix); i += 2 {
				pix[i], pix[i+1] = frame.Pix[i+1], frame.Pix[i]
			}
			var err error
			source, err = hugio.ToSizeReader(bytes.NewReader(pix))
			if err != nil {
				return err
			}
			params["width"] = frame.Rect.Dx()
			params["height"] = frame.Rect.Dy()
			params["stride"] = frame.Stride
			params["depth"] = img.GetDepth()
			if params["depth"] == 0 {
				params["depth"] = 10 // Default to 10-bit for HDR.
			}
		case *image.NRGBA:
			cmd = "encodeNRGBA"
			var err error
			source, err = hugio.ToSizeReader(bytes.NewReader(frame.Pix))
			if err != nil {
				return err
			}
			params["width"] = frame.Rect.Dx()
			params["height"] = frame.Rect.Dy()
			params["stride"] = frame.Stride
			params["depth"] = 8
		default:
			// Convert to NRGBA64 for HDR or NRGBA for SDR.
			b := firstFrame.Bounds()
			if img.GetDepth() > 8 {
				newImg := image.NewNRGBA64(image.Rect(0, 0, b.Dx(), b.Dy()))
				draw.Draw(newImg, newImg.Bounds(), firstFrame, b.Min, draw.Src)
				cmd = "encodeNRGBA"
				pix := make([]byte, len(newImg.Pix))
				for i := 0; i < len(pix); i += 2 {
					pix[i], pix[i+1] = newImg.Pix[i+1], newImg.Pix[i]
				}
				var err error
				source, err = hugio.ToSizeReader(bytes.NewReader(pix))
				if err != nil {
					return err
				}
				params["width"] = newImg.Rect.Dx()
				params["height"] = newImg.Rect.Dy()
				params["stride"] = newImg.Stride
				params["depth"] = img.GetDepth()
			} else {
				newImg := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
				draw.Draw(newImg, newImg.Bounds(), firstFrame, b.Min, draw.Src)
				cmd = "encodeNRGBA"
				var err error
				source, err = hugio.ToSizeReader(bytes.NewReader(newImg.Pix))
				if err != nil {
					return err
				}
				params["width"] = newImg.Rect.Dx()
				params["height"] = newImg.Rect.Dy()
				params["stride"] = newImg.Stride
				params["depth"] = 8
			}
		}
	case *image.NRGBA64:
		// 16-bit RGBA for HDR images.
		// Go's NRGBA64 is big-endian, but libavif expects little-endian.
		// Copy and swap bytes.
		cmd = "encodeNRGBA"
		pix := make([]byte, len(img.Pix))
		for i := 0; i < len(pix); i += 2 {
			pix[i], pix[i+1] = img.Pix[i+1], img.Pix[i]
		}
		var err error
		source, err = hugio.ToSizeReader(bytes.NewReader(pix))
		if err != nil {
			return err
		}
		params["width"] = img.Rect.Dx()
		params["height"] = img.Rect.Dy()
		params["stride"] = img.Stride
		params["depth"] = 10 // Encode HDR as 10-bit AVIF.
	case *image.RGBA64:
		// 16-bit RGBA for HDR images.
		// Go's RGBA64 is big-endian, but libavif expects little-endian.
		// Copy and swap bytes.
		cmd = "encodeNRGBA"
		pix := make([]byte, len(img.Pix))
		for i := 0; i < len(pix); i += 2 {
			pix[i], pix[i+1] = img.Pix[i+1], img.Pix[i]
		}
		var err error
		source, err = hugio.ToSizeReader(bytes.NewReader(pix))
		if err != nil {
			return err
		}
		params["width"] = img.Rect.Dx()
		params["height"] = img.Rect.Dy()
		params["stride"] = img.Stride
		params["depth"] = 10 // Encode HDR as 10-bit AVIF.
	case *image.NRGBA:
		cmd = "encodeNRGBA"
		var err error
		source, err = hugio.ToSizeReader(bytes.NewReader(img.Pix))
		if err != nil {
			return err
		}
		params["width"] = img.Rect.Dx()
		params["height"] = img.Rect.Dy()
		params["stride"] = img.Stride
		params["depth"] = 8
	case *image.Gray:
		cmd = "encodeGray"
		var err error
		source, err = hugio.ToSizeReader(bytes.NewReader(img.Pix))
		if err != nil {
			return err
		}
		params["width"] = img.Rect.Dx()
		params["height"] = img.Rect.Dy()
		params["stride"] = img.Stride
		params["depth"] = 8
	case *image.Gray16:
		// 16-bit grayscale for HDR.
		cmd = "encodeGray"
		var err error
		source, err = hugio.ToSizeReader(bytes.NewReader(img.Pix))
		if err != nil {
			return err
		}
		params["width"] = img.Rect.Dx()
		params["height"] = img.Rect.Dy()
		params["stride"] = img.Stride
		params["depth"] = 10
	default:
		// Check if source is HDR (RGBA64 color model) to preserve quality.
		if src.ColorModel() == color.RGBA64Model || src.ColorModel() == color.NRGBA64Model {
			b := src.Bounds()
			newImg := image.NewNRGBA64(image.Rect(0, 0, b.Dx(), b.Dy()))
			draw.Draw(newImg, newImg.Bounds(), src, b.Min, draw.Src)
			return d.Encode(w, newImg, options)
		}
		// Convert to NRGBA and try again.
		b := src.Bounds()
		newImg := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(newImg, newImg.Bounds(), src, b.Min, draw.Src)
		return d.Encode(w, newImg, options)
	}

	message := Message[AvifInput]{
		Header: Header{
			Version:       1,
			Command:       cmd,
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: AvifInput{
			Source:      source,
			Destination: w,
			Options:     options,
			Params:      params,
		},
	}

	_, err = dd.Execute(context.Background(), message)
	return err
}
