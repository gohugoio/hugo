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
    "github.com/spf13/hugo/template/bundle"
    "html/template"
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
    Inner  template.HTML
    Page   *Page
}

type Shortcodes map[string]ShortcodeFunc

func ShortcodesHandle(stringToParse string, p *Page, t bundle.Template) string {
    leadStart := strings.Index(stringToParse, `{{%`)
    if leadStart >= 0 {
        leadEnd := strings.Index(stringToParse[leadStart:], `%}}`) + leadStart
        if leadEnd > leadStart {
            name, par := SplitParams(stringToParse[leadStart+3 : leadEnd])
            tmpl := GetTemplate(name, t)
            if tmpl == nil {
                return stringToParse
            }
            params := Tokenize(par)
            // Always look for closing tag.
            endStart, endEnd := FindEnd(stringToParse[leadEnd:], name)
            var data = &ShortcodeWithPage{Params: params, Page: p}
            if endStart > 0 {
                s := stringToParse[leadEnd+3 : leadEnd+endStart]
                data.Inner = template.HTML(CleanP(ShortcodesHandle(s, p, t)))
                remainder := CleanP(stringToParse[leadEnd+endEnd:])

                return CleanP(stringToParse[:leadStart]) +
                    ShortcodeRender(tmpl, data) +
                    CleanP(ShortcodesHandle(remainder, p, t))
            }
            return CleanP(stringToParse[:leadStart]) +
                ShortcodeRender(tmpl, data) +
                CleanP(ShortcodesHandle(stringToParse[leadEnd+3:], p,
                    t))
        }
    }
    return stringToParse
}

// Clean up odd behavior when closing tag is on first line
// or opening tag is on the last line due to extra line in markdown file
func CleanP(str string) string {
    if strings.HasSuffix(strings.TrimSpace(str), "<p>") {
        idx := strings.LastIndex(str, "<p>")
        str = str[:idx]
    }

    if strings.HasPrefix(strings.TrimSpace(str), "</p>") {
        str = str[strings.Index(str, "</p>")+5:]
    }

    return str
}

func FindEnd(str string, name string) (int, int) {
    var endPos int
    var startPos int
    var try []string

    try = append(try, "{{% /"+name+" %}}")
    try = append(try, "{{% /"+name+"%}}")
    try = append(try, "{{%/"+name+"%}}")
    try = append(try, "{{%/"+name+" %}}")

    lowest := len(str)
    for _, x := range try {
        start := strings.Index(str, x)
        if start < lowest && start > 0 {
            startPos = start
            endPos = startPos + len(x)
        }
    }

    return startPos, endPos
}

func GetTemplate(name string, t bundle.Template) *template.Template {
    if x := t.Lookup("shortcodes/" + name + ".html"); x != nil {
        return x
    }
    return t.Lookup("_internal/shortcodes/" + name + ".html")
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

    // if don't need to parse, don't parse.
    if strings.Index(in, " ") < 0 && strings.Index(in, "=") < 1 {
        return append(final, in)
    }

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

        // Adjusted to handle htmlencoded and non htmlencoded input
        if !strings.HasPrefix(v, "&ldquo;") && !strings.HasPrefix(v, "\"") && !inQuote {
            final = append(final, v)
        } else if inQuote && (strings.HasSuffix(v, "&rdquo;") ||
            strings.HasSuffix(v, "\"")) && !strings.HasSuffix(v, "\\\"") {
            if strings.HasSuffix(v, "\"") {
                first[i] = v[:len(v)-1]
            } else {
                first[i] = v[:len(v)-7]
            }
            final = append(final, strings.Join(first[start:i+1], " "))
            inQuote = false
        } else if (strings.HasPrefix(v, "&ldquo;") ||
            strings.HasPrefix(v, "\"")) && !inQuote {
            if strings.HasSuffix(v, "&rdquo;") || strings.HasSuffix(v,
                "\"") {
                if strings.HasSuffix(v, "\"") {
                    if len(v) > 1 {
                        final = append(final, v[1:len(v)-1])
                    } else {
                        final = append(final, "")
                    }
                } else {
                    final = append(final, v[7:len(v)-7])
                }
            } else {
                start = i
                if strings.HasPrefix(v, "\"") {
                    first[i] = v[1:]
                } else {
                    first[i] = v[7:]
                }
                inQuote = true
            }
        }

        // No closing "... just make remainder the final token
        if inQuote && i == len(first) {
            final = append(final, first[start:]...)
        }
    }

    if len(keys) > 0 && (len(keys) != len(final)) {
        panic("keys and final different lengths")
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

func ShortcodeRender(tmpl *template.Template, data *ShortcodeWithPage) string {
    buffer := new(bytes.Buffer)
    err := tmpl.Execute(buffer, data)
    if err != nil {
        fmt.Println("error processing shortcode", tmpl.Name(), "\n ERR:", err)
    }
    return buffer.String()
}
