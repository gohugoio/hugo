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

package helpers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
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
func Highlight(code, lang, optsStr string) string {

	if !HasPygments() {
		jww.WARN.Println("Highlighting requires Pygments to be installed and in the path")
		return code
	}

	options, err := parsePygmentsOpts(optsStr)

	if err != nil {
		jww.ERROR.Print(err.Error())
		return code
	}

	// Try to read from cache first
	hash := sha1.New()
	io.WriteString(hash, code)
	io.WriteString(hash, lang)
	io.WriteString(hash, options)

	fs := hugofs.OsFs

	cacheDir := viper.GetString("CacheDir")
	var cachefile string

	if cacheDir != "" {
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

	cmd := exec.Command(pygmentsBin, "-l"+lang, "-fhtml", "-O", options)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		jww.ERROR.Print(stderr.String())
		return code
	}

	if cachefile != "" {
		// Write cache file
		if err := WriteToDisk(cachefile, bytes.NewReader(out.Bytes()), fs); err != nil {
			jww.ERROR.Print(stderr.String())
		}
	}

	return out.String()
}

var pygmentsKeywords = make(map[string]bool)

func init() {
	pygmentsKeywords["style"] = true
	pygmentsKeywords["encoding"] = true
	pygmentsKeywords["noclasses"] = true
	pygmentsKeywords["hl_lines"] = true
	pygmentsKeywords["linenos"] = true
	pygmentsKeywords["classprefix"] = true
}

func parsePygmentsOpts(in string) (string, error) {

	in = strings.Trim(in, " ")

	style := viper.GetString("PygmentsStyle")

	noclasses := "true"
	if viper.GetBool("PygmentsUseClasses") {
		noclasses = "false"
	}

	if len(in) == 0 {
		return fmt.Sprintf("style=%s,noclasses=%s,encoding=utf8", style, noclasses), nil
	}

	options := make(map[string]string)

	o := strings.Split(in, ",")
	for _, v := range o {
		keyVal := strings.Split(v, "=")
		key := strings.ToLower(strings.Trim(keyVal[0], " "))
		if len(keyVal) != 2 || !pygmentsKeywords[key] {
			return "", fmt.Errorf("invalid Pygments option: %s", key)
		}
		options[key] = keyVal[1]
	}

	if _, ok := options["style"]; !ok {
		options["style"] = style
	}

	if _, ok := options["noclasses"]; !ok {
		options["noclasses"] = noclasses
	}

	if _, ok := options["encoding"]; !ok {
		options["encoding"] = "utf8"
	}

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
	return optionsStr, nil
}
