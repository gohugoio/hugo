// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"fmt"
	"strconv"

	"github.com/gohugoio/hugo/helpers"

	"html/template"
	"path"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/pkg/errors"

	"reflect"

	"regexp"
	"sort"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/gohugoio/hugo/resources/page"

	_errors "github.com/pkg/errors"

	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/common/urls"
	"github.com/gohugoio/hugo/output"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/tpl"
)

var (
	_ urls.RefLinker  = (*ShortcodeWithPage)(nil)
	_ pageWrapper     = (*ShortcodeWithPage)(nil)
	_ text.Positioner = (*ShortcodeWithPage)(nil)
)

// ShortcodeWithPage is the "." context in a shortcode template.
type ShortcodeWithPage struct {
	Params        interface{}
	Inner         template.HTML
	Page          page.Page
	Parent        *ShortcodeWithPage
	Name          string
	IsNamedParams bool

	// Zero-based ordinal in relation to its parent. If the parent is the page itself,
	// this ordinal will represent the position of this shortcode in the page content.
	Ordinal int

	// pos is the position in bytes in the source file. Used for error logging.
	posInit   sync.Once
	posOffset int
	pos       text.Position

	scratch *maps.Scratch
}

// Position returns this shortcode's detailed position. Note that this information
// may be expensive to calculate, so only use this in error situations.
func (scp *ShortcodeWithPage) Position() text.Position {
	scp.posInit.Do(func() {
		if p, ok := mustUnwrapPage(scp.Page).(pageContext); ok {
			scp.pos = p.posOffset(scp.posOffset)
		}
	})
	return scp.pos
}

// Site returns information about the current site.
func (scp *ShortcodeWithPage) Site() page.Site {
	return scp.Page.Site()
}

// Ref is a shortcut to the Ref method on Page. It passes itself as a context
// to get better error messages.
func (scp *ShortcodeWithPage) Ref(args map[string]interface{}) (string, error) {
	return scp.Page.RefFrom(args, scp)
}

// RelRef is a shortcut to the RelRef method on Page. It passes itself as a context
// to get better error messages.
func (scp *ShortcodeWithPage) RelRef(args map[string]interface{}) (string, error) {
	return scp.Page.RelRefFrom(args, scp)
}

// Scratch returns a scratch-pad scoped for this shortcode. This can be used
// as a temporary storage for variables, counters etc.
func (scp *ShortcodeWithPage) Scratch() *maps.Scratch {
	if scp.scratch == nil {
		scp.scratch = maps.NewScratch()
	}
	return scp.scratch
}

// Get is a convenience method to look up shortcode parameters by its key.
func (scp *ShortcodeWithPage) Get(key interface{}) interface{} {
	if scp.Params == nil {
		return nil
	}
	if reflect.ValueOf(scp.Params).Len() == 0 {
		return nil
	}

	var x reflect.Value

	switch key.(type) {
	case int64, int32, int16, int8, int:
		if reflect.TypeOf(scp.Params).Kind() == reflect.Map {
			// We treat this as a non error, so people can do similar to
			// {{ $myParam := .Get "myParam" | default .Get 0 }}
			// Without having to do additional checks.
			return nil
		} else if reflect.TypeOf(scp.Params).Kind() == reflect.Slice {
			idx := int(reflect.ValueOf(key).Int())
			ln := reflect.ValueOf(scp.Params).Len()
			if idx > ln-1 {
				return ""
			}
			x = reflect.ValueOf(scp.Params).Index(idx)
		}
	case string:
		if reflect.TypeOf(scp.Params).Kind() == reflect.Map {
			x = reflect.ValueOf(scp.Params).MapIndex(reflect.ValueOf(key))
			if !x.IsValid() {
				return ""
			}
		} else if reflect.TypeOf(scp.Params).Kind() == reflect.Slice {
			// We treat this as a non error, so people can do similar to
			// {{ $myParam := .Get "myParam" | default .Get 0 }}
			// Without having to do additional checks.
			return nil
		}
	}

	return x.Interface()

}

func (scp *ShortcodeWithPage) page() page.Page {
	return scp.Page
}

// Note - this value must not contain any markup syntax
const shortcodePlaceholderPrefix = "HAHAHUGOSHORTCODE"

func createShortcodePlaceholder(id string, ordinal int) string {
	return shortcodePlaceholderPrefix + "-" + id + strconv.Itoa(ordinal) + "-HBHB"
}

