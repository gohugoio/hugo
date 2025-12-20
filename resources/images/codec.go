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

package images

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/gohugoio/hugo/common/himage"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// Decoder defines the decoding of an image format.
// These matches the globalmpackage image functions.
type Decoder interface {
	Decode(r io.Reader) (image.Image, error)
	DecodeConfig(r io.Reader) (image.Config, error)
}

// Encoder defines the encoding of an image format.
type Encoder interface {
	Encode(w io.Writer, src image.Image) error
}

type ToEncoder interface {
	EncodeTo(conf ImageConfig, w io.Writer, src image.Image) error
}

// CodecStdlib defines both decoding and encoding of an image format as defined by the standard library.
type CodecStdlib interface {
	Decoder
	Encoder
}

// Codec is a generic image codec supporting multiple formats.
type Codec struct {
	webp CodecStdlib
}

func newCodec(webp CodecStdlib) *Codec {
	return &Codec{webp: webp}
}

func (d *Codec) EncodeTo(conf ImageConfig, w io.Writer, img image.Image) error {
	switch conf.TargetFormat {
	case JPEG:
		var rgba *image.RGBA
		quality := conf.Quality

		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			return jpeg.Encode(w, rgba, &jpeg.Options{Quality: quality})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	case PNG:
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(w, img)
	case GIF:
		if anim, ok := img.(himage.AnimatedImage); ok {
			if g, ok := anim.GetRaw().(*gif.GIF); ok {
				return gif.EncodeAll(w, g)
			}

			// Animated image, but not a GIF. Convert it.
			frames := anim.GetFrames()
			if len(frames) == 0 {
				return gif.Encode(w, img, &gif.Options{NumColors: 256})
			}

			frameDurations := anim.GetFrameDurations()
			if len(frameDurations) != len(frames) {
				return errors.New("gif: number of frame durations does not match number of frames")
			}

			outGif := &gif.GIF{
				Delay: himage.FrameDurationsToGifDelays(frameDurations),
			}

			outGif.LoopCount = anim.GetLoopCount()

			for _, frame := range frames {
				bounds := frame.Bounds()
				palettedImage := image.NewPaletted(bounds, palette.Plan9)
				draw.Draw(palettedImage, palettedImage.Rect, frame, bounds.Min, draw.Src)
				outGif.Image = append(outGif.Image, palettedImage)
			}

			return gif.EncodeAll(w, outGif)
		}
		return gif.Encode(w, img, &gif.Options{
			NumColors: 256,
		})
	case TIFF:
		return tiff.Encode(w, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
	case BMP:
		return bmp.Encode(w, img)
	case WEBP:
		return d.webp.Encode(w, img)
	default:
		return errors.New("format not supported")
	}
}

func (d *Codec) DecodeFormat(f Format, r io.Reader) (image.Image, error) {
	switch f {
	case JPEG, PNG:
		// We reworked this decode/encode setup to get full WebP support in v0.153.0.
		// In the first take of that we used f to decide whether to call png.Decode or jpeg.Decode here,
		// but testing it on some sites, it seems that it's not uncommon to store JPEGs with PNG extensions and vice versa.
		// So, to reduce some noise in that release, we fallback to the standard library here,
		// which will read the magic bytes and decode accordingly.
		img, _, err := image.Decode(r)
		return img, err
	case GIF:
		g, err := gif.DecodeAll(r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode gif: %w", err)
		}
		if len(g.Delay) > 1 {
			return &giphy{gif: g, Image: g.Image[0]}, nil
		}
		return g.Image[0], nil
	case TIFF:
		return tiff.Decode(r)
	case BMP:
		return bmp.Decode(r)
	case WEBP:
		img, err := d.webp.Decode(r)
		if err == nil {
			return img, nil
		}
		if rs, ok := r.(io.ReadSeeker); ok {
			// See issue 14288. Turns out it's not uncommon to e.g. name their PNG files with a WEBP extension.
			// With the old Go's webp decoder, this didn't fail (it looked for the file header),
			// but now some error has surfaced.
			// To reduce some noise, we try to reset and decode again using the standard library.
			_, err2 := rs.Seek(0, io.SeekStart)
			if err2 != nil {
				return nil, err
			}
			img, _, err2 = image.Decode(rs)
			if err2 == nil {
				return img, nil
			}
		}
		return nil, err
	default:
		return nil, errors.New("format not supported")
	}
}

func (d *Codec) Decode(r io.Reader) (image.Image, error) {
	rr := toPeekReader(r)
	format, err := formatFromImage(rr)
	if err != nil {
		return nil, err
	}
	if format != 0 {
		return d.DecodeFormat(format, rr)
	}

	// Fallback to the standard image.Decode.
	img, _, err := image.Decode(rr)
	return img, err
}

func (d *Codec) DecodeConfig(r io.Reader) (image.Config, string, error) {
	rr := toPeekReader(r)
	format, err := formatFromImage(rr)
	if err != nil {
		return image.Config{}, "", err
	}
	if format == WEBP {
		cfg, err := d.webp.DecodeConfig(rr)
		return cfg, "webp", err
	}

	// Fallback to the standard image.DecodeConfig.
	conf, name, err := image.DecodeConfig(rr)
	return conf, name, err
}

// toPeekReader converts an io.Reader to a peekReader.
func toPeekReader(r io.Reader) peekReader {
	if rr, ok := r.(peekReader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

// A peekReader is an io.Reader that can also peek ahead.
type peekReader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

const (
	// The WebP file header is 12 bytes long and starts with "RIFF" followed by
	// 4 bytes indicating the file size, followed by "WEBP" and the VP8 chunk header.
	// We use '?' as a wildcard for the 4 size bytes.
	magicWebp = "RIFF????WEBPVP8"
	// The GIF file header is 6 bytes long and starts with "GIF87a" or "GIF89a".
	magicGif = "GIF8???"
)

type magicFormat struct {
	magic  string
	format Format
}

var magicFormats = []magicFormat{
	{magic: magicWebp, format: WEBP},
	{magic: magicGif, format: GIF},
}

// formatFromImage determines the image format from the magic bytes.
// Note that this is only a partial implementation,
// as we currently only need WebP and GIF detection.
// The others can be handled by the standard library.
func formatFromImage(r peekReader) (Format, error) {
	for _, mf := range magicFormats {
		magicLen := len(mf.magic)
		b, err := r.Peek(magicLen)
		if err == nil && match(mf.magic, b) {
			return mf.format, nil
		}
	}
	return 0, nil
}

func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i := 0; i < len(magic); i++ {
		if magic[i] != '?' && magic[i] != b[i] {
			return false
		}
	}
	return true
}

var (
	_ himage.AnimatedImage       = (*giphy)(nil)
	_ himage.ImageConfigProvider = (*giphy)(nil)
)

type giphy struct {
	image.Image
	gif *gif.GIF
}

func (g *giphy) GetRaw() any {
	return g.gif
}

func (g *giphy) GetLoopCount() int {
	return g.gif.LoopCount
}

func (g *giphy) GetFrames() []image.Image {
	frames := make([]image.Image, len(g.gif.Image))
	for i, frame := range g.gif.Image {
		frames[i] = frame
	}
	return frames
}

func (g *giphy) GetImageConfig() image.Config {
	return g.gif.Config
}

func (g *giphy) SetFrames(frames []image.Image) {
	if len(frames) == 0 {
		panic("frames cannot be empty")
	}
	g.gif.Image = make([]*image.Paletted, len(frames))
	for i, frame := range frames {
		g.gif.Image[i] = frame.(*image.Paletted)
	}
	g.Image = g.gif.Image[0]
}

func (g *giphy) SetWidthHeight(width, height int) {
	g.gif.Config.Width = width
	g.gif.Config.Height = height
}

func (g *giphy) GetFrameDurations() []int {
	return himage.GifDelaysToFrameDurations(g.gif.Delay)
}
