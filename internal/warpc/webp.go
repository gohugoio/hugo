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
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"maps"

	"github.com/gohugoio/hugo/common/himage"
	"github.com/gohugoio/hugo/common/hugio"
)

var (
	_ SourceProvider      = WebpInput{}
	_ DestinationProvider = WebpInput{}
)

type CommonImageProcessingParams struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
	Stride int `json:"stride,omitempty"`

	// For animated images.
	FrameDurations []int `json:"frameDurations,omitempty"`
	LoopCount      int   `json:"loopCount,omitempty"`
}

/*
If you're reading this and questioning the protocol on top of WASM using JSON and streams instead of WAS Call's with pointers written to linear memory:

- The goal of this is to eventually make it into one or more RPC plugin APIs.
- Passing pointers around has a number of challenges in that context.
- One would be that it's not possible to pass pointers to child processes (e.g. non-WASM plugins).

Also, you would think that this JSON/streams approach would be significantly slower than using pointers directly, but in practice,
at least for WebP, the difference is negligible, see below output from a test run:

[pointers] DecodeWebp took 18.168375ms
[pointers] EncodeWebp took 13.959458ms
[pointers] DecodeWebpConfig took 93.083µs

[streams] DecodeWebp took 17.192917ms
[streams] EncodeWebp took 14.084792ms
[streams] DecodeWebpConfig took 54.334µs

Also note that the placement of this code in this internal package is also temporary. We 1. Need to get the WASM RPC plugin infrastructure in place, and 2. Need to decide on the final API shape for image processing plugins.
*/
type WebpInput struct {
	Source      hugio.SizeReader `json:"-"`       // Will be sent in a separate stream.
	Destination io.Writer        `json:"-"`       // Will be used to write the result to.
	Options     map[string]any   `json:"options"` // Config options.
	Params      map[string]any   `json:"params"`  // Command params (width, height, etc.).
}

func (w WebpInput) GetSource() hugio.SizeReader {
	return w.Source
}

func (w WebpInput) GetDestination() io.Writer {
	return w.Destination
}

type WebpOutput struct {
	Params CommonImageProcessingParams `json:"params"`
}

type WebpCodec struct {
	d func() (Dispatcher[WebpInput, WebpOutput], error)
}