type shortcode struct {
	name      string
	isInline  bool          // inline shortcode. Any inner will be a Go template.
	isClosing bool          // whether a closing tag was provided
	inner     []interface{} // string or nested shortcode
	params    interface{}   // map or array
	ordinal   int
	err       error

	info tpl.Info

	// If set, the rendered shortcode is sent as part of the surrounding content
	// to Blackfriday and similar.
	// Before Hug0 0.55 we didn't send any shortcode output to the markup
	// renderer, and this flag told Hugo to process the {{ .Inner }} content
	// separately.
	// The old behaviour can be had by starting your shortcode template with:
	//    {{ $_hugo_config := `{ "version": 1 }`}}
	doMarkup bool

	// the placeholder in the source when passed to Blackfriday etc.
	// This also identifies the rendered shortcode.
	placeholder string

	pos    int // the position in bytes in the source file
	length int // the length in bytes in the source file
}

func (s shortcode) insertPlaceholder() bool {
	return !s.doMarkup || s.configVersion() == 1
}

func (s shortcode) configVersion() int {
	if s.info == nil {
		// Not set for inline shortcodes.
		return 2
	}

	return s.info.ParseInfo().Config.Version
}

func (s shortcode) innerString() string {
	var sb strings.Builder

	for _, inner := range s.inner {
		sb.WriteString(inner.(string))
	}

	return sb.String()
}

func (sc shortcode) String() string {
	// for testing (mostly), so any change here will break tests!
	var params interface{}
	switch v := sc.params.(type) {
	case map[string]interface{}:
		// sort the keys so test assertions won't fail
		var keys []string
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var tmp = make(map[string]interface{})

		for _, k := range keys {
			tmp[k] = v[k]
		}
		params = tmp

	default:
		// use it as is
		params = sc.params
	}

	return fmt.Sprintf("%s(%q, %t){%s}", sc.name, params, sc.doMarkup, sc.inner)
}

type shortcodeHandler struct {
	p *pageState

	s *Site

	// Ordered list of shortcodes for a page.
	shortcodes []*shortcode

	// All the shortcode names in this set.
	nameSet map[string]bool

	// Configuration
	enableInlineShortcodes bool
}

func newShortcodeHandler(p *pageState, s *Site, placeholderFunc func() string) *shortcodeHandler {

	sh := &shortcodeHandler{
		p:                      p,
		s:                      s,
		enableInlineShortcodes: s.enableInlineShortcodes,
		shortcodes:             make([]*shortcode, 0, 4),
		nameSet:                make(map[string]bool),
	}

	return sh
}

const (
	innerNewlineRegexp = "\n"
	innerCleanupRegexp = `\A<p>(.*)</p>\n\z`
	innerCleanupExpand = "$1"
)

