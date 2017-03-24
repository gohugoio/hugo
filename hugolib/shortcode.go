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
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"

	bp "github.com/spf13/hugo/bufferpool"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/tpl"
)

// ShortcodeWithPage is the "." context in a shortcode template.
type ShortcodeWithPage struct {
	Params        interface{}
	Inner         template.HTML
	InnerPage     template.HTML
	Page          *Page
	Parent        *ShortcodeWithPage
	IsNamedParams bool
	scratch       *Scratch
}

// Site returns information about the current site.
func (scp *ShortcodeWithPage) Site() *SiteInfo {
	return scp.Page.Site
}

// Ref is a shortcut to the Ref method on Page.
func (scp *ShortcodeWithPage) Ref(ref string) (string, error) {
	return scp.Page.Ref(ref)
}

// RelRef is a shortcut to the RelRef method on Page.
func (scp *ShortcodeWithPage) RelRef(ref string) (string, error) {
	return scp.Page.RelRef(ref)
}

// Scratch returns a scratch-pad scoped for this shortcode. This can be used
// as a temporary storage for variables, counters etc.
func (scp *ShortcodeWithPage) Scratch() *Scratch {
	if scp.scratch == nil {
		scp.scratch = newScratch()
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
			return "error: cannot access named params by position"
		} else if reflect.TypeOf(scp.Params).Kind() == reflect.Slice {
			idx := int(reflect.ValueOf(key).Int())
			ln := reflect.ValueOf(scp.Params).Len()
			if idx > ln-1 {
				helpers.DistinctErrorLog.Printf("No shortcode param at .Get %d in page %s, have params: %v", idx, scp.Page.FullFilePath(), scp.Params)
				return fmt.Sprintf("error: index out of range for positional param at position %d", idx)
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
			if reflect.ValueOf(scp.Params).Len() == 1 && reflect.ValueOf(scp.Params).Index(0).String() == "" {
				return nil
			}
			return "error: cannot access positional params by string name"
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
	name      string
	inner     []interface{} // string or nested shortcode
	params    interface{}   // map or array
	err       error
	doMarkup  bool
	reusePage bool
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

type shortcodeHandler struct {
	// Maps the shortcodeplaceholder with the shortcode rendering func.
	contentShortCodes map[string]func() (string, error)

	// Maps the shortcodeplaceholder with the actual shortcode.
	shortcodes map[string]shortcode

	// All the shortcode names in this set.
	nameSet map[string]bool
}

func newShortcodeHandler() *shortcodeHandler {
	return &shortcodeHandler{
		contentShortCodes: make(map[string]func() (string, error)),
		shortcodes:        make(map[string]shortcode),
		nameSet:           make(map[string]bool),
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
func isInnerShortcode(t *template.Template) (bool, error) {
	isInnerShortcodeCache.RLock()
	m, ok := isInnerShortcodeCache.m[t.Name()]
	isInnerShortcodeCache.RUnlock()

	if ok {
		return m, nil
	}

	isInnerShortcodeCache.Lock()
	defer isInnerShortcodeCache.Unlock()
	if t.Tree == nil {
		return false, errors.New("Template failed to compile")
	}
	match, _ := regexp.MatchString("{{.*?\\.Inner.*?}}", t.Tree.Root.String())
	isInnerShortcodeCache.m[t.Name()] = match

	return match, nil
}

func clearIsInnerShortcodeCache() {
	isInnerShortcodeCache.Lock()
	defer isInnerShortcodeCache.Unlock()
	isInnerShortcodeCache.m = make(map[string]bool)
}

func createShortcodePlaceholder(id int) string {
	return fmt.Sprintf("HAHA%s-%dHBHB", shortcodePlaceholderPrefix, id)
}

// HandleAndReplaceShortcodes does all in one go: extract, render and replace
// used for Page Reuse shortcodes
func (page *Page) HandleAndReplaceShortcodes() (string, error) {
	stringToParse := string(page.rawContent)
	shortcodeState := newShortcodeHandler()
	tmpContent, err := shortcodeState.extractAndRenderShortcodes(stringToParse, page)
	if err != nil {
		return "", err
	}
	if len(shortcodeState.contentShortCodes) > 0 {
		shortcodes, err := executeShortcodeFuncMap(shortcodeState.contentShortCodes)
		if err != nil {
			return "", err
		}
		tmpContentWithTokensReplaced, err := replaceShortcodeTokens([]byte(tmpContent), shortcodePlaceholderPrefix, shortcodes)
		if err != nil {
			return "", fmt.Errorf("Failed to replace shortcode tokens in %s:\n%s", page.BaseFileName(), err.Error())
		}
		return string(tmpContentWithTokensReplaced), nil
	}
	return tmpContent, nil
}

const innerNewlineRegexp = "\n"
const innerCleanupRegexp = `\A<p>(.*)</p>\n\z`
const innerCleanupExpand = "$1"

func renderShortcode(sc shortcode, parent *ShortcodeWithPage, p *Page) string {
	tmpl := getShortcodeTemplate(sc.name, p.s.Tmpl)

	if tmpl == nil && !sc.reusePage {
		p.s.Log.ERROR.Printf("Unable to locate template for shortcode '%s' in page %q", sc.name, p.Path())
		return ""
	}

	data := &ShortcodeWithPage{Params: sc.params, Page: p, Parent: parent}
	var reusedPageRef string
	if sc.params != nil {
		data.IsNamedParams = reflect.TypeOf(sc.params).Kind() == reflect.Map
		if sc.reusePage {
			if data.IsNamedParams {
				reusedPageRef = sc.params.(map[string]string)["ref"]
			} else {
				paramsArr := sc.params.([]string)
				if len(paramsArr) > 0 {
					reusedPageRef = paramsArr[0]
				}
			}
		}
	}
	if sc.reusePage && reusedPageRef == "" {
		p.s.Log.ERROR.Printf("No page ref specified for Reuse shortcode '%s' in page %s",
			sc.name, p.BaseFileName())
		return ""
	}
	var reuseContent []byte
	if len(sc.inner) > 0 || sc.reusePage {
		var inner string
		for _, innerData := range sc.inner {
			switch innerData.(type) {
			case string:
				inner += innerData.(string)
			case shortcode:
				inner += renderShortcode(innerData.(shortcode), data, p)
			default:
				p.s.Log.ERROR.Printf("Illegal state on shortcode rendering of %q in page %q. Illegal type in inner data: %s ",
					sc.name, p.Path(), reflect.TypeOf(innerData))
				return ""
			}
		}
		if sc.reusePage {
			reusedPage, err := p.Site.LookupPage(reusedPageRef, p)
			if err != nil {
				p.s.Log.ERROR.Printf("Unable to locate a page for Reuse shortcode '%s' in page %s: %s",
					sc.name, p.BaseFileName(), err)
				return ""
			}
			if reusedPage == nil {
				p.s.Log.ERROR.Printf("Unable to locate a page for Reuse shortcode '%s' in page %s",
					sc.name, p.BaseFileName())
				return ""
			}
			reusedPageWithShortcodes, err := reusedPage.HandleAndReplaceShortcodes()
			if err != nil {
				p.s.Log.ERROR.Printf("Unable to handle shortcodes in the reused page '%s'",
					reusedPage.BaseFileName())
				return ""
			}
			reuseContent = p.renderReusedPage([]byte(reusedPageWithShortcodes), reusedPage)
		}

		var newInner []byte
		if sc.doMarkup {
			newInner = p.s.ContentSpec.RenderBytes(&helpers.RenderingContext{
				Content: []byte(inner), PageFmt: p.determineMarkupType(),
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
			switch p.determineMarkupType() {
			case "unknown", "markdown":
				if match, _ := regexp.MatchString(innerNewlineRegexp, inner); !match {
					cleaner, err := regexp.Compile(innerCleanupRegexp)
					if err == nil {
						newInner = cleaner.ReplaceAll(newInner, []byte(innerCleanupExpand))
					}
				}
			}
		} else {
			newInner = []byte(inner)
		}
		data.Inner = template.HTML(newInner)
		if sc.reusePage {
			data.InnerPage = template.HTML(reuseContent)
		}
	}

	if tmpl != nil {
		return renderShortcodeWithPage(tmpl, data)
	}
	return string(reuseContent)
}

func (s *shortcodeHandler) extractAndRenderShortcodes(stringToParse string, p *Page) (string, error) {
	content, err := s.extractShortcodes(stringToParse, p)

	if err != nil {
		//  try to render what we have whilst logging the error
		p.s.Log.ERROR.Println(err.Error())
	}

	s.contentShortCodes = renderShortcodes(s.shortcodes, p)

	return content, err

}

var emptyShortcodeFn = func() (string, error) { return "", nil }

func executeShortcodeFuncMap(funcs map[string]func() (string, error)) (map[string]string, error) {
	result := make(map[string]string)

	for k, v := range funcs {
		s, err := v()
		if err != nil {
			return nil, fmt.Errorf("Failed to execute shortcode with key %s: %s", k, err)
		}
		result[k] = s
	}

	return result, nil
}

func renderShortcodes(shortcodes map[string]shortcode, p *Page) map[string]func() (string, error) {
	renderedShortcodes := make(map[string]func() (string, error))

	for key, sc := range shortcodes {
		if sc.err != nil {
			// need to have something to replace with
			renderedShortcodes[key] = emptyShortcodeFn
		} else {
			shorctode := sc
			renderedShortcodes[key] = func() (string, error) { return renderShortcode(shorctode, nil, p), nil }
		}
	}

	return renderedShortcodes
}

var errShortCodeIllegalState = errors.New("Illegal shortcode state")

// pageTokens state:
// - before: positioned just before the shortcode start
// - after: shortcode(s) consumed (plural when they are nested)
func (s *shortcodeHandler) extractShortcode(pt *pageTokens, p *Page) (shortcode, error) {
	sc := shortcode{}
	var isInner = false

	var currItem item
	var cnt = 0

Loop:
	for {
		currItem = pt.next()

		switch currItem.typ {
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup, tLeftDelimScReuse:
			next := pt.peek()
			if next.typ == tScClose {
				continue
			}

			if cnt > 0 {
				// nested shortcode; append it to inner content
				pt.backup3(currItem, next)
				nested, err := s.extractShortcode(pt, p)
				if nested.name != "" {
					s.nameSet[nested.name] = true
				}
				if err == nil {
					sc.inner = append(sc.inner, nested)
				} else {
					return sc, err
				}

			} else {
				sc.doMarkup = currItem.typ == tLeftDelimScWithMarkup || currItem.typ == tLeftDelimScReuse
				sc.reusePage = currItem.typ == tLeftDelimScReuse
			}

			cnt++

		case tRightDelimScWithMarkup, tRightDelimScNoMarkup, tRightDelimScReuse:
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
			if next.typ == tRightDelimScWithMarkup || next.typ == tRightDelimScNoMarkup || next.typ == tRightDelimScReuse {
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
			tmpl := getShortcodeTemplate(sc.name, p.s.Tmpl)
			{
			}
			if tmpl == nil {
				if sc.reusePage {
					continue
				}
				return sc, fmt.Errorf("Unable to locate template for shortcode %q in page %q", sc.name, p.Path())
			}

			var err error
			isInner, err = isInnerShortcode(tmpl)
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

func (s *shortcodeHandler) extractShortcodes(stringToParse string, p *Page) (string, error) {

	startIdx := strings.Index(stringToParse, "{{")

	// short cut for docs with no shortcodes
	if startIdx < 0 {
		return stringToParse, nil
	}

	// the parser takes a string;
	// since this is an internal API, it could make sense to use the mutable []byte all the way, but
	// it seems that the time isn't really spent in the byte copy operations, and the impl. gets a lot cleaner
	pt := &pageTokens{lexer: newShortcodeLexer("parse-page", stringToParse, pos(startIdx))}

	id := 1 // incremented id, will be appended onto temp. shortcode placeholders

	result := bp.GetBuffer()
	defer bp.PutBuffer(result)
	//var result bytes.Buffer

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currItem item
	var currShortcode shortcode

Loop:
	for {
		currItem = pt.next()

		switch currItem.typ {
		case tText:
			result.WriteString(currItem.val)
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup, tLeftDelimScReuse:
			// let extractShortcode handle left delim (will do so recursively)
			pt.backup()

			currShortcode, err := s.extractShortcode(pt, p)

			if currShortcode.name != "" {
				s.nameSet[currShortcode.name] = true
			}

			if err != nil {
				return result.String(), err
			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			placeHolder := createShortcodePlaceholder(id)
			result.WriteString(placeHolder)
			s.shortcodes[placeHolder] = currShortcode
			id++
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

func getShortcodeTemplate(name string, t tpl.Template) *template.Template {
	if x := t.Lookup("shortcodes/" + name + ".html"); x != nil {
		return x
	}
	if x := t.Lookup("theme/shortcodes/" + name + ".html"); x != nil {
		return x
	}
	return t.Lookup("_internal/shortcodes/" + name + ".html")
}

func renderShortcodeWithPage(tmpl *template.Template, data *ShortcodeWithPage) string {
	buffer := bp.GetBuffer()
	defer bp.PutBuffer(buffer)

	isInnerShortcodeCache.RLock()
	err := tmpl.Execute(buffer, data)
	isInnerShortcodeCache.RUnlock()
	if err != nil {
		data.Page.s.Log.ERROR.Println("error processing shortcode", tmpl.Name(), "\n ERR:", err)
		data.Page.s.Log.WARN.Println(data)
	}
	return buffer.String()
}
