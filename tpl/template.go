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

// Package tpl contains template functions and related types.
package tpl

import (
	"context"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/bep/helpers/contexthelpers"
	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/hhtml"

	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"

	texttemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
)

// Template is the common interface between text/template and html/template.
type Template interface {
	Name() string
	Prepare() (*texttemplate.Template, error)
}

// RenderingContext represents the currently rendered site/language.
type RenderingContext struct {
	Site       site
	SiteOutIdx int
}

type (
	contextKey uint8
)

const (
	contextKeyDependencyManagerScopedProvider contextKey = iota
	contextKeyDependencyScope
	contextKeyPage
	contextKeyIsInGoldmark
	cntextKeyCurrentTemplateInfo
	contextKeyPartialDecoratorIDStack
	contextKeyIsInPartialCached
)

// Context manages values passed in the context to templates.
var Context = struct {
	DependencyManagerScopedProvider    contexthelpers.ContextDispatcher[identity.DependencyManagerScopedProvider]
	GetDependencyManagerInCurrentScope func(context.Context) identity.Manager
	DependencyScope                    contexthelpers.ContextDispatcher[int]
	Page                               contexthelpers.ContextDispatcher[page]
	IsInGoldmark                       contexthelpers.ContextDispatcher[bool]
	CurrentTemplate                    contexthelpers.ContextDispatcher[*CurrentTemplateInfo]
	PartialDecoratorIDStack            contexthelpers.ContextDispatcher[*collections.Stack[*StringBool]]
	IsInPartialCached                  contexthelpers.ContextDispatcher[bool]
}{
	DependencyManagerScopedProvider: contexthelpers.NewContextDispatcher[identity.DependencyManagerScopedProvider](contextKeyDependencyManagerScopedProvider),
	DependencyScope:                 contexthelpers.NewContextDispatcher[int](contextKeyDependencyScope),
	Page:                            contexthelpers.NewContextDispatcher[page](contextKeyPage),
	IsInGoldmark:                    contexthelpers.NewContextDispatcher[bool](contextKeyIsInGoldmark),
	CurrentTemplate:                 contexthelpers.NewContextDispatcher[*CurrentTemplateInfo](cntextKeyCurrentTemplateInfo),
	PartialDecoratorIDStack:         contexthelpers.NewContextDispatcher[*collections.Stack[*StringBool]](contextKeyPartialDecoratorIDStack),
	IsInPartialCached:               contexthelpers.NewContextDispatcher[bool](contextKeyIsInPartialCached),
}

func init() {
	Context.GetDependencyManagerInCurrentScope = func(ctx context.Context) identity.Manager {
		idmsp := Context.DependencyManagerScopedProvider.Get(ctx)
		if idmsp != nil {
			return idmsp.GetDependencyManagerForScope(Context.DependencyScope.Get(ctx))
		}
		return nil
	}
}

// StringBool is a helper struct to hold a string and a bool value.
type StringBool struct {
	Str  string
	Bool bool
}

type page interface {
	IsBranch() bool
}

type site interface {
	Language() *langs.Language
}

const (
	// HugoDeferredTemplatePrefix is the prefix for placeholders for deferred templates.
	HugoDeferredTemplatePrefix = "__hdeferred/"
	// HugoDeferredTemplateSuffix is the suffix for placeholders for deferred templates.
	HugoDeferredTemplateSuffix = "__d="
)

// StripHTML strips out all HTML tags, comments and the content of script and
// style elements in s. Paragraph ends and line breaks are replaced with a
// newline, other newlines with a space, text in adjacent block elements is
// separated by a space, and repeated whitespace is collapsed.
// Strings without any '<' are returned unchanged.
func StripHTML(s string) string {
	if strings.IndexByte(s, '<') == -1 {
		return s
	}

	w := stripHTMLWriter{}
	w.b.Grow(len(s))
	sc := hhtml.NewScanner(s, 0)
	for {
		textStart := sc.Pos()
		t, ok := sc.Next(len(s))
		if !ok {
			w.writeText(s[textStart:])
			break
		}
		w.writeText(s[textStart:t.Start])
		switch {
		case t.Type == hhtml.StartTag && t.Name == "br", t.Type == hhtml.EndTag && t.Name == "p":
			w.writeSpace('\n')
		case t.Type == hhtml.StartTag && (t.Name == "script" || t.Name == "style"):
			if !t.SelfClosing {
				sc.SkipRawText(len(s), t.Name)
			}
		case hhtml.IsBlockElement(t.Name):
			// Keep text in e.g. adjacent list items apart.
			w.needsSpace = true
		}
	}

	return w.b.String()
}

type stripHTMLWriter struct {
	b        strings.Builder
	wasSpace bool
	// needsSpace is set when a space should separate the next text
	// from any text before it.
	needsSpace bool
}

func (w *stripHTMLWriter) writeText(s string) {
	for i := 0; i < len(s); {
		c := s[i]
		if c >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if unicode.IsSpace(r) {
				if !w.wasSpace {
					w.b.WriteString(s[i : i+size])
					w.wasSpace = true
				}
			} else {
				w.writeNonSpace(s[i : i+size])
			}
			i += size
			continue
		}
		if isASCIISpace(c) {
			if c == '\n' {
				c = ' '
			}
			w.writeSpace(c)
			i++
			continue
		}
		j := i + 1
		for j < len(s) && s[j] < utf8.RuneSelf && !isASCIISpace(s[j]) {
			j++
		}
		w.writeNonSpace(s[i:j])
		i = j
	}
}

func (w *stripHTMLWriter) writeNonSpace(s string) {
	if w.needsSpace && !w.wasSpace && w.b.Len() > 0 {
		w.b.WriteByte(' ')
	}
	w.b.WriteString(s)
	w.wasSpace = false
	w.needsSpace = false
}

func (w *stripHTMLWriter) writeSpace(c byte) {
	if !w.wasSpace {
		w.b.WriteByte(c)
		w.wasSpace = true
	}
}

func isASCIISpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// DeferredExecution holds the template and data for a deferred execution.
type DeferredExecution struct {
	Mu           sync.Mutex
	Ctx          context.Context
	TemplatePath string
	Data         any

	Executed bool
	Result   string
}

type CurrentTemplateInfoOps interface {
	CurrentTemplateInfoCommonOps
	Base() CurrentTemplateInfoCommonOps
}

type CurrentTemplateInfoCommonOps interface {
	// Template name.
	Name() string
	// Template source filename.
	// Will be empty for internal templates.
	Filename() string
}

// CurrentTemplateInfo as returned in templates.Current.
type CurrentTemplateInfo struct {
	Parent *CurrentTemplateInfo
	Level  int
	Key    string
	CurrentTemplateInfoOps
}

// CurrentTemplateInfos is a slice of CurrentTemplateInfo.
type CurrentTemplateInfos []*CurrentTemplateInfo

// Reverse creates a copy of the slice and reverses it.
func (c CurrentTemplateInfos) Reverse() CurrentTemplateInfos {
	if len(c) == 0 {
		return c
	}
	r := make(CurrentTemplateInfos, len(c))
	copy(r, c)
	slices.Reverse(r)
	return r
}

// Ancestors returns the ancestors of the current template.
func (ti *CurrentTemplateInfo) Ancestors() CurrentTemplateInfos {
	var ancestors []*CurrentTemplateInfo
	for ti.Parent != nil {
		ti = ti.Parent
		ancestors = append(ancestors, ti)
	}
	return ancestors
}
