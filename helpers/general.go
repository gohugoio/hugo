// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/spf13/afero"

	"github.com/jdkato/prose/transform"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// TCPListen starts listening on a valid TCP port.
func TCPListen() (net.Listener, *net.TCPAddr, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}
	addr := l.Addr()
	if a, ok := addr.(*net.TCPAddr); ok {
		return l, a, nil
	}
	l.Close()
	return nil, nil, fmt.Errorf("unable to obtain a valid tcp port: %v", addr)
}

// FirstUpper returns a string with the first character as upper case.
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// UniqueStrings returns a new slice with any duplicates removed.
func UniqueStrings(s []string) []string {
	unique := make([]string, 0, len(s))
	for i, val := range s {
		var seen bool
		for j := 0; j < i; j++ {
			if s[j] == val {
				seen = true
				break
			}
		}
		if !seen {
			unique = append(unique, val)
		}
	}
	return unique
}

// UniqueStringsReuse returns a slice with any duplicates removed.
// It will modify the input slice.
func UniqueStringsReuse(s []string) []string {
	result := s[:0]
	for i, val := range s {
		var seen bool

		for j := 0; j < i; j++ {
			if s[j] == val {
				seen = true
				break
			}
		}

		if !seen {
			result = append(result, val)
		}
	}
	return result
}

// UniqueStringsSorted returns a sorted slice with any duplicates removed.
// It will modify the input slice.
func UniqueStringsSorted(s []string) []string {
	if len(s) == 0 {
		return nil
	}
	ss := sort.StringSlice(s)
	ss.Sort()
	i := 0
	for j := 1; j < len(s); j++ {
		if !ss.Less(i, j) {
			continue
		}
		i++
		s[i] = s[j]
	}

	return s[:i+1]
}

// ReaderToBytes takes an io.Reader argument, reads from it
// and returns bytes.
func ReaderToBytes(lines io.Reader) []byte {
	if lines == nil {
		return []byte{}
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)

	b.ReadFrom(lines)

	bc := make([]byte, b.Len())
	copy(bc, b.Bytes())
	return bc
}

// ReaderToString is the same as ReaderToBytes, but returns a string.
func ReaderToString(lines io.Reader) string {
	if lines == nil {
		return ""
	}
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	b.ReadFrom(lines)
	return b.String()
}

// ReaderContains reports whether subslice is within r.
func ReaderContains(r io.Reader, subslice []byte) bool {
	if r == nil || len(subslice) == 0 {
		return false
	}

	bufflen := len(subslice) * 4
	halflen := bufflen / 2
	buff := make([]byte, bufflen)
	var err error
	var n, i int

	for {
		i++
		if i == 1 {
			n, err = io.ReadAtLeast(r, buff[:halflen], halflen)
		} else {
			if i != 2 {
				// shift left to catch overlapping matches
				copy(buff[:], buff[halflen:])
			}
			n, err = io.ReadAtLeast(r, buff[halflen:], halflen)
		}

		if n > 0 && bytes.Contains(buff, subslice) {
			return true
		}

		if err != nil {
			break
		}
	}
	return false
}

// GetTitleFunc returns a func that can be used to transform a string to
// title case.
//
// # The supported styles are
//
// - "Go" (strings.Title)
// - "AP" (see https://www.apstylebook.com/)
// - "Chicago" (see https://www.chicagomanualofstyle.org/home.html)
// - "FirstUpper" (only the first character is upper case)
// - "None" (no transformation)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
	switch strings.ToLower(style) {
	case "go":
		//lint:ignore SA1019 keep for now.
		return strings.Title
	case "chicago":
		tc := transform.NewTitleConverter(transform.ChicagoStyle)
		return tc.Title
	case "none":
		return func(s string) string { return s }
	case "firstupper":
		return FirstUpper
	default:
		tc := transform.NewTitleConverter(transform.APStyle)
		return tc.Title
	}
}

// HasStringsPrefix tests whether the string slice s begins with prefix slice s.
func HasStringsPrefix(s, prefix []string) bool {
	return len(s) >= len(prefix) && compareStringSlices(s[0:len(prefix)], prefix)
}

// HasStringsSuffix tests whether the string slice s ends with suffix slice s.
func HasStringsSuffix(s, suffix []string) bool {
	return len(s) >= len(suffix) && compareStringSlices(s[len(s)-len(suffix):], suffix)
}

func compareStringSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// SliceToLower goes through the source slice and lowers all values.
func SliceToLower(s []string) []string {
	if s == nil {
		return nil
	}

	l := make([]string, len(s))
	for i, v := range s {
		l[i] = strings.ToLower(v)
	}

	return l
}

// IsWhitespace determines if the given rune is whitespace.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// PrintFs prints the given filesystem to the given writer starting from the given path.
// This is useful for debugging.
func PrintFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}

	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(fmt.Sprintf("error: path %q: %s", path, err))
		}
		path = filepath.ToSlash(path)
		if path == "" {
			path = "."
		}
		fmt.Fprintln(w, path, info.IsDir())
		return nil
	})
}

// FormatByteCount pretty formats b.
func FormatByteCount(bc uint64) string {
	const (
		Gigabyte = 1 << 30
		Megabyte = 1 << 20
		Kilobyte = 1 << 10
	)
	switch {
	case bc > Gigabyte || -bc > Gigabyte:
		return fmt.Sprintf("%.2f GB", float64(bc)/Gigabyte)
	case bc > Megabyte || -bc > Megabyte:
		return fmt.Sprintf("%.2f MB", float64(bc)/Megabyte)
	case bc > Kilobyte || -bc > Kilobyte:
		return fmt.Sprintf("%.2f KB", float64(bc)/Kilobyte)
	}
	return fmt.Sprintf("%d B", bc)
}
