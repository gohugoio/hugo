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

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/common/hcontext"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/langs"

	htmltemplate "github.com/gohugoio/hugo/tpl/internal/go_templates/htmltemplate"
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
)

// Context manages values passed in the context to templates.
var Context = struct {
	DependencyManagerScopedProvider    hcontext.ContextDispatcher[identity.DependencyManagerScopedProvider]
	GetDependencyManagerInCurrentScope func(context.Context) identity.Manager
	DependencyScope                    hcontext.ContextDispatcher[int]
	Page                               hcontext.ContextDispatcher[page]
	IsInGoldmark                       hcontext.ContextDispatcher[bool]
	CurrentTemplate                    hcontext.ContextDispatcher[*CurrentTemplateInfo]
}{
	DependencyManagerScopedProvider: hcontext.NewContextDispatcher[identity.DependencyManagerScopedProvider](contextKeyDependencyManagerScopedProvider),
	DependencyScope:                 hcontext.NewContextDispatcher[int](contextKeyDependencyScope),
	Page:                            hcontext.NewContextDispatcher[page](contextKeyPage),
	IsInGoldmark:                    hcontext.NewContextDispatcher[bool](contextKeyIsInGoldmark),
	CurrentTemplate:                 hcontext.NewContextDispatcher[*CurrentTemplateInfo](cntextKeyCurrentTemplateInfo),
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

type page interface {
	IsNode() bool
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

const hugoNewLinePlaceholder = "___hugonl_"

var stripHTMLReplacerPre = strings.NewReplacer("\n", " ", "</p>", hugoNewLinePlaceholder, "<br>", hugoNewLinePlaceholder, "<br />", hugoNewLinePlaceholder)

// StripHTML strips out all HTML tags in s.
func StripHTML(s string) string {
	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}

	pre := stripHTMLReplacerPre.Replace(s)
	preReplaced := pre != s

	s = htmltemplate.StripTags(pre)

	if preReplaced {
		s = strings.ReplaceAll(s, hugoNewLinePlaceholder, "\n")
	}

	var wasSpace bool
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	for _, r := range s {
		isSpace := unicode.IsSpace(r)
		if !(isSpace && wasSpace) {
			b.WriteRune(r)
		}
		wasSpace = isSpace
	}

	if b.Len() > 0 {
		s = b.String()
	}

	return s
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
