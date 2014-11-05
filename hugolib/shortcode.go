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
	"github.com/spf13/hugo/helpers"
	jww "github.com/spf13/jwalterweatherman"
	"html/template"
	"reflect"
	"strconv"
	"strings"
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

type shortCode struct {
	name     string
	inner    string
	params   interface{} // map or array
	doMarkup bool
}

func (sc shortCode) String() string {
	// for testing (mostly), so any change here will break tests!
	return fmt.Sprintf("%s(%q, %t){%s}", sc.name, sc.params, sc.doMarkup, sc.inner)
}

// all in  one go: extract, render and replace
// only used for testing
func ShortcodesHandle(stringToParse string, page *Page, t Template) string {

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

func extractAndRenderShortcodes(stringToParse string, p *Page, t Template) (string, map[string]string) {

	content, shortcodes, err := extractShortcodes(stringToParse, p)
	renderedShortcodes := make(map[string]string)

	if err != nil {
		//  try to render what we have whilst logging the error
		jww.ERROR.Println(err.Error())
	}

	for key, sc := range shortcodes {
		var data = &ShortcodeWithPage{Params: sc.params, Page: p}
		if sc.inner != "" {
			if sc.doMarkup {
				data.Inner = template.HTML(helpers.RenderBytes([]byte(sc.inner), p.guessMarkupType(), p.UniqueId()))
			} else {
				data.Inner = template.HTML(sc.inner)
			}

		}

		tmpl := GetTemplate(sc.name, t)

		if tmpl == nil {
			jww.ERROR.Printf("Unable to locate template for shortcode '%s' in page %s", sc.name, p.BaseFileName())
			continue
		}
		renderedShortcode := ShortcodeRender(tmpl, data)
		renderedShortcodes[key] = renderedShortcode
	}

	return content, renderedShortcodes

}

func extractShortcodes(stringToParse string, p *Page) (string, map[string]shortCode, error) {

	shortCodes := make(map[string]shortCode)

	startIdx := strings.Index(stringToParse, "{{")

	// short cut for docs with no shortcodes
	if startIdx < 0 {
		return stringToParse, shortCodes, nil
	}

	// the parser takes a string;
	// since this is an internal API, it could make sense to use the mutable []byte all the way, but
	// it seems that the time isn't really spent in the byte copy operations, and the impl. gets a lot cleaner
	t := pageTokens{lexer: newShortcodeLexer("parse-page", stringToParse, pos(startIdx))}

	id := 1 // incremented id, will be appended onto temp. shortcode placeholders
	var result bytes.Buffer

	// the parser is guaranteed to return items in proper order or fail, so …
	// … it's safe to keep some "global" state
	var currItem item
	var currShortcode shortCode

Loop:
	for {
		currItem = t.next()

		switch currItem.typ {
		case tText:
			result.WriteString(currItem.val)
		case tLeftDelimScWithMarkup, tLeftDelimScNoMarkup:
			currShortcode = shortCode{}
			currShortcode.doMarkup = currItem.typ == tLeftDelimScWithMarkup
		case tRightDelimScWithMarkup, tRightDelimScNoMarkup:
			// need 3-token look-ahead here in the worst case looking for
			// some shortcode inner content
			if t.peek().typ == tText {
				textToken := t.next()
				if t.isLeftShortcodeDelim(t.peek()) {
					delim := t.next()
					if t.peek().typ == tScClose {
						currShortcode.inner = textToken.val
						// consume the shortcode close
						t.consume(3)
					} else {
						// shortcode is open
						t.backup3(textToken, delim)
					}
				} else {
					// let the text item be handled elsewhere
					t.backup2(textToken)
				}
			} else if t.isLeftShortcodeDelim(t.peek()) {
				// check for empty inner content
				delim := t.next()
				if t.peek().typ == tScClose {
					t.consume(3)
				} else {
					// new shortcode
					t.backup2(delim)
				}

			}

			if currShortcode.params == nil {
				currShortcode.params = make([]string, 0)
			}

			// wrap it in a block level element to let it be left alone by the markup engine
			placeHolder := fmt.Sprintf("<div>%s-%d</div>", shortcodePlaceholderPrefix, id)
			result.WriteString(placeHolder)
			shortCodes[placeHolder] = currShortcode
			id++
		case tScName:
			currShortcode.name = currItem.val
		case tScParam:
			if !t.isValueNext() {
				continue
			} else if t.peek().typ == tScParamVal {
				// named params
				if currShortcode.params == nil {
					params := make(map[string]string)
					params[currItem.val] = t.next().val
					currShortcode.params = params
				} else {
					params := currShortcode.params.(map[string]string)
					params[currItem.val] = t.next().val
				}
			} else {
				// positional params
				if currShortcode.params == nil {
					var params []string
					params = append(params, currItem.val)
					currShortcode.params = params
				} else {
					params := currShortcode.params.([]string)
					params = append(params, currItem.val)
					currShortcode.params = params
				}
			}

		case tEOF:
			break Loop
		case tError:
			return result.String(), shortCodes, fmt.Errorf("%s:%d: %s",
				p.BaseFileName(), (p.lineNumRawContentStart() + t.lexer.lineNum() - 1), currItem)
		}
	}

	return result.String(), shortCodes, nil

}

// Replace prefixed shortcode tokens (HUGOSHORTCODE-1, HUGOSHORTCODE-2) with the real content.
// This assumes that all tokens exist in the input string and that they are in order.
// numReplacements = -1 will do len(replacements), and it will always start from the beginning (1)
// wrappendInDiv = true means that the token is wrapped in a <div></div>
func replaceShortcodeTokens(source []byte, prefix string, numReplacements int, wrappedInDiv bool, replacements map[string]string) ([]byte, error) {

	if numReplacements < 0 {
		numReplacements = len(replacements)
	}

	if numReplacements == 0 {
		return source, nil
	}

	newLen := len(source)

	for i := 1; i <= numReplacements; i++ {
		key := prefix + "-" + strconv.Itoa(i)

		if wrappedInDiv {
			key = "<div>" + key + "</div>"
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
		if wrappedInDiv {
			oldVal = "<div>" + oldVal + "</div>"
		}
		newVal := []byte(replacements[oldVal])
		j := start

		k := bytes.Index(source[start:], []byte(oldVal))
		if k < 0 {
			// this should never happen, but let the caller decide to panic or not
			return nil, fmt.Errorf("illegal state in content; shortcode token #%d is missing or out of order", tokenNum)
		}
		j += k

		width += copy(buff[width:], source[start:j])
		width += copy(buff[width:], newVal)
		start = j + len(oldVal)
	}
	width += copy(buff[width:], source[start:])
	return buff[0:width], nil
}

func GetTemplate(name string, t Template) *template.Template {
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
