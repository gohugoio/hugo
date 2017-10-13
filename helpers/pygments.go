// Copyright 2016 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

const pygmentsBin = "pygmentize"

// hasPygments checks to see if Pygments is installed and available
// on the system.
func hasPygments() bool {
	if _, err := exec.LookPath(pygmentsBin); err != nil {
		return false
	}
	return true
}

type highlighters struct {
	cs          *ContentSpec
	ignoreCache bool
	cacheDir    string
}

func newHiglighters(cs *ContentSpec) highlighters {
	return highlighters{cs: cs, ignoreCache: cs.cfg.GetBool("ignoreCache"), cacheDir: cs.cfg.GetString("cacheDir")}
}

func (h highlighters) chromaHighlight(code, lang, optsStr string) (string, error) {
	opts, err := h.cs.parsePygmentsOpts(optsStr)
	if err != nil {
		jww.ERROR.Print(err.Error())
		return code, err
	}

	style, found := opts["style"]
	if !found || style == "" {
		style = "friendly"
	}

	f, err := h.cs.chromaFormatterFromOptions(opts)
	if err != nil {
		jww.ERROR.Print(err.Error())
		return code, err
	}

	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	err = chromaHighlight(b, code, lang, style, f)
	if err != nil {
		jww.ERROR.Print(err.Error())
		return code, err
	}

	return h.injectCodeTag(`<div class="highlight">`+b.String()+"</div>", lang), nil
}

func (h highlighters) pygmentsHighlight(code, lang, optsStr string) (string, error) {
	options, err := h.cs.createPygmentsOptionsString(optsStr)

	if err != nil {
		jww.ERROR.Print(err.Error())
		return code, nil
	}

	// Try to read from cache first
	hash := sha1.New()
	io.WriteString(hash, code)
	io.WriteString(hash, lang)
	io.WriteString(hash, options)

	fs := hugofs.Os

	var cachefile string

	if !h.ignoreCache && h.cacheDir != "" {
		cachefile = filepath.Join(h.cacheDir, fmt.Sprintf("pygments-%x", hash.Sum(nil)))

		exists, err := Exists(cachefile, fs)
		if err != nil {
			jww.ERROR.Print(err.Error())
			return code, nil
		}
		if exists {
			f, err := fs.Open(cachefile)
			if err != nil {
				jww.ERROR.Print(err.Error())
				return code, nil
			}

			s, err := ioutil.ReadAll(f)
			if err != nil {
				jww.ERROR.Print(err.Error())
				return code, nil
			}

			return string(s), nil
		}
	}

	// No cache file, render and cache it
	var out bytes.Buffer
	var stderr bytes.Buffer

	var langOpt string
	if lang == "" {
		langOpt = "-g" // Try guessing the language
	} else {
		langOpt = "-l" + lang
	}

	cmd := exec.Command(pygmentsBin, langOpt, "-fhtml", "-O", options)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		jww.ERROR.Print(stderr.String())
		return code, err
	}

	str := string(normalizeExternalHelperLineFeeds([]byte(out.String())))

	str = h.injectCodeTag(str, lang)

	if !h.ignoreCache && cachefile != "" {
		// Write cache file
		if err := WriteToDisk(cachefile, strings.NewReader(str), fs); err != nil {
			jww.ERROR.Print(stderr.String())
		}
	}

	return str, nil
}

var preRe = regexp.MustCompile(`(?s)(.*?<pre.*?>)(.*?)(</pre>)`)

func (h highlighters) injectCodeTag(code, lang string) string {
	if lang == "" {
		return code
	}
	codeTag := fmt.Sprintf(`<code class="language-%s" data-lang="%s">`, lang, lang)
	return preRe.ReplaceAllString(code, fmt.Sprintf("$1%s$2</code>$3", codeTag))
}

func chromaHighlight(w io.Writer, source, lexer, style string, f chroma.Formatter) error {
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	if f == nil {
		f = formatters.Fallback
	}

	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return f.Format(w, s, it)
}

var pygmentsKeywords = make(map[string]bool)

func init() {
	pygmentsKeywords["encoding"] = true
	pygmentsKeywords["outencoding"] = true
	pygmentsKeywords["nowrap"] = true
	pygmentsKeywords["full"] = true
	pygmentsKeywords["title"] = true
	pygmentsKeywords["style"] = true
	pygmentsKeywords["noclasses"] = true
	pygmentsKeywords["classprefix"] = true
	pygmentsKeywords["cssclass"] = true
	pygmentsKeywords["cssstyles"] = true
	pygmentsKeywords["prestyles"] = true
	pygmentsKeywords["linenos"] = true
	pygmentsKeywords["hl_lines"] = true
	pygmentsKeywords["linenostart"] = true
	pygmentsKeywords["linenostep"] = true
	pygmentsKeywords["linenospecial"] = true
	pygmentsKeywords["nobackground"] = true
	pygmentsKeywords["lineseparator"] = true
	pygmentsKeywords["lineanchors"] = true
	pygmentsKeywords["linespans"] = true
	pygmentsKeywords["anchorlinenos"] = true
	pygmentsKeywords["startinline"] = true
}

