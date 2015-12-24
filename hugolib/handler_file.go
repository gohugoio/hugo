// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"regexp"

	"log"

	"github.com/bamiaux/rez"
	"github.com/dchest/cssmin"
	"github.com/oliamb/cutter"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/tpl"
)

func init() {
	RegisterHandler(new(cssHandler))
	RegisterHandler(new(imageHandler))
	RegisterHandler(new(defaultHandler))
}

type basicFileHandler Handle

func (h basicFileHandler) Read(f *source.File, s *Site) HandledResult {
	return HandledResult{file: f}
}

func (h basicFileHandler) PageConvert(*Page, tpl.Template) HandledResult {
	return HandledResult{}
}

type defaultHandler struct{ basicFileHandler }

func (h defaultHandler) Extensions() []string { return []string{"*"} }
func (h defaultHandler) FileConvert(f *source.File, s *Site) HandledResult {
	s.WriteDestFile(f.Path(), f.Contents)
	return HandledResult{file: f}
}

type cssHandler struct{ basicFileHandler }

func (h cssHandler) Extensions() []string { return []string{"css"} }
func (h cssHandler) FileConvert(f *source.File, s *Site) HandledResult {
	x := cssmin.Minify(f.Bytes())
	s.WriteDestFile(f.Path(), helpers.BytesToReader(x))
	return HandledResult{file: f}
}

type imageHandler struct{ basicFileHandler }

func (h imageHandler) Extensions() []string { return []string{"jpg", "jpeg", "png"} }
func (h imageHandler) FileConvert(file *source.File, s *Site) HandledResult {
	props := helpers.ParseImageProps(file.BaseFileName())

	if props.Width != 0 || props.Height != 0 {
		input, ext, _ := image.Decode(file.Contents)
		inputRect := input.Bounds()

		width := inputRect.Max.X
		height := inputRect.Max.Y
		if props.Width > 0 {
			width = props.Width
		}
		if props.Height > 0 {
			height = props.Height
		}

		var output image.Image
		var err error
		if props.Crop {
			output, err = cutter.Crop(input, cutter.Config{
				Width:  width,
				Height: height,
			})
		} else {
			output = image.NewRGBA(image.Rect(0, 0, width, height))
			err = rez.Convert(output, input, rez.NewBicubicFilter())
		}

		if err != nil {
			log.Println("Error converting/cropping image", err)
		}

		buf := new(bytes.Buffer)
		switch ext {
		case "jpg", "jpeg":
			_ = jpeg.Encode(buf, output, nil)
		case "png":
			_ = png.Encode(buf, output)
		}

		re := regexp.MustCompile(`\.(w|h)\d(-)?.+\.`)
		path := file.Path()
		path = re.ReplaceAllString(path, ".")
		file = source.NewFileWithContents(path, bytes.NewReader(buf.Bytes()))
	}
	s.WriteDestFile(file.Path(), file.Contents)
	return HandledResult{file: file}
}