func renderShortcode(
	level int,
	s *Site,
	tplVariants tpl.TemplateVariants,
	sc *shortcode,
	parent *ShortcodeWithPage,
	p *pageState) (string, bool, error) {

	var tmpl tpl.Template

	// Tracks whether this shortcode or any of its children has template variations
	// in other languages or output formats. We are currently only interested in
	// the output formats, so we may get some false positives -- we
	// should improve on that.
	var hasVariants bool

	if sc.isInline {
		if !p.s.enableInlineShortcodes {
			return "", false, nil
		}
		templName := path.Join("_inline_shortcode", p.File().Path(), sc.name)
		if sc.isClosing {
			templStr := sc.innerString()

			var err error
			tmpl, err = s.TextTmpl().Parse(templName, templStr)
			if err != nil {
				fe := herrors.ToFileError("html", err)
				l1, l2 := p.posOffset(sc.pos).LineNumber, fe.Position().LineNumber
				fe = herrors.ToFileErrorWithLineNumber(fe, l1+l2-1)
				return "", false, p.wrapError(fe)
			}

		} else {
			// Re-use of shortcode defined earlier in the same page.
			var found bool
			tmpl, found = s.TextTmpl().Lookup(templName)
			if !found {
				return "", false, _errors.Errorf("no earlier definition of shortcode %q found", sc.name)
			}
		}
	} else {
		var found, more bool
		tmpl, found, more = s.Tmpl().LookupVariant(sc.name, tplVariants)
		if !found {
			s.Log.ERROR.Printf("Unable to locate template for shortcode %q in page %q", sc.name, p.File().Path())
			return "", false, nil
		}
		hasVariants = hasVariants || more
	}

	data := &ShortcodeWithPage{Ordinal: sc.ordinal, posOffset: sc.pos, Params: sc.params, Page: newPageForShortcode(p), Parent: parent, Name: sc.name}
	if sc.params != nil {
		data.IsNamedParams = reflect.TypeOf(sc.params).Kind() == reflect.Map
	}

	if len(sc.inner) > 0 {
		var inner string
		for _, innerData := range sc.inner {
			switch innerData := innerData.(type) {
			case string:
				inner += innerData
			case *shortcode:
				s, more, err := renderShortcode(level+1, s, tplVariants, innerData, data, p)
				if err != nil {
					return "", false, err
				}
				hasVariants = hasVariants || more
				inner += s
			default:
				s.Log.ERROR.Printf("Illegal state on shortcode rendering of %q in page %q. Illegal type in inner data: %s ",
					sc.name, p.File().Path(), reflect.TypeOf(innerData))
				return "", false, nil
			}
		}

		// Pre Hugo 0.55 this was the behaviour even for the outer-most
		// shortcode.
		if sc.doMarkup && (level > 0 || sc.configVersion() == 1) {
			var err error
			b, err := p.pageOutput.cp.renderContent([]byte(inner), false)

			if err != nil {
				return "", false, err
			}

			newInner := b.Bytes()

			// If the type is “” (unknown) or “markdown”, we assume the markdown
			// generation has been performed. Given the input: `a line`, markdown
			// specifies the HTML `<p>a line</p>\n`. When dealing with documents as a
			// whole, this is OK. When dealing with an `{{ .Inner }}` block in Hugo,
			// this is not so good. This code does two things:
			//
			// 1.  Check to see if inner has a newline in it. If so, the Inner data is
			//     unchanged.
			// 2   If inner does not have a newline, strip the wrapping <p> block and
			//     the newline.
			switch p.m.markup {
			case "", "markdown":
				if match, _ := regexp.MatchString(innerNewlineRegexp, inner); !match {
					cleaner, err := regexp.Compile(innerCleanupRegexp)

					if err == nil {
						newInner = cleaner.ReplaceAll(newInner, []byte(innerCleanupExpand))
					}
				}
			}

			// TODO(bep) we may have plain text inner templates.
			data.Inner = template.HTML(newInner)
		} else {
			data.Inner = template.HTML(inner)
		}

	}

	result, err := renderShortcodeWithPage(s.Tmpl(), tmpl, data)

	if err != nil && sc.isInline {
		fe := herrors.ToFileError("html", err)
		l1, l2 := p.posFromPage(sc.pos).LineNumber, fe.Position().LineNumber
		fe = herrors.ToFileErrorWithLineNumber(fe, l1+l2-1)
		return "", false, fe
	}

	return result, hasVariants, err
}

func (s *shortcodeHandler) hasShortcodes() bool {
	return s != nil && len(s.shortcodes) > 0
}

func (s *shortcodeHandler) renderShortcodesForPage(p *pageState, f output.Format) (map[string]string, bool, error) {

	rendered := make(map[string]string)

	tplVariants := tpl.TemplateVariants{
		Language:     p.Language().Lang,
		OutputFormat: f,
	}

	var hasVariants bool

	for _, v := range s.shortcodes {
		s, more, err := renderShortcode(0, s.s, tplVariants, v, nil, p)
		if err != nil {
			err = p.parseError(_errors.Wrapf(err, "failed to render shortcode %q", v.name), p.source.parsed.Input(), v.pos)
			return nil, false, err
		}
		hasVariants = hasVariants || more
		rendered[v.placeholder] = s

	}

	return rendered, hasVariants, nil
}

var errShortCodeIllegalState = errors.New("Illegal shortcode state")

func (s *shortcodeHandler) parseError(err error, input []byte, pos int) error {
	if s.p != nil {
		return s.p.parseError(err, input, pos)
	}
	return err
}

