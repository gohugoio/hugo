// Copyright © 2013 Steve Francia <spf@spf13.com>.
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
	"strings"
	"unicode"
)

var _ = fmt.Println

type ShortcodeFunc func([]string) string

type Shortcode struct {
	Name string
	Func ShortcodeFunc
}

type ShortcodeWithPage struct {
	Params interface{}
	Page   *Page
}

type Shortcodes map[string]ShortcodeFunc

func ShortcodesHandle(stringToParse string, p *Page, t Template) string {
	posStart := strings.Index(stringToParse, "{{%")
	if posStart > 0 {
		posEnd := strings.Index(stringToParse[posStart:], "%}}") + posStart
		if posEnd > posStart {
			name, par := SplitParams(stringToParse[posStart+3 : posEnd])
			params := Tokenize(par)
			var data = &ShortcodeWithPage{Params: params, Page: p}
			newString := stringToParse[:posStart] + ShortcodeRender(name, data, t) + ShortcodesHandle(stringToParse[posEnd+3:], p, t)
			return newString
		}
	}
	return stringToParse
}

func StripShortcodes(stringToParse string) string {
	posStart := strings.Index(stringToParse, "{{%")
	if posStart > 0 {
		posEnd := strings.Index(stringToParse[posStart:], "%}}") + posStart
		if posEnd > posStart {
			newString := stringToParse[:posStart] + StripShortcodes(stringToParse[posEnd+3:])
			return newString
		}
	}
	return stringToParse
}

func Tokenize(in string) interface{} {
	first := strings.Fields(in)
	var final = make([]string, 0)
	var keys = make([]string, 0)
	inQuote := false
	start := 0

	for i, v := range first {
		index := strings.Index(v, "=")

		if !inQuote {
			if index > 1 {
				keys = append(keys, v[:index])
				v = v[index+1:]
			}
		}

		if !strings.HasPrefix(v, "&ldquo;") && !inQuote {
			final = append(final, v)
		} else if inQuote && strings.HasSuffix(v, "&rdquo;") && !strings.HasSuffix(v, "\\\"") {
			first[i] = v[:len(v)-7]
			final = append(final, strings.Join(first[start:i+1], " "))
			inQuote = false
		} else if strings.HasPrefix(v, "&ldquo;") && !inQuote {
			if strings.HasSuffix(v, "&rdquo;") {
				final = append(final, v[7:len(v)-7])
			} else {
				start = i
				first[i] = v[7:]
				inQuote = true
			}
		}

		// No closing "... just make remainder the final token
		if inQuote && i == len(first) {
			final = append(final, first[start:len(first)]...)
		}
	}

	if len(keys) > 0 {
		var m = make(map[string]string)
		for i, k := range keys {
			m[k] = final[i]
		}

		return m
	}

	return final
}

func SplitParams(in string) (name string, par2 string) {
	i := strings.IndexFunc(strings.TrimSpace(in), unicode.IsSpace)
	if i < 1 {
		return strings.TrimSpace(in), ""
	}

	return strings.TrimSpace(in[:i+1]), strings.TrimSpace(in[i+1:])
}

func ShortcodeRender(name string, data *ShortcodeWithPage, t Template) string {
	buffer := new(bytes.Buffer)
	t.ExecuteTemplate(buffer, "shortcodes/"+name+".html", data)
	return buffer.String()
}
