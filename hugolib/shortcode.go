// Copyright © 2013-14 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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
	"html/template"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/tpl"
	jww "github.com/spf13/jwalterweatherman"
)

type ShortcodeFunc func([]string) string

type Shortcode struct {
	Name string
	Func ShortcodeFunc
}

type ShortcodeWithPage struct {
	Params interface{}
	Inner  template.HTML
	Page   *Page
}

func (scp *ShortcodeWithPage) Ref(ref string) (string, error) {
	return scp.Page.Ref(ref)
}

func (scp *ShortcodeWithPage) RelRef(ref string) (string, error) {
	return scp.Page.RelRef(ref)
}

func (scp *ShortcodeWithPage) Get(key interface{}) interface{} {
	if reflect.ValueOf(scp.Params).Len() == 0 {
		return nil
	}

	var x reflect.Value

	switch key.(type) {
	case int64, int32, int16, int8, int:
		if reflect.TypeOf(scp.Params).Kind() == reflect.Map {
			return "error: cannot access named params by position"
		} else if reflect.TypeOf(scp.Params).Kind() == reflect.Slice {
			x = reflect.ValueOf(scp.Params).Index(int(reflect.ValueOf(key).Int()))
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
	name     string
	inner    []interface{} // string or nested shortcode
	params   interface{}   // map or array
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

// all in  one go: extract, render and replace
// only used for testing
func ShortcodesHandle(stringToParse string, page *Page, t tpl.Template) string {
	tmpContent, tmpShortcodes := extractAndRenderShortcodes(stringToParse, page, t)

	if len(tmpShortcodes) > 0 {
		tmpContentWithTokensReplaced, err := replaceShortcodeTokens([]byte(tmpContent), shortcodePlaceholderPrefix, -1, true, tmpShortcodes)

		if err != nil {
			jww.ERROR.Printf("Fail to replace short code tokens in %s:\n%s", page.BaseFileName(), err.Error())
		} else {
			return string(tmpContentWithTokensReplaced)
		}
	}

	return string(tmpContent)
}

var isInnerShortcodeCache = make(map[string]bool)

// to avoid potential costly look-aheads for closing tags we look inside the template itself
// we could change the syntax to self-closing tags, but that would make users cry
// the value found is cached
func isInnerShortcode(t *template.Template) bool {
	if m, ok := isInnerShortcodeCache[t.Name()]; ok {
		return m
	}

	match, _ := regexp.MatchString("{{.*?\\.Inner.*?}}", t.Tree.Root.String())
	isInnerShortcodeCache[t.Name()] = match

	return match
}

func createShortcodePlaceholder(id int) string {
	return fmt.Sprintf("{@{@%s-%d@}@}", shortcodePlaceholderPrefix, id)
}

const innerNewlineRegexp = "\n"
const innerCleanupRegexp = `\A<p>(.*)</p>\n\z`
const innerCleanupExpand = "$1"

func renderShortcode(sc shortcode, p *Page, t tpl.Template) string {
	var data = &ShortcodeWithPage{Params: sc.params, Page: p}
	tmpl := GetTemplate(sc.name, t)

	if tmpl == nil {
		jww.ERROR.Printf("Unable to locate template for shortcode '%s' in page %s", sc.name, p.BaseFileName())
		return ""
	}

	if len(sc.inner) > 0 {
		var inner string
		for _, innerData := range sc.inner {
			switch innerData.(type) {
			case string:
				inner += innerData.(string)
			case shortcode:
				inner += renderShortcode(innerData.(shortcode), p, t)
			default:
				jww.ERROR.Printf("Illegal state on shortcode rendering of '%s' in page %s. Illegal type in inner data: %s ",
					sc.name, p.BaseFileName(), reflect.TypeOf(innerData))
				return ""
			}
		}

		if sc.doMarkup {
			newInner := helpers.RenderBytes(helpers.RenderingContext{
				Content: []byte(inner), PageFmt: p.guessMarkupType(),
				DocumentId: p.UniqueId(), ConfigFlags: p.getRenderingConfigFlags()})

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
			switch p.guessMarkupType() {
			case "unknown", "markdown":
				if match, _ := regexp.MatchString(innerNewlineRegexp, inner); !match {
					cleaner, err := regexp.Compile(innerCleanupRegexp)

					if err == nil {
						newInner = cleaner.ReplaceAll(newInner, []byte(innerCleanupExpand))
					}
				}
			}

			data.Inner = template.HTML(newInner)
		} else {
			data.Inner = template.HTML(inner)
		}

	}

	return ShortcodeRender(tmpl, data)
}

func extractAndRenderShortcodes(stringToParse string, p *Page, t tpl.Template) (string, map[string]string) {

	content, shortcodes, err := extractShortcodes(stringToParse, p, t)
	renderedShortcodes := make(map[string]string)

	if err != nil {
		//  try to render what we have whilst logging the error
		jww.ERROR.Println(err.Error())
	}

	for key, sc := range shortcodes {
		if sc.err != nil {
			// need to have something to replace with
			renderedShortcodes[key] = ""
		} else {
			renderedShortcodes[key] = renderShortcode(sc, p, t)
		}
	}

	return content, renderedShortcodes

}

// pageTokens state:
// - before: positioned just before the shortcode start
// - after: shortcode(s) consumed (plural when they are nested)
func extractShortcode(pt *pageTokens, p *Page, t tpl.Template) (shortcode, error) {
	sc := shortcode{}
	var isInner = false

	var currItem item
	var cnt = 0

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
				nested, err := extractShortcode(pt, p, t)
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
			if !isInner {
				next := pt.peek()
				if next.typ == tError {
					// return that error, more specific
					continue
				}
				return sc, fmt.Errorf("Shortcode '%s' has no .Inner, yet a closing tag was provided", next.val)
			}
			pt.consume(2)
			return sc, nil
		case tText:
			sc.inner = append(sc.inner, currItem.val)
		case tScName:
			sc.name = currItem.val
			tmpl := GetTemplate(sc.name, t)

			if tmpl == nil {
				return sc, fmt.Errorf("Unable to locate template for shortcode '%s' in page %s", sc.name, p.BaseFileName())
			}
			isInner = isInnerShortcode(tmpl)

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
					params := sc.params.(map[string]string)
					params[currItem.val] = pt.next().val
				}
			} else {
				// positional params
				if sc.params == nil {
					var params []string
					params = append(params, currItem.val)
					sc.params = params
				} else {
					params := sc.params.([]string)
					params = append(params, currItem.val)
					sc.params = params
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

func extractShortcodes(stringToParse string, p *Page, t tpl.Template) (string, map[string]shortcode, error) {

	shortCodes := make(map[string]shortcode)

	startIdx := strings.Index(stringToParse, "{{")

	// short cut for docs with no shortcodes
	if startIdx < 0 {
		return stringToParse, shortCodes, nil
	}

	// the parser takes a string;
	// since this is an internal API, it could make sense to use the mutable []byte all the way, but
	// it seems that the time isn't really spent in the byte copy operations, and the impl. gets a lot cleaner
	pt := &pageTokens{lexer: newShortcodeLexer("parse-page", stringToParse, pos(startIdx))}

	id := 1 // incremented id, will be appended onto temp. shortcode placeholders
	var result bytes.Buffer

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currItem item
	var currShortcode shortcode
	var err error

Loop:
	for {
		currItem = pt.next()

		switch currItem.typ {
		case tText:
			result.WriteString(currItem.val)
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup:
			// let extractShortcode handle left delim (will do so recursively)
			pt.backup()
			if currShortcode, err = extractShortcode(pt, p, t); err != nil {
				return result.String(), shortCodes, err
			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			placeHolder := createShortcodePlaceholder(id)
			result.WriteString(placeHolder)
			shortCodes[placeHolder] = currShortcode
			id++
		case tEOF:
			break Loop
		case tError:
			err := fmt.Errorf("%s:%d: %s",
				p.BaseFileName(), (p.lineNumRawContentStart() + pt.lexer.lineNum() - 1), currItem)
			currShortcode.err = err
			return result.String(), shortCodes, err
		}
	}

	return result.String(), shortCodes, nil

}

// Replace prefixed shortcode tokens (HUGOSHORTCODE-1, HUGOSHORTCODE-2) with the real content.
// This assumes that all tokens exist in the input string and that they are in order.
// numReplacements = -1 will do len(replacements), and it will always start from the beginning (1)
// wrapped = true means that the token has been wrapped in {@{@/@}@}
func replaceShortcodeTokens(source []byte, prefix string, numReplacements int, wrapped bool, replacements map[string]string) ([]byte, error) {

	if numReplacements < 0 {
		numReplacements = len(replacements)
	}

	if numReplacements == 0 {
		return source, nil
	}

	newLen := len(source)

	for i := 1; i <= numReplacements; i++ {
		key := prefix + "-" + strconv.Itoa(i)

		if wrapped {
			key = "{@{@" + key + "@}@}"
		}
		val := []byte(replacements[key])

		newLen += (len(val) - len(key))
	}

	buff := make([]byte, newLen)

	width := 0
	start := 0

	for i := 0; i < numReplacements; i++ {
		tokenNum := i + 1
		oldVal := prefix + "-" + strconv.Itoa(tokenNum)
		if wrapped {
			oldVal = "{@{@" + oldVal + "@}@}"
		}
		newVal := []byte(replacements[oldVal])
		j := start

		k := bytes.Index(source[start:], []byte(oldVal))

		if k < 0 {
			// this should never happen, but let the caller decide to panic or not
			return nil, fmt.Errorf("illegal state in content; shortcode token #%d is missing or out of order (%q)", tokenNum, source)
		}
		j += k

		width += copy(buff[width:], source[start:j])
		width += copy(buff[width:], newVal)
		start = j + len(oldVal)
	}
	width += copy(buff[width:], source[start:])
	return buff[0:width], nil
}

func GetTemplate(name string, t tpl.Template) *template.Template {
	if x := t.Lookup("shortcodes/" + name + ".html"); x != nil {
		return x
	}
	if x := t.Lookup("theme/shortcodes/" + name + ".html"); x != nil {
		return x
	}
	return t.Lookup("_internal/shortcodes/" + name + ".html")
}

func ShortcodeRender(tmpl *template.Template, data *ShortcodeWithPage) string {
	buffer := new(bytes.Buffer)
	err := tmpl.Execute(buffer, data)
	if err != nil {
		jww.ERROR.Println("error processing shortcode", tmpl.Name(), "\n ERR:", err)
		jww.WARN.Println(data)
	}
	return buffer.String()
}
