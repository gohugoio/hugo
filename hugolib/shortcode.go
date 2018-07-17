// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/output"

	"github.com/gohugoio/hugo/media"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
)

// ShortcodeWithPage is the "." context in a shortcode template.
type ShortcodeWithPage struct {
	Params        interface{}
	Inner         template.HTML
	Page          *PageWithoutContent
	Parent        *ShortcodeWithPage
	IsNamedParams bool

	// Zero-based ordinal in relation to its parent. If the parent is the page itself,
	// this ordinal will represent the position of this shortcode in the page content.
	Ordinal int

	scratch *maps.Scratch
}

// Site returns information about the current site.
func (scp *ShortcodeWithPage) Site() *SiteInfo {
	return scp.Page.Site
}

// Ref is a shortcut to the Ref method on Page.
func (scp *ShortcodeWithPage) Ref(args map[string]interface{}) (string, error) {
	return scp.Page.Ref(args)
}

// RelRef is a shortcut to the RelRef method on Page.
func (scp *ShortcodeWithPage) RelRef(args map[string]interface{}) (string, error) {
	return scp.Page.RelRef(args)
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

	switch x.Kind() {
	case reflect.String:
		return x.String()
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		return x.Int()
	default:
		return x
	}

}

// Note - this value must not contain any markup syntax
const shortcodePlaceholderPrefix = "HUGOSHORTCODE"

type shortcode struct {
	name     string
	inner    []interface{} // string or nested shortcode
	params   interface{}   // map or array
	ordinal  int
	err      error
	doMarkup bool
}

func (sc shortcode) String() string {
	// for testing (mostly), so any change here will break tests!
	var params interface{}
	switch v := sc.params.(type) {
	case map[string]string:
		// sort the keys so test assertions won't fail
		var keys []string
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var tmp = make([]string, len(keys))

		for i, k := range keys {
			tmp[i] = k + ":" + v[k]
		}
		params = tmp

	default:
		// use it as is
		params = sc.params
	}

	return fmt.Sprintf("%s(%q, %t){%s}", sc.name, params, sc.doMarkup, sc.inner)
}

// We may have special shortcode templates for AMP etc.
// Note that in the below, OutputFormat may be empty.
// We will try to look for the most specific shortcode template available.
type scKey struct {
	Lang                 string
	OutputFormat         string
	Suffix               string
	ShortcodePlaceholder string
}

func newScKey(m media.Type, shortcodeplaceholder string) scKey {
	return scKey{Suffix: m.Suffix(), ShortcodePlaceholder: shortcodeplaceholder}
}

func newScKeyFromLangAndOutputFormat(lang string, o output.Format, shortcodeplaceholder string) scKey {
	return scKey{Lang: lang, Suffix: o.MediaType.Suffix(), OutputFormat: o.Name, ShortcodePlaceholder: shortcodeplaceholder}
}

func newDefaultScKey(shortcodeplaceholder string) scKey {
	return newScKey(media.HTMLType, shortcodeplaceholder)
}

type shortcodeHandler struct {
	init sync.Once

	p *PageWithoutContent

	// This is all shortcode rendering funcs for all potential output formats.
	contentShortcodes *orderedMap

	// This map contains the new or changed set of shortcodes that need
	// to be rendered for the current output format.
	contentShortcodesDelta *orderedMap

	// This maps the shorcode placeholders with the rendered content.
	// We will do (potential) partial re-rendering per output format,
	// so keep this for the unchanged.
	renderedShortcodes map[string]string

	// Maps the shortcodeplaceholder with the actual shortcode.
	shortcodes *orderedMap

	// All the shortcode names in this set.
	nameSet map[string]bool

	placeholderID   int
	placeholderFunc func() string
}

func (s *shortcodeHandler) nextPlaceholderID() int {
	s.placeholderID++
	return s.placeholderID
}

