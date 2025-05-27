// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"context"
	"encoding/xml"
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/dreampuf/mermaid.go"
	"github.com/spf13/cast"
)

type mermaidDiagram struct {
	Content string
	width   int
	height  int
}

type SVGRoot struct {
	Width    string   `xml:"width,attr"`
	Height   string   `xml:"height,attr"`
	ViewBox  string   `xml:"viewBox,attr"`
	innerSVG innerSVG `xml:"svg"`
}

type innerSVG struct {
	XML string `xml:",innerxml"`
}

func (m mermaidDiagram) Inner() template.HTML {
	parsed := SVGRoot{}
	xml.Unmarshal([]byte(m.Content), &parsed)
	inner, _ := xml.Marshal(parsed.innerSVG)
	return template.HTML(inner)
}

func (m mermaidDiagram) Wrapped() template.HTML {
	return template.HTML(m.Content)
}

func (m mermaidDiagram) Width() int {
	return m.width
}

func (m mermaidDiagram) Height() int {
	return m.height
}

// Mermaid creates a new SVG diagram from input v.
func (d *Namespace) Mermaid(v any) SVGDiagram {
	var r io.Reader

	switch vv := v.(type) {
	case io.Reader:
		r = vv
	case []byte:
		r = bytes.NewReader(vv)
	default:
		r = strings.NewReader(cast.ToString(v))
	}

	re, _ := mermaid_go.NewRenderEngine(context.TODO())
	defer re.Cancel()

	buf := new(strings.Builder)
	io.Copy(buf, r)
	result, _ := re.Render(buf.String())

	parsed := SVGRoot{}
	xml.Unmarshal([]byte(result), &parsed)

	if parsed.Height == "" || parsed.Height[len(parsed.Height)-1:] == "%" {
		parsed.Height = strings.Split(parsed.ViewBox, " ")[3]
	}

	if parsed.Width == "" || parsed.Width[len(parsed.Width)-1:] == "%" {
		parsed.Width = strings.Split(parsed.ViewBox, " ")[2]
	}
	height, _ := strconv.Atoi(parsed.Height)
	width, _ := strconv.Atoi(parsed.Width)

	return mermaidDiagram{
		Content: result,
		height:  height,
		width:   width,
	}
}
