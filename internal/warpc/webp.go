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
	"io"
	"time"

	"github.com/gohugoio/hugo/common/hdebug"
	"github.com/gohugoio/hugo/common/hugio"
)

var (
	_ SourceProvider      = WebpInput{}
	_ DestinationProvider = WebpInput{}
)

type CommonImageProcessingOptions struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
	Stride int `json:"stride,omitempty"`
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
	Source      hugio.SizeReader `json:"-"` // Will be sent in a separate stream.
	Destination io.Writer        `json:"-"` // Will be used to write the result to.
	Options     map[string]any   `json:"options"`

	// TODO1 config optioGetSourcens.
}

func (w WebpInput) GetSource() hugio.SizeReader {
	return w.Source
}

func (w WebpInput) GetDestination() io.Writer {
	return w.Destination
}

type WebpOutput struct {
	Options CommonImageProcessingOptions
}

func stopClock(what string, start time.Time) {
	hdebug.Printf("%s took %s", what, time.Since(start))
}

// TODO1 in webp wasm scrip, do a bare clone of libwebp.

// Decode reads a WEBP image from r and returns it as an image.Image.
// TODO1 remember to increment the webp hash id so people get new wasm versions.
func (d *Dispatchers) DecodeWebp(r io.Reader) (image.Image, error) {
	dd, err := d.Webp()
	if err != nil {
		return nil, err
	}
	start := time.Now()
	defer stopClock("DecodeWebp", start)

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
			ID:            d.id.Add(1),
			Command:       "decode",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: WebpInput{
			Source:      source,
			Destination: &destination,
			Options:     map[string]any{
				// TODO1
			},
		},
	}

	out, err := dd.Execute(context.Background(), message)
	if err != nil {
		return nil, err
	}

	w, h, stride := out.Data.Options.Width, out.Data.Options.Height, out.Data.Options.Stride
	if w == 0 || h == 0 || stride == 0 {
		return nil, fmt.Errorf("received invalid image dimensions: %dx%d stride %d", w, h, stride)
	}

	img := &image.RGBA{
		Pix:    destination.Bytes(),
		Stride: stride,
		Rect:   image.Rect(0, 0, w, h),
	}

	return img, nil
}

func (d *Dispatchers) DecodeWebpConfig(r io.Reader) (image.Config, error) {
	dd, err := d.Webp()
	if err != nil {
		return image.Config{}, err
	}
	start := time.Now()
	defer stopClock("DecodeWebpConfig", start)

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
			ID:            d.id.Add(1),
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
		Width:      out.Data.Options.Width,
		Height:     out.Data.Options.Height,
		ColorModel: color.RGBAModel, // TODO1
	}, nil
}

func (d *Dispatchers) EncodeWebp(w io.Writer, src image.Image) error {
	dd, err := d.Webp()
	if err != nil {
		return err
	}
	start := time.Now()
	defer stopClock("EncodeWebp", start)

	var (
		bounds     = src.Bounds()
		imageBytes []byte
		stride     int
	)

	switch v := src.(type) {
	case *image.RGBA:
		imageBytes = v.Pix
		stride = v.Stride
	default:
		hdebug.Panicf("unsupported %T", src)
	}

	if len(imageBytes) == 0 {
		return fmt.Errorf("no image bytes extracted from %T", src)
	}

	// Commands:
	// encodeNRGBA
	// encodeGray TODO1
	// decode
	// config
	message := Message[WebpInput]{
		Header: Header{
			Version:       1,
			ID:            d.id.Add(1),
			Command:       "encodeNRGBA",
			RequestKinds:  []string{MessageKindJSON, MessageKindBlob},
			ResponseKinds: []string{MessageKindJSON, MessageKindBlob},
		},

		Data: WebpInput{
			Source:      bytes.NewReader(imageBytes),
			Destination: w,
			Options: map[string]any{
				"width":  bounds.Max.X,
				"height": bounds.Max.Y,
				"stride": stride,
			},
		},
	}

	_, err = dd.Execute(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}
