// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/afero"

	"github.com/jdkato/prose/transform"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// Strips carriage returns from third-party / external processes (useful for Windows)
func normalizeExternalHelperLineFeeds(content []byte) []byte {
	return bytes.Replace(content, []byte("\r"), []byte(""), -1)
}

// FindAvailablePort returns an available and valid TCP port.
func FindAvailablePort() (*net.TCPAddr, error) {
	l, err := net.Listen("tcp", ":0")
	if err == nil {
		defer l.Close()
		addr := l.Addr()
		if a, ok := addr.(*net.TCPAddr); ok {
			return a, nil
		}
		return nil, fmt.Errorf("Unable to obtain a valid tcp port. %v", addr)
	}
	return nil, err
}

// InStringArray checks if a string is an element of a slice of strings
// and returns a boolean value.
func InStringArray(arr []string, el string) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

// GuessType attempts to guess the type of file from a given string.
func GuessType(in string) string {
	switch strings.ToLower(in) {
	case "md", "markdown", "mdown":
		return "markdown"
	case "asciidoc", "adoc", "ad":
		return "asciidoc"
	case "mmark":
		return "mmark"
	case "rst":
		return "rst"
	case "pandoc", "pdc":
		return "pandoc"
	case "html", "htm":
		return "html"
	case "org":
		return "org"
	}

	return "unknown"
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
	var unique []string
	set := map[string]interface{}{}
	for _, val := range s {
		if _, ok := set[val]; !ok {
			unique = append(unique, val)
			set[val] = val
		}
	}
	return unique
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

	bc := make([]byte, b.Len(), b.Len())
	copy(bc, b.Bytes())
	return bc
}

// ToLowerMap makes all the keys in the given map lower cased and will do so
// recursively.
// Notes:
// * This will modify the map given.
// * Any nested map[interface{}]interface{} will be converted to map[string]interface{}.
func ToLowerMap(m map[string]interface{}) {
	for k, v := range m {
		switch v.(type) {
		case map[interface{}]interface{}:
			v = cast.ToStringMap(v)
			ToLowerMap(v.(map[string]interface{}))
		case map[string]interface{}:
			ToLowerMap(v.(map[string]interface{}))
		}

		lKey := strings.ToLower(k)
		if k != lKey {
			delete(m, k)
			m[lKey] = v
		}

	}
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
// The supported styles are
//
// - "Go" (strings.Title)
// - "AP" (see https://www.apstylebook.com/)
// - "Chicago" (see http://www.chicagomanualofstyle.org/home.html)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
	switch strings.ToLower(style) {
	case "go":
		return strings.Title
	case "chicago":
		tc := transform.NewTitleConverter(transform.ChicagoStyle)
		return tc.Title
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

// ThemeSet checks whether a theme is in use or not.
func (p *PathSpec) ThemeSet() bool {
	return p.theme != ""
}

type logPrinter interface {
	// Println is the only common method that works in all of JWWs loggers.
	Println(a ...interface{})
}

// DistinctLogger ignores duplicate log statements.
type DistinctLogger struct {
	sync.RWMutex
	logger logPrinter
	m      map[string]bool
}

// Println will log the string returned from fmt.Sprintln given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Println(v ...interface{}) {
	// fmt.Sprint doesn't add space between string arguments
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.print(logStatement)
}

// Printf will log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
// Note: A newline is appended.
func (l *DistinctLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.print(logStatement)
}

func (l *DistinctLogger) print(logStatement string) {
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	if !l.m[logStatement] {
		l.logger.Println(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}

// NewDistinctErrorLogger creates a new DistinctLogger that logs ERRORs
func NewDistinctErrorLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.ERROR}
}

// NewDistinctWarnLogger creates a new DistinctLogger that logs WARNs
func NewDistinctWarnLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.WARN}
}

// NewDistinctFeedbackLogger creates a new DistinctLogger that can be used
// to give feedback to the user while not spamming with duplicates.
func NewDistinctFeedbackLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.FEEDBACK}
}

var (
	// DistinctErrorLog can be used to avoid spamming the logs with errors.
	DistinctErrorLog = NewDistinctErrorLogger()

	// DistinctWarnLog can be used to avoid spamming the logs with warnings.
	DistinctWarnLog = NewDistinctWarnLogger()

	// DistinctFeedbackLog can be used to avoid spamming the logs with info messages.
	DistinctFeedbackLog = NewDistinctFeedbackLogger()
)

// InitLoggers sets up the global distinct loggers.
func InitLoggers() {
	DistinctErrorLog = NewDistinctErrorLogger()
	DistinctWarnLog = NewDistinctWarnLogger()
	DistinctFeedbackLog = NewDistinctFeedbackLogger()
}

// Deprecated informs about a deprecation, but only once for a given set of arguments' values.
// If the err flag is enabled, it logs as an ERROR (will exit with -1) and the text will
// point at the next Hugo release.
// The idea is two remove an item in two Hugo releases to give users and theme authors
// plenty of time to fix their templates.
func Deprecated(object, item, alternative string, err bool) {
	if err {
		DistinctErrorLog.Printf("%s's %s is deprecated and will be removed in Hugo %s. %s.", object, item, CurrentHugoVersion.Next().ReleaseVersion(), alternative)

	} else {
		// Make sure the users see this while avoiding build breakage. This will not lead to an os.Exit(-1)
		DistinctFeedbackLog.Printf("WARNING: %s's %s is deprecated and will be removed in a future release. %s.", object, item, alternative)
	}
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

// MD5String takes a string and returns its MD5 hash.
func MD5String(f string) string {
	h := md5.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum([]byte{}))
}

// MD5FromFileFast creates a MD5 hash from the given file. It only reads parts of
// the file for speed, so don't use it if the files are very subtly different.
// It will not close the file.
func MD5FromFileFast(f afero.File) (string, error) {
	const (
		// Do not change once set in stone!
		maxChunks = 8
		peekSize  = 64
		seek      = 2048
	)

	h := md5.New()
	buff := make([]byte, peekSize)

	for i := 0; i < maxChunks; i++ {
		if i > 0 {
			_, err := f.Seek(seek, 0)
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", err
			}
		}

		_, err := io.ReadAtLeast(f, buff, peekSize)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				h.Write(buff)
				break
			}
			return "", err
		}
		h.Write(buff)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// MD5FromFile creates a MD5 hash from the given file.
// It will not close the file.
func MD5FromFile(f afero.File) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// IsWhitespace determines if the given rune is whitespace.
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

// NormalizeHugoFlags facilitates transitions of Hugo command-line flags,
// e.g. --baseUrl to --baseURL, --uglyUrls to --uglyURLs
func NormalizeHugoFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "baseUrl":
		name = "baseURL"
		break
	case "uglyUrls":
		name = "uglyURLs"
		break
	}
	return pflag.NormalizedName(name)
}

// DiffStringSlices returns the difference between two string slices.
// Useful in tests.
// See:
// http://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings-in-golang
func DiffStringSlices(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}

// DiffString splits the strings into fields and runs it into DiffStringSlices.
// Useful for tests.
func DiffStrings(s1, s2 string) []string {
	return DiffStringSlices(strings.Fields(s1), strings.Fields(s2))
}
