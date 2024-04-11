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

package attributes

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/spf13/cast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// Markdown attributes used as options by the Chroma highlighter.
var chromaHighlightProcessingAttributes = map[string]bool{
	"anchorLineNos":      true,
	"guessSyntax":        true,
	"hl_Lines":           true,
	"hl_inline":          true,
	"lineAnchors":        true,
	"lineNos":            true,
	"lineNoStart":        true,
	"lineNumbersInTable": true,
	"noClasses":          true,
	"nohl":               true,
	"style":              true,
	"tabWidth":           true,
}

func init() {
	for k, v := range chromaHighlightProcessingAttributes {
		chromaHighlightProcessingAttributes[strings.ToLower(k)] = v
	}
}

type AttributesOwnerType int

const (
	AttributesOwnerGeneral AttributesOwnerType = iota
	AttributesOwnerCodeBlockChroma
	AttributesOwnerCodeBlockCustom
)

func New(astAttributes []ast.Attribute, ownerType AttributesOwnerType) *AttributesHolder {
	var (
		attrs []Attribute
		opts  []Attribute
	)
	for _, v := range astAttributes {
		nameLower := strings.ToLower(string(v.Name))
		if strings.HasPrefix(string(nameLower), "on") {
			continue
		}
		var vv any
		switch vvv := v.Value.(type) {
		case bool, float64:
			vv = vvv
		case []any:
			// Highlight line number hlRanges.
			var hlRanges [][2]int
			for _, l := range vvv {
				if ln, ok := l.(float64); ok {
					hlRanges = append(hlRanges, [2]int{int(ln) - 1, int(ln) - 1})
				} else if rng, ok := l.([]uint8); ok {
					slices := strings.Split(string([]byte(rng)), "-")
					lhs, err := strconv.Atoi(slices[0])
					if err != nil {
						continue
					}
					rhs := lhs
					if len(slices) > 1 {
						rhs, err = strconv.Atoi(slices[1])
						if err != nil {
							continue
						}
					}
					hlRanges = append(hlRanges, [2]int{lhs - 1, rhs - 1})
				}
			}
			vv = hlRanges
		case []byte:
			// Note that we don't do any HTML escaping here.
			// We used to do that, but that changed in #9558.
			// Now it's up to the templates to decide.
			vv = string(vvv)
		default:
			panic(fmt.Sprintf("not implemented: %T", vvv))
		}

		if ownerType == AttributesOwnerCodeBlockChroma && chromaHighlightProcessingAttributes[nameLower] {
			attr := Attribute{Name: string(v.Name), Value: vv}
			opts = append(opts, attr)
		} else {
			attr := Attribute{Name: nameLower, Value: vv}
			attrs = append(attrs, attr)
		}

	}

	return &AttributesHolder{
		attributes: attrs,
		options:    opts,
	}
}

type Attribute struct {
	Name  string
	Value any
}

func (a Attribute) ValueString() string {
	return cast.ToString(a.Value)
}

// Empty holds no attributes.
var Empty = &AttributesHolder{}

type AttributesHolder struct {
	// What we get from Goldmark.
	attributes []Attribute

	// Attributes considered to be an option (code blocks)
	options []Attribute

	// What we send to the render hooks.
	attributesMapInit sync.Once
	attributesMap     map[string]any
	optionsMapInit    sync.Once
	optionsMap        map[string]any
}

type Attributes map[string]any

func (a *AttributesHolder) Attributes() map[string]any {
	a.attributesMapInit.Do(func() {
		a.attributesMap = make(map[string]any)
		for _, v := range a.attributes {
			a.attributesMap[v.Name] = v.Value
		}
	})
	return a.attributesMap
}

func (a *AttributesHolder) Options() map[string]any {
	a.optionsMapInit.Do(func() {
		a.optionsMap = make(map[string]any)
		for _, v := range a.options {
			a.optionsMap[v.Name] = v.Value
		}
	})
	return a.optionsMap
}

func (a *AttributesHolder) AttributesSlice() []Attribute {
	return a.attributes
}

func (a *AttributesHolder) OptionsSlice() []Attribute {
	return a.options
}

// RenderASTAttributes writes the AST attributes to the given as attributes to an HTML element.
// This is used by the default HTML renderers, e.g. for headings etc. where no hook template could be found.
// This performs HTML escaping of string attributes.
func RenderASTAttributes(w hugio.FlexiWriter, attributes ...ast.Attribute) {
	for _, attr := range attributes {

		a := strings.ToLower(string(attr.Name))
		if strings.HasPrefix(a, "on") {
			continue
		}

		_, _ = w.WriteString(" ")
		_, _ = w.Write(attr.Name)
		_, _ = w.WriteString(`="`)

		switch v := attr.Value.(type) {
		case []byte:
			_, _ = w.Write(util.EscapeHTML(v))
		default:
			w.WriteString(cast.ToString(v))
		}

		_ = w.WriteByte('"')
	}
}

// RenderAttributes Render writes the attributes to the given as attributes to an HTML element.
// This is used for the default codeblock rendering.
// This performs HTML escaping of string attributes.
func RenderAttributes(w hugio.FlexiWriter, skipClass bool, attributes ...Attribute) {
	for _, attr := range attributes {
		a := strings.ToLower(string(attr.Name))
		if skipClass && a == "class" {
			continue
		}
		_, _ = w.WriteString(" ")
		_, _ = w.WriteString(attr.Name)
		_, _ = w.WriteString(`="`)

		switch v := attr.Value.(type) {
		case []byte:
			_, _ = w.Write(util.EscapeHTML(v))
		default:
			w.WriteString(cast.ToString(v))
		}

		_ = w.WriteByte('"')
	}
}