// pageTokens state:
// - before: positioned just before the shortcode start
// - after: shortcode(s) consumed (plural when they are nested)
func (s *shortcodeHandler) extractShortcode(ordinal, level int, pt *pageparser.Iterator) (*shortcode, error) {
	if s == nil {
		panic("handler nil")
	}
	sc := &shortcode{ordinal: ordinal}

	var cnt = 0
	var nestedOrdinal = 0
	var nextLevel = level + 1

	fail := func(err error, i pageparser.Item) error {
		return s.parseError(err, pt.Input(), i.Pos)
	}

Loop:
	for {
		currItem := pt.Next()
		switch {
		case currItem.IsLeftShortcodeDelim():
			next := pt.Peek()
			if next.IsShortcodeClose() {
				continue
			}

			if cnt > 0 {
				// nested shortcode; append it to inner content
				pt.Backup()
				nested, err := s.extractShortcode(nestedOrdinal, nextLevel, pt)
				nestedOrdinal++
				if nested != nil && nested.name != "" {
					s.nameSet[nested.name] = true
				}

				if err == nil {
					sc.inner = append(sc.inner, nested)
				} else {
					return sc, err
				}

			} else {
				sc.doMarkup = currItem.IsShortcodeMarkupDelimiter()
			}

			cnt++

		case currItem.IsRightShortcodeDelim():
			// we trust the template on this:
			// if there's no inner, we're done
			if !sc.isInline && !sc.info.ParseInfo().IsInner {
				return sc, nil
			}

		case currItem.IsShortcodeClose():
			next := pt.Peek()
			if !sc.isInline && !sc.info.ParseInfo().IsInner {
				if next.IsError() {
					// return that error, more specific
					continue
				}

				return sc, fail(_errors.Errorf("shortcode %q has no .Inner, yet a closing tag was provided", next.Val), next)
			}
			if next.IsRightShortcodeDelim() {
				// self-closing
				pt.Consume(1)
			} else {
				sc.isClosing = true
				pt.Consume(2)
			}

			return sc, nil
		case currItem.IsText():
			sc.inner = append(sc.inner, currItem.ValStr())
		case currItem.Type == pageparser.TypeEmoji:
			// TODO(bep) avoid the duplication of these "text cases", to prevent
			// more of #6504 in the future.
			val := currItem.ValStr()
			if emoji := helpers.Emoji(val); emoji != nil {
				sc.inner = append(sc.inner, string(emoji))
			} else {
				sc.inner = append(sc.inner, val)
			}
		case currItem.IsShortcodeName():

			sc.name = currItem.ValStr()

			// Check if the template expects inner content.
			// We pick the first template for an arbitrary output format
			// if more than one. It is "all inner or no inner".
			tmpl, found, _ := s.s.Tmpl().LookupVariant(sc.name, tpl.TemplateVariants{})
			if !found {
				return nil, _errors.Errorf("template for shortcode %q not found", sc.name)
			}

			sc.info = tmpl.(tpl.Info)
		case currItem.IsInlineShortcodeName():
			sc.name = currItem.ValStr()
			sc.isInline = true
		case currItem.IsShortcodeParam():
			if !pt.IsValueNext() {
				continue
			} else if pt.Peek().IsShortcodeParamVal() {
				// named params
				if sc.params == nil {
					params := make(map[string]interface{})
					params[currItem.ValStr()] = pt.Next().ValTyped()
					sc.params = params
				} else {
					if params, ok := sc.params.(map[string]interface{}); ok {
						params[currItem.ValStr()] = pt.Next().ValTyped()
					} else {
						return sc, errShortCodeIllegalState
					}

				}
			} else {
				// positional params
				if sc.params == nil {
					var params []interface{}
					params = append(params, currItem.ValTyped())
					sc.params = params
				} else {
					if params, ok := sc.params.([]interface{}); ok {
						params = append(params, currItem.ValTyped())
						sc.params = params
					} else {
						return sc, errShortCodeIllegalState
					}

				}
			}
		case currItem.IsDone():
			// handled by caller
			pt.Backup()
			break Loop

		}
	}
	return sc, nil
}

// Replace prefixed shortcode tokens with the real content.
// Note: This function will rewrite the input slice.
func replaceShortcodeTokens(source []byte, replacements map[string]string) ([]byte, error) {

	if len(replacements) == 0 {
		return source, nil
	}

	start := 0

	pre := []byte(shortcodePlaceholderPrefix)
	post := []byte("HBHB")
	pStart := []byte("<p>")
	pEnd := []byte("</p>")

	k := bytes.Index(source[start:], pre)

	for k != -1 {
		j := start + k
		postIdx := bytes.Index(source[j:], post)
		if postIdx < 0 {
			// this should never happen, but let the caller decide to panic or not
			return nil, errors.New("illegal state in content; shortcode token missing end delim")
		}

		end := j + postIdx + 4

		newVal := []byte(replacements[string(source[j:end])])

		// Issue #1148: Check for wrapping p-tags <p>
		if j >= 3 && bytes.Equal(source[j-3:j], pStart) {
			if (k+4) < len(source) && bytes.Equal(source[end:end+4], pEnd) {
				j -= 3
				end += 4
			}
		}

		// This and other cool slice tricks: https://github.com/golang/go/wiki/SliceTricks
		source = append(source[:j], append(newVal, source[end:]...)...)
		start = j
		k = bytes.Index(source[start:], pre)

	}

	return source, nil
}

func renderShortcodeWithPage(h tpl.TemplateHandler, tmpl tpl.Template, data *ShortcodeWithPage) (string, error) {
	buffer := bp.GetBuffer()
	defer bp.PutBuffer(buffer)

	err := h.Execute(tmpl, buffer, data)
	if err != nil {
		return "", _errors.Wrap(err, "failed to process shortcode")
	}
	return buffer.String(), nil
}
