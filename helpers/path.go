// Copyright © 2014 Steve Francia <spf@spf13.com>.
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
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

var sanitizeRegexp = regexp.MustCompile("[^a-zA-Z0-9./_-]")

// Take a string with any characters and replace it so the string could be used in a path.
// E.g. Social Media -> social-media
func MakePath(s string) string {
	return UnicodeSanitize(strings.ToLower(strings.Replace(strings.TrimSpace(s), " ", "-", -1)))
}

func Sanitize(s string) string {
	return sanitizeRegexp.ReplaceAllString(s, "")
}

func UnicodeSanitize(s string) string {
	source := []rune(s)
	target := make([]rune, 0, len(source))

	for _, r := range source {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '/' || r == '_' || r == '-' {
			target = append(target, r)
		}
	}

	return string(target)
}

func ReplaceExtension(path string, newExt string) string {
	f, _ := FileAndExt(path)
	return f + "." + newExt
}

// Check if Exists && is Directory
func DirExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileAndExt(in string) (name string, ext string) {
	ext = path.Ext(in)
	base := path.Base(in)

	if strings.Contains(base, ".") {
		name = base[:strings.LastIndex(base, ".")]
	} else {
		name = in
	}

	return
}

func PathPrep(ugly bool, in string) string {
	if ugly {
		return Uglify(in)
	} else {
		return PrettifyPath(in)
	}
}

// /section/name.html -> /section/name/index.html
// /section/name/  -> /section/name/index.html
// /section/name/index.html -> /section/name/index.html
func PrettifyPath(in string) string {
	if path.Ext(in) == "" {
		// /section/name/  -> /section/name/index.html
		if len(in) < 2 {
			return "/"
		}
		return path.Join(path.Clean(in), "index.html")
	} else {
		name, ext := FileAndExt(in)
		if name == "index" {
			// /section/name/index.html -> /section/name/index.html
			return path.Clean(in)
		} else {
			// /section/name.html -> /section/name/index.html
			return path.Join(path.Dir(in), name, "index"+ext)
		}
	}
}