func parseOptions(defaults map[string]string, in string) (map[string]string, error) {
	in = strings.Trim(in, " ")
	opts := make(map[string]string)

	if defaults != nil {
		for k, v := range defaults {
			opts[k] = v
		}
	}

	if in == "" {
		return opts, nil
	}

	for _, v := range strings.Split(in, ",") {
		keyVal := strings.Split(v, "=")
		key := strings.ToLower(strings.Trim(keyVal[0], " "))
		if len(keyVal) != 2 || !pygmentsKeywords[key] {
			return opts, fmt.Errorf("invalid Pygments option: %s", key)
		}
		opts[key] = keyVal[1]
	}

	return opts, nil
}

func createOptionsString(options map[string]string) string {
	var keys []string
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var optionsStr string
	for i, k := range keys {
		optionsStr += fmt.Sprintf("%s=%s", k, options[k])
		if i < len(options)-1 {
			optionsStr += ","
		}
	}

	return optionsStr
}

func parseDefaultPygmentsOpts(cfg config.Provider) (map[string]string, error) {
	options, err := parseOptions(nil, cfg.GetString("pygmentsOptions"))
	if err != nil {
		return nil, err
	}

	if cfg.IsSet("pygmentsStyle") {
		options["style"] = cfg.GetString("pygmentsStyle")
	}

	if cfg.IsSet("pygmentsUseClasses") {
		if cfg.GetBool("pygmentsUseClasses") {
			options["noclasses"] = "false"
		} else {
			options["noclasses"] = "true"
		}

	}

	if _, ok := options["encoding"]; !ok {
		options["encoding"] = "utf8"
	}

	return options, nil
}

func (cs *ContentSpec) chromaFormatterFromOptions(pygmentsOpts map[string]string) (chroma.Formatter, error) {
	var options = []html.Option{html.TabWidth(4)}

	if pygmentsOpts["noclasses"] == "false" {
		options = append(options, html.WithClasses())
	}

	lineNumbers := pygmentsOpts["linenos"]
	if lineNumbers != "" {
		options = append(options, html.WithLineNumbers())
		if lineNumbers != "inline" {
			options = append(options, html.LineNumbersInTable())
		}
	}

	startLineStr := pygmentsOpts["linenostart"]
	var startLine = 1
	if startLineStr != "" {

		line, err := strconv.Atoi(strings.TrimSpace(startLineStr))
		if err == nil {
			startLine = line
			options = append(options, html.BaseLineNumber(startLine))
		}
	}

	hlLines := pygmentsOpts["hl_lines"]

	if hlLines != "" {
		ranges, err := hlLinesToRanges(startLine, hlLines)

		if err == nil {
			options = append(options, html.HighlightLines(ranges))
		}
	}

	return html.New(options...), nil
}

func (cs *ContentSpec) parsePygmentsOpts(in string) (map[string]string, error) {
	opts, err := parseOptions(cs.defatultPygmentsOpts, in)
	if err != nil {
		return nil, err
	}
	return opts, nil

}

func (cs *ContentSpec) createPygmentsOptionsString(in string) (string, error) {
	opts, err := cs.parsePygmentsOpts(in)
	if err != nil {
		return "", err
	}
	return createOptionsString(opts), nil
}

// startLine compansates for https://github.com/alecthomas/chroma/issues/30
func hlLinesToRanges(startLine int, s string) ([][2]int, error) {
	var ranges [][2]int
	s = strings.TrimSpace(s)

	if s == "" {
		return ranges, nil
	}

	// Variants:
	// 1 2 3 4
	// 1-2 3-4
	// 1-2 3
	// 1 3-4
	// 1    3-4
	fields := strings.Split(s, " ")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		numbers := strings.Split(field, "-")
		var r [2]int
		first, err := strconv.Atoi(numbers[0])
		if err != nil {
			return ranges, err
		}
		first = first + startLine - 1
		r[0] = first
		if len(numbers) > 1 {
			second, err := strconv.Atoi(numbers[1])
			if err != nil {
				return ranges, err
			}
			second = second + startLine - 1
			r[1] = second
		} else {
			r[1] = first
		}

		ranges = append(ranges, r)
	}
	return ranges, nil

}