func (s *shortcodeHandler) createShortcodePlaceholder() string {
	if s.placeholderFunc != nil {
		return s.placeholderFunc()
	}
	return fmt.Sprintf("HAHA%s-%p-%d-HBHB", shortcodePlaceholderPrefix, s.p.Page, s.nextPlaceholderID())
}

func newShortcodeHandler(p *Page) *shortcodeHandler {
	return &shortcodeHandler{
		p:                  p.withoutContent(),
		contentShortcodes:  newOrderedMap(),
		shortcodes:         newOrderedMap(),
		nameSet:            make(map[string]bool),
		renderedShortcodes: make(map[string]string),
	}
}

// TODO(bep) make it non-global
var isInnerShortcodeCache = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

// to avoid potential costly look-aheads for closing tags we look inside the template itself
// we could change the syntax to self-closing tags, but that would make users cry
// the value found is cached
func isInnerShortcode(t tpl.TemplateExecutor) (bool, error) {
	isInnerShortcodeCache.RLock()
	m, ok := isInnerShortcodeCache.m[t.Name()]
	isInnerShortcodeCache.RUnlock()

	if ok {
		return m, nil
	}

	isInnerShortcodeCache.Lock()
	defer isInnerShortcodeCache.Unlock()
	match, _ := regexp.MatchString("{{.*?\\.Inner.*?}}", t.Tree())
	isInnerShortcodeCache.m[t.Name()] = match

	return match, nil
}

func clearIsInnerShortcodeCache() {
	isInnerShortcodeCache.Lock()
	defer isInnerShortcodeCache.Unlock()
	isInnerShortcodeCache.m = make(map[string]bool)
}

const innerNewlineRegexp = "\n"
const innerCleanupRegexp = `\A<p>(.*)</p>\n\z`
const innerCleanupExpand = "$1"

func prepareShortcodeForPage(placeholder string, sc *shortcode, parent *ShortcodeWithPage, p *PageWithoutContent) map[scKey]func() (string, error) {

	m := make(map[scKey]func() (string, error))
	lang := p.Lang()

	for _, f := range p.outputFormats {
		// The most specific template will win.
		key := newScKeyFromLangAndOutputFormat(lang, f, placeholder)
		m[key] = func() (string, error) {
			return renderShortcode(key, sc, nil, p), nil
		}
	}

	return m
}