// Decode reads a WEBP image from r and returns it as an image.Image.
// Note that animated WebP images are returnes as an himage.AnimatedImage.
func (d *WebpCodec) Decode(r io.Reader) (image.Image, error) {
	dd, err := d.d()
	if err != nil {
		return nil, err
	}

	source, err := hugio.ToSizeReader(r)
	if err != nil {
		return nil, err
	}

	var destination bytes.Buffer

	// Commands:
	// encodeNRGBA
	// encodeGray
	// decode
	// config
	message := Message[WebpInput]{
		Header: Header{
			Version:       1,
			Command:       "decode",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: WebpInput{
			Source:      source,
			Destination: &destination,
			Options:     map[string]any{},
		},
	}

	out, err := dd.Execute(context.Background(), message)
	if err != nil {
		return nil, err
	}

	w, h, stride := out.Data.Params.Width, out.Data.Params.Height, out.Data.Params.Stride
	if w == 0 || h == 0 || stride == 0 {
		return nil, fmt.Errorf("received invalid image dimensions: %dx%d stride %d", w, h, stride)
	}

	if len(out.Data.Params.FrameDurations) > 0 {
		// Animated WebP.
		img := &WEBP{
			frameDurations: out.Data.Params.FrameDurations,
			loopCount:      out.Data.Params.LoopCount,
		}

		frameSize := stride * h
		frames := make([]image.Image, len(destination.Bytes())/frameSize)
		for i := 0; i < len(frames); i++ {
			frameBytes := destination.Bytes()[i*frameSize : (i+1)*frameSize]
			frameImg := &image.NRGBA{
				Pix:    frameBytes,
				Stride: stride,
				Rect:   image.Rect(0, 0, w, h),
			}
			frames[i] = frameImg
		}
		img.SetFrames(frames)
		return img, nil
	}

	img := &image.NRGBA{
		Pix:    destination.Bytes(),
		Stride: stride,
		Rect:   image.Rect(0, 0, w, h),
	}

	return img, nil
}

func (d *WebpCodec) DecodeConfig(r io.Reader) (image.Config, error) {
	dd, err := d.d()
	if err != nil {
		return image.Config{}, err
	}

	// Avoid reading the entire image for config only.
	const webpMaxHeaderSize = 32
	b := make([]byte, webpMaxHeaderSize)
	_, err = r.Read(b)
	if err != nil {
		return image.Config{}, err
	}

	message := Message[WebpInput]{
		Header: Header{
			Version:       1,
			Command:       "config",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON},
		},

		Data: WebpInput{
			Source: bytes.NewReader(b),
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

func (d *WebpCodec) Encode(w io.Writer, img image.Image, opts map[string]any) error {
	b := img.Bounds()
	if b.Dx() >= 1<<16 || b.Dy() >= 1<<16 {
		return errors.New("webp: image is too large to encode")
	}

	dd, err := d.d()
	if err != nil {
		return err
	}

	const (
		commandEncodeNRGBA = "encodeNRGBA"
		commandEncodeGray  = "encodeGray"
	)

	var (
		bounds         = img.Bounds()
		imageBytes     []byte
		stride         int
		frameDurations []int
		loopCount      int
		command        string
	)

	switch v := img.(type) {
	case *image.RGBA:
		imageBytes = v.Pix
		stride = v.Stride
		command = commandEncodeNRGBA
	case *WEBP:
		// Animated WebP.
		frames := v.GetFrames()
		if len(frames) == 0 {
			return errors.New("webp: animated image has no frames")
		}
		firstFrame := frames[0]
		bounds = firstFrame.Bounds()
		nrgba := convertToNRGBA(firstFrame)
		frameSize := nrgba.Stride * bounds.Dy()
		imageBytes = make([]byte, frameSize*len(frames))
		stride = nrgba.Stride
		for i, frame := range frames {
			var nrgbaFrame *image.NRGBA
			if i == 0 {
				nrgbaFrame = nrgba
			} else {
				nrgbaFrame = convertToNRGBA(frame)
			}
			copy(imageBytes[i*frameSize:(i+1)*frameSize], nrgbaFrame.Pix)
		}
		frameDurations = v.GetFrameDurations()
		loopCount = v.loopCount
		command = commandEncodeNRGBA
	case *image.Gray:
		imageBytes = v.Pix
		stride = v.Stride
		command = commandEncodeGray
	case himage.AnimatedImage:
		frames := v.GetFrames()
		if len(frames) == 0 {
			return errors.New("webp: animated image has no frames")
		}
		firstFrame := frames[0]
		bounds = firstFrame.Bounds()
		nrgba := convertToNRGBA(firstFrame)
		frameSize := nrgba.Stride * bounds.Dy()
		imageBytes = make([]byte, frameSize*len(frames))
		stride = nrgba.Stride
		for i, frame := range frames {
			var nrgbaFrame *image.NRGBA
			if i == 0 {
				nrgbaFrame = nrgba
			} else {
				nrgbaFrame = convertToNRGBA(frame)
			}
			copy(imageBytes[i*frameSize:(i+1)*frameSize], nrgbaFrame.Pix)
		}
		frameDurations = v.GetFrameDurations()
		loopCount = v.GetLoopCount()
		command = commandEncodeNRGBA
	default:
		nrgba := convertToNRGBA(img)
		imageBytes = nrgba.Pix
		stride = nrgba.Stride
		command = commandEncodeNRGBA

	}

	if len(imageBytes) == 0 {
		return fmt.Errorf("no image bytes extracted from %T", img)
	}

	opts = maps.Clone(opts)
	opts["useSharpYuv"] = true // Use sharp (and slow) RGB->YUV conversion.

	message := Message[WebpInput]{
		Header: Header{
			Version:       1,
			Command:       command,
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: WebpInput{
			Source:      bytes.NewReader(imageBytes),
			Destination: w,
			Options:     opts,
			Params: map[string]any{
				"width":          bounds.Max.X,
				"height":         bounds.Max.Y,
				"stride":         stride,
				"frameDurations": frameDurations,
				"loopCount":      loopCount,
			},
		},
	}

	_, err = dd.Execute(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func convertToNRGBA(src image.Image) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}

var _ himage.AnimatedImage = (*WEBP)(nil)

// WEBP represents an animated WebP image.
// The naming deliberately matches the fields in the standard library image/gif package.
type WEBP struct {
	image.Image    // The first frame.
	frames         []image.Image
	frameDurations []int
	loopCount      int
}

func (w *WEBP) GetLoopCount() int {
	return w.loopCount
}

func (w *WEBP) GetFrames() []image.Image {
	return w.frames
}

func (w *WEBP) GetFrameDurations() []int {
	return w.frameDurations
}

func (w *WEBP) GetRaw() any {
	return w
}

func (w *WEBP) SetFrames(frames []image.Image) {
	if len(frames) == 0 {
		panic("frames cannot be empty")
	}
	w.frames = frames
	w.Image = frames[0]
}

func (w *WEBP) SetWidthHeight(width, height int) {
	// No-op for WEBP.
}
