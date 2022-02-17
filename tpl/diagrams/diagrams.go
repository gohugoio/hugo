// Copyright 2022 The Hugo Authors. All rights reserved.
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

package diagrams

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/bep/goat"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

type SVGDiagram interface {
	Body() template.HTML
	SVG() template.HTML
	Width() int
	Height() int
}

type goatDiagram struct {
	d goat.SVG
}

func (d goatDiagram) Body() template.HTML {
	return template.HTML(d.d.Body)
}

func (d goatDiagram) SVG() template.HTML {
	return template.HTML(d.d.String())
}

func (d goatDiagram) Width() int {
	return d.d.Width
}

func (d goatDiagram) Height() int {
	return d.d.Height
}

type Diagrams struct {
	d *deps.Deps
}

func (d *Diagrams) Goat(v interface{}) SVGDiagram {
	var r io.Reader

	switch vv := v.(type) {
	case io.Reader:
		r = vv
	case []byte:
		r = bytes.NewReader(vv)
	default:
		r = strings.NewReader(cast.ToString(v))
	}

	return goatDiagram{
		d: goat.BuildSVG(r),
	}
}
