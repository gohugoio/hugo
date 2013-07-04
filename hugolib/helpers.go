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
	"github.com/kr/pretty"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var sanitizeRegexp = regexp.MustCompile("[^a-zA-Z0-9/_-]")

// TODO: Make these wrappers private
// Wrapper around Fprintf taking verbose flag in account.
func Printvf(format string, a ...interface{}) {
	//if *verbose {
	fmt.Fprintf(os.Stderr, format, a...)
	//}
}

func Printer(x interface{}) {
	fmt.Printf("%#v", pretty.Formatter(x))
	fmt.Println("")
}

// Wrapper around Fprintln taking verbose flag in account.
func Printvln(a ...interface{}) {
	//if *verbose {
	fmt.Fprintln(os.Stderr, a...)
	//}
}

func FatalErr(str string) {
	fmt.Println(str)
	os.Exit(1)
}

func PrintErr(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}

func Error(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}

func interfaceToStringToDate(i interface{}) time.Time {
	s := interfaceToString(i)
	d, e := time.Parse("02 Jan 06 15:04 MST", s)

	if e != nil {
		d, e = time.Parse("2006-01-02", s)
	}

	if e != nil {
		d, e = time.Parse("02 Jan 06", s)
	}

	return d

}

func interfaceToBool(i interface{}) bool {
	switch b := i.(type) {
	case bool:
		return b
	default:
		Error("Only Boolean values are supported for this JSON key")
	}

	return false

}

func interfaceArrayToStringArray(i interface{}) []string {
	var a []string

	switch vv := i.(type) {
	case []interface{}:
		for _, u := range vv {
			a = append(a, interfaceToString(u))
		}
	}

	return a
}

func interfaceToString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	default:
		Error("Only Strings are supported for this JSON key")
	}

	return ""
}

// Check if Exists && is Directory
func dirExists(path string) (bool, error) {
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
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func mkdirIf(path string) {
	err := os.Mkdir(path, 0777)
	if err != nil && os.IsNotExist(err) {
		fmt.Println(err)
	}
}

func Urlize(url string) string {
	return Sanitize(strings.ToLower(strings.Replace(strings.TrimSpace(url), " ", "-", -1)))
}

func Gt(a interface{}, b interface{}) bool {
	var left, right int64
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, 64)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, 64)
	}

	return left > right
}

func IsSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}

func ReturnWhenSet(a interface{}, index int) interface{} {
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Slice:
		if av.Len() > index {

			avv := av.Index(index)
			switch avv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return avv.Int()
			case reflect.String:
				return avv.String()
			}
		}
	}

	return ""
}

func Sanitize(s string) string {
	return sanitizeRegexp.ReplaceAllString(s, "")
}

func fileExt(path string) (file, ext string) {
	if strings.Contains(path, ".") {
		i := len(path) - 1
		for path[i] != '.' {
			i--
		}
		return path[:i], path[i+1:]
	}
	return path, ""
}

func replaceExtension(path string, newExt string) string {
	f, _ := fileExt(path)
	return f + "." + newExt
}

func TotalWords(s string) int {
	return len(strings.Fields(s))
}

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	for _, f := range strings.Fields(s) {
		m[f] += 1
	}

	return m
}

func StripHTML(s string) string {
	output := ""

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "</p>", " \n", -1)
		s = strings.Replace(s, "<br>", " \n", -1)
		s = strings.Replace(s, "</br>", " \n", -1)

		// Walk through the string removing all tags
		b := new(bytes.Buffer)
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}
	return output
}

func TruncateWords(s string, max int) string {
	words := strings.Fields(s)
	if max > len(words) {
		return strings.Join(words, " ")
	}

	return strings.Join(words[:max], " ")
}

func TruncateWordsToWholeSentence(s string, max int) string {
	words := strings.Fields(s)
	if max > len(words) {
		return strings.Join(words, " ")
	}

	for counter, word := range words[max:] {
		if strings.HasSuffix(word, ".") ||
			strings.HasSuffix(word, "?") ||
			strings.HasSuffix(word, ".\"") ||
			strings.HasSuffix(word, "!") {
			return strings.Join(words[:max+counter+1], " ")
		}
	}

	return strings.Join(words[:max], " ")
}

func MakePermalink(domain string, path string) string {
	return strings.TrimRight(domain, "/") + "/" + strings.TrimLeft(path, "/")
}
