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
	"sort"
	"strings"

	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
)

const pygmentsBin = "pygmentize"

// HasPygments checks to see if Pygments is installed and available
// on the system.
func HasPygments() bool {
	if _, err := exec.LookPath(pygmentsBin); err != nil {
		return false
	}
	return true
}

// Highlight takes some code and returns highlighted code.
func Highlight(cfg config.Provider, code, lang, optsStr string) string {
	if !HasPygments() {
		jww.WARN.Println("Highlighting requires Pygments to be installed and in the path")
		return code
	}

	options, err := parsePygmentsOpts(cfg, optsStr)

	if err != nil {
		jww.ERROR.Print(err.Error())
		return code
	}

	// Try to read from cache first
	hash := sha1.New()
	io.WriteString(hash, code)
	io.WriteString(hash, lang)
	io.WriteString(hash, options)

	fs := hugofs.Os

	ignoreCache := cfg.GetBool("ignoreCache")
	cacheDir := cfg.GetString("cacheDir")
	var cachefile string

	if !ignoreCache && cacheDir != "" {
		cachefile = filepath.Join(cacheDir, fmt.Sprintf("pygments-%x", hash.Sum(nil)))

		exists, err := Exists(cachefile, fs)
		if err != nil {
			jww.ERROR.Print(err.Error())
			return code
		}
		if exists {
			f, err := fs.Open(cachefile)
			if err != nil {
				jww.ERROR.Print(err.Error())
				return code
			}

			s, err := ioutil.ReadAll(f)
			if err != nil {
				jww.ERROR.Print(err.Error())
				return code
			}

			return string(s)
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
		return code
	}

	str := string(normalizeExternalHelperLineFeeds([]byte(out.String())))

	// inject code tag into Pygments output
	if lang != "" && strings.Contains(str, "<pre>") {
		codeTag := fmt.Sprintf(`<pre><code class="language-%s" data-lang="%s">`, lang, lang)
		str = strings.Replace(str, "<pre>", codeTag, 1)
		str = strings.Replace(str, "</pre>", "</code></pre>", 1)
	}

	if !ignoreCache && cachefile != "" {
		// Write cache file
		if err := WriteToDisk(cachefile, strings.NewReader(str), fs); err != nil {
			jww.ERROR.Print(stderr.String())
		}
	}

	return str
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

func parseOptions(options map[string]string, in string) error {
	in = strings.Trim(in, " ")

	if in == "" {
		return nil
	}

	for _, v := range strings.Split(in, ",") {
		keyVal := strings.Split(v, "=")
		key := strings.ToLower(strings.Trim(keyVal[0], " "))
		if len(keyVal) != 2 || !pygmentsKeywords[key] {
			return fmt.Errorf("invalid Pygments option: %s", key)
		}
		options[key] = keyVal[1]
	}

	return nil
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
	options := make(map[string]string)
	err := parseOptions(options, cfg.GetString("pygmentsOptions"))
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

func parsePygmentsOpts(cfg config.Provider, in string) (string, error) {
	options, err := parseDefaultPygmentsOpts(cfg)
	if err != nil {
		return "", err
	}

	err = parseOptions(options, in)
	if err != nil {
		return "", err
	}

	return createOptionsString(options), nil
}