func renderShortcode(
	tmplKey scKey,
	sc *shortcode,
	parent *ShortcodeWithPage,
	p *PageWithoutContent) string {

	tmpl := getShortcodeTemplateForTemplateKey(tmplKey, sc.name, p.s.Tmpl)
	if tmpl == nil {
		p.s.Log.ERROR.Printf("Unable to locate template for shortcode %q in page %q", sc.name, p.Path())
		return ""
	}

	data := &ShortcodeWithPage{Ordinal: sc.ordinal, Params: sc.params, Page: p, Parent: parent}
	if sc.params != nil {
		data.IsNamedParams = reflect.TypeOf(sc.params).Kind() == reflect.Map
	}

	if len(sc.inner) > 0 {
		var inner string
		for _, innerData := range sc.inner {
			switch innerData.(type) {
			case string:
				inner += innerData.(string)
			case *shortcode:
				inner += renderShortcode(tmplKey, innerData.(*shortcode), data, p)
			default:
				p.s.Log.ERROR.Printf("Illegal state on shortcode rendering of %q in page %q. Illegal type in inner data: %s ",
					sc.name, p.Path(), reflect.TypeOf(innerData))
				return ""
			}
		}

		if sc.doMarkup {
			newInner := p.s.ContentSpec.RenderBytes(&helpers.RenderingContext{
				Content:      []byte(inner),
				PageFmt:      p.Markup,
				Cfg:          p.Language(),
				DocumentID:   p.UniqueID(),
				DocumentName: p.Path(),
				Config:       p.getRenderingConfig()})

			// If the type is “unknown” or “markdown”, we assume the markdown
			// generation has been performed. Given the input: `a line`, markdown
			// specifies the HTML `<p>a line</p>\n`. When dealing with documents as a
			// whole, this is OK. When dealing with an `{{ .Inner }}` block in Hugo,
			// this is not so good. This code does two things:
			//
			// 1.  Check to see if inner has a newline in it. If so, the Inner data is
			//     unchanged.
			// 2   If inner does not have a newline, strip the wrapping <p> block and
			//     the newline. This was previously tricked out by wrapping shortcode
			//     substitutions in <div>HUGOSHORTCODE-1</div> which prevents the
			//     generation, but means that you can’t use shortcodes inside of
			//     markdown structures itself (e.g., `[foo]({{% ref foo.md %}})`).
			switch p.Markup {
			case "unknown", "markdown":
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

	return renderShortcodeWithPage(tmpl, data)
}

// The delta represents new output format-versions of the shortcodes,
// which, combined with the ones that do not have alternative representations,
// builds a complete set ready for a full rebuild of the Page content.
// This method returns false if there are no new shortcode variants in the
// current rendering context's output format. This mean we can safely reuse
// the content from the previous output format, if any.
func (s *shortcodeHandler) updateDelta() bool {
	s.init.Do(func() {
		s.contentShortcodes = createShortcodeRenderers(s.shortcodes, s.p.withoutContent())
	})

	if !s.p.shouldRenderTo(s.p.s.rc.Format) {
		// TODO(bep) add test for this re translations
		return false
	}
	of := s.p.s.rc.Format
	contentShortcodes := s.contentShortcodesForOutputFormat(of)

	if s.contentShortcodesDelta == nil || s.contentShortcodesDelta.Len() == 0 {
		s.contentShortcodesDelta = contentShortcodes
		return true
	}

	delta := newOrderedMap()

	for _, k := range contentShortcodes.Keys() {
		if !s.contentShortcodesDelta.Contains(k) {
			v, _ := contentShortcodes.Get(k)
			delta.Add(k, v)
		}
	}

	s.contentShortcodesDelta = delta

	return delta.Len() > 0
}

func (s *shortcodeHandler) clearDelta() {
	if s == nil {
		return
	}
	s.contentShortcodesDelta = newOrderedMap()
}

func (s *shortcodeHandler) contentShortcodesForOutputFormat(f output.Format) *orderedMap {
	contentShortcodesForOuputFormat := newOrderedMap()
	lang := s.p.Lang()

	for _, key := range s.shortcodes.Keys() {
		shortcodePlaceholder := key.(string)

		key := newScKeyFromLangAndOutputFormat(lang, f, shortcodePlaceholder)
		renderFn, found := s.contentShortcodes.Get(key)

		if !found {
			key.OutputFormat = ""
			renderFn, found = s.contentShortcodes.Get(key)
		}

		// Fall back to HTML
		if !found && key.Suffix != "html" {
			key.Suffix = "html"
			renderFn, found = s.contentShortcodes.Get(key)
		}

		if !found {
			panic(fmt.Sprintf("Shortcode %q could not be found", shortcodePlaceholder))
		}
		contentShortcodesForOuputFormat.Add(newScKeyFromLangAndOutputFormat(lang, f, shortcodePlaceholder), renderFn)
	}

	return contentShortcodesForOuputFormat
}

func (s *shortcodeHandler) executeShortcodesForDelta(p *PageWithoutContent) error {

	for _, k := range s.contentShortcodesDelta.Keys() {
		render := s.contentShortcodesDelta.getShortcodeRenderer(k)
		renderedShortcode, err := render()
		if err != nil {
			return fmt.Errorf("Failed to execute shortcode in page %q: %s", p.Path(), err)
		}

		s.renderedShortcodes[k.(scKey).ShortcodePlaceholder] = renderedShortcode
	}

	return nil

}

func createShortcodeRenderers(shortcodes *orderedMap, p *PageWithoutContent) *orderedMap {

	shortcodeRenderers := newOrderedMap()

	for _, k := range shortcodes.Keys() {
		v := shortcodes.getShortcode(k)
		prepared := prepareShortcodeForPage(k.(string), v, nil, p)
		for kk, vv := range prepared {
			shortcodeRenderers.Add(kk, vv)
		}
	}

	return shortcodeRenderers
}

var errShortCodeIllegalState = errors.New("Illegal shortcode state")

// pageTokens state:
// - before: positioned just before the shortcode start
// - after: shortcode(s) consumed (plural when they are nested)
func (s *shortcodeHandler) extractShortcode(ordinal int, pt *pageTokens, p *PageWithoutContent) (*shortcode, error) {
	sc := &shortcode{ordinal: ordinal}
	var isInner = false

	var currItem item
	var cnt = 0
	var nestedOrdinal = 0

Loop:
	for {
		currItem = pt.next()

		switch currItem.typ {
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup:
			next := pt.peek()
			if next.typ == tScClose {
				continue
			}

			if cnt > 0 {
				// nested shortcode; append it to inner content
				pt.backup3(currItem, next)
				nested, err := s.extractShortcode(nestedOrdinal, pt, p)
				nestedOrdinal++
				if nested.name != "" {
					s.nameSet[nested.name] = true
				}
				if err == nil {
					sc.inner = append(sc.inner, nested)
				} else {
					return sc, err
				}

			} else {
				sc.doMarkup = currItem.typ == tLeftDelimScWithMarkup
			}

			cnt++

		case tRightDelimScWithMarkup, tRightDelimScNoMarkup:
			// we trust the template on this:
			// if there's no inner, we're done
			if !isInner {
				return sc, nil
			}

		case tScClose:
			next := pt.peek()
			if !isInner {
				if next.typ == tError {
					// return that error, more specific
					continue
				}
				return sc, fmt.Errorf("Shortcode '%s' in page '%s' has no .Inner, yet a closing tag was provided", next.val, p.FullFilePath())
			}
			if next.typ == tRightDelimScWithMarkup || next.typ == tRightDelimScNoMarkup {
				// self-closing
				pt.consume(1)
			} else {
				pt.consume(2)
			}

			return sc, nil
		case tText:
			sc.inner = append(sc.inner, currItem.val)
		case tScName:
			sc.name = currItem.val
			// We pick the first template for an arbitrary output format
			// if more than one. It is "all inner or no inner".
			tmpl := getShortcodeTemplateForTemplateKey(scKey{}, sc.name, p.s.Tmpl)
			if tmpl == nil {
				return sc, fmt.Errorf("Unable to locate template for shortcode %q in page %q", sc.name, p.Path())
			}

			var err error
			isInner, err = isInnerShortcode(tmpl.(tpl.TemplateExecutor))
			if err != nil {
				return sc, fmt.Errorf("Failed to handle template for shortcode %q for page %q: %s", sc.name, p.Path(), err)
			}

		case tScParam:
			if !pt.isValueNext() {
				continue
			} else if pt.peek().typ == tScParamVal {
				// named params
				if sc.params == nil {
					params := make(map[string]string)
					params[currItem.val] = pt.next().val
					sc.params = params
				} else {
					if params, ok := sc.params.(map[string]string); ok {
						params[currItem.val] = pt.next().val
					} else {
						return sc, errShortCodeIllegalState
					}

				}
			} else {
				// positional params
				if sc.params == nil {
					var params []string
					params = append(params, currItem.val)
					sc.params = params
				} else {
					if params, ok := sc.params.([]string); ok {
						params = append(params, currItem.val)
						sc.params = params
					} else {
						return sc, errShortCodeIllegalState
					}

				}
			}

		case tError, tEOF:
			// handled by caller
			pt.backup()
			break Loop

		}
	}
	return sc, nil
}

func (s *shortcodeHandler) extractShortcodes(stringToParse string, p *PageWithoutContent) (string, error) {

	startIdx := strings.Index(stringToParse, "{{")

	// short cut for docs with no shortcodes
	if startIdx < 0 {
		return stringToParse, nil
	}

	// the parser takes a string;
	// since this is an internal API, it could make sense to use the mutable []byte all the way, but
	// it seems that the time isn't really spent in the byte copy operations, and the impl. gets a lot cleaner
	pt := &pageTokens{lexer: newShortcodeLexer("parse-page", stringToParse, pos(startIdx))}

	result := bp.GetBuffer()
	defer bp.PutBuffer(result)
	//var result bytes.Buffer

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currItem item
	var currShortcode shortcode
	var ordinal int

Loop:
	for {
		currItem = pt.next()

		switch currItem.typ {
		case tText:
			result.WriteString(currItem.val)
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup:
			// let extractShortcode handle left delim (will do so recursively)
			pt.backup()

			currShortcode, err := s.extractShortcode(ordinal, pt, p)

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if err != nil {
				return result.String(), err
			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			placeHolder := s.createShortcodePlaceholder()
			result.WriteString(placeHolder)
			ordinal++
			s.shortcodes.Add(placeHolder, currShortcode)
		case tEOF:
			break Loop
		case tError:
			err := fmt.Errorf("%s:%d: %s",
				p.FullFilePath(), (p.lineNumRawContentStart() + pt.lexer.lineNum() - 1), currItem)
			currShortcode.err = err
			return result.String(), err
		}
	}

	return result.String(), nil

}

// Replace prefixed shortcode tokens (HUGOSHORTCODE-1, HUGOSHORTCODE-2) with the real content.
// Note: This function will rewrite the input slice.
func replaceShortcodeTokens(source []byte, prefix string, replacements map[string]string) ([]byte, error) {

	if len(replacements) == 0 {
		return source, nil
	}

	sourceLen := len(source)
	start := 0

	pre := []byte("HAHA" + prefix)
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
			if (k+4) < sourceLen && bytes.Equal(source[end:end+4], pEnd) {
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

func getShortcodeTemplateForTemplateKey(key scKey, shortcodeName string, t tpl.TemplateFinder) tpl.Template {
	isInnerShortcodeCache.RLock()
	defer isInnerShortcodeCache.RUnlock()

	var names []string

	suffix := strings.ToLower(key.Suffix)
	outFormat := strings.ToLower(key.OutputFormat)
	lang := strings.ToLower(key.Lang)

	if outFormat != "" && suffix != "" {
		if lang != "" {
			names = append(names, fmt.Sprintf("%s.%s.%s.%s", shortcodeName, lang, outFormat, suffix))
		}
		names = append(names, fmt.Sprintf("%s.%s.%s", shortcodeName, outFormat, suffix))
	}

	if suffix != "" {
		if lang != "" {
			names = append(names, fmt.Sprintf("%s.%s.%s", shortcodeName, lang, suffix))
		}
		names = append(names, fmt.Sprintf("%s.%s", shortcodeName, suffix))
	}

	names = append(names, shortcodeName)

	for _, name := range names {

		if x, found := t.Lookup("shortcodes/" + name); found {
			return x
		}
		if x, found := t.Lookup("theme/shortcodes/" + name); found {
			return x
		}
		if x, found := t.Lookup("_internal/shortcodes/" + name); found {
			return x
		}
	}
	return nil
}

func renderShortcodeWithPage(tmpl tpl.Template, data *ShortcodeWithPage) string {
	buffer := bp.GetBuffer()
	defer bp.PutBuffer(buffer)

	isInnerShortcodeCache.RLock()
	err := tmpl.Execute(buffer, data)
	isInnerShortcodeCache.RUnlock()
	if err != nil {
		data.Page.s.Log.ERROR.Printf("error processing shortcode %q for page %q: %s", tmpl.Name(), data.Page.Path(), err)
	}
	return buffer.String()
}
