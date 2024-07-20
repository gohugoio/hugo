// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package strings provides template functions for manipulating strings.
package strings

import (
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gohugoio/hugo/common/text"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/tpl"
	"github.com/rogpeppe/go-internal/diff"

	"github.com/spf13/cast"
)

// New returns a new instance of the strings-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	return &Namespace{deps: d}
}

// Namespace provides template functions for the "strings" namespace.
// Most functions mimic the Go stdlib, but the order of the parameters may be
// different to ease their use in the Go template system.
type Namespace struct {
	deps *deps.Deps
}

// CountRunes returns the number of runes in s, excluding whitespace.
func (ns *Namespace) CountRunes(s any) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert content to string: %w", err)
	}

	counter := 0
	for _, r := range tpl.StripHTML(ss) {
		if !helpers.IsWhitespace(r) {
			counter++
		}
	}

	return counter, nil
}

// RuneCount returns the number of runes in s.
func (ns *Namespace) RuneCount(s any) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert content to string: %w", err)
	}
	return utf8.RuneCountInString(ss), nil
}

// CountWords returns the approximate word count in s.
func (ns *Namespace) CountWords(s any) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert content to string: %w", err)
	}

	isCJKLanguage, err := regexp.MatchString(`\p{Han}|\p{Hangul}|\p{Hiragana}|\p{Katakana}`, ss)
	if err != nil {
		return 0, fmt.Errorf("failed to match regex pattern against string: %w", err)
	}

	if !isCJKLanguage {
		return len(strings.Fields(tpl.StripHTML(ss))), nil
	}

	counter := 0
	for _, word := range strings.Fields(tpl.StripHTML(ss)) {
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			counter++
		} else {
			counter += runeCount
		}
	}

	return counter, nil
}

// Count counts the number of non-overlapping instances of substr in s.
// If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func (ns *Namespace) Count(substr, s any) (int, error) {
	substrs, err := cast.ToStringE(substr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert substr to string: %w", err)
	}
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert s to string: %w", err)
	}
	return strings.Count(ss, substrs), nil
}

// Chomp returns a copy of s with all trailing newline characters removed.
func (ns *Namespace) Chomp(s any) (any, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	res := text.Chomp(ss)
	switch s.(type) {
	case template.HTML:
		return template.HTML(res), nil
	default:
		return res, nil
	}
}

// Contains reports whether substr is in s.
func (ns *Namespace) Contains(s, substr any) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	su, err := cast.ToStringE(substr)
	if err != nil {
		return false, err
	}

	return strings.Contains(ss, su), nil
}

// ContainsAny reports whether any Unicode code points in chars are within s.
func (ns *Namespace) ContainsAny(s, chars any) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sc, err := cast.ToStringE(chars)
	if err != nil {
		return false, err
	}

	return strings.ContainsAny(ss, sc), nil
}

// ContainsNonSpace reports whether s contains any non-space characters as defined
// by Unicode's White Space property,
// <docsmeta>{"newIn": "0.111.0" }</docsmeta>
func (ns *Namespace) ContainsNonSpace(s any) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	for _, r := range ss {
		if !unicode.IsSpace(r) {
			return true, nil
		}
	}
	return false, nil
}

// Diff returns an anchored diff of the two texts old and new in the “unified
// diff” format. If old and new are identical, Diff returns an empty string.
func (ns *Namespace) Diff(oldname string, old any, newname string, new any) (string, error) {
	olds, err := cast.ToStringE(old)
	if err != nil {
		return "", err
	}
	news, err := cast.ToStringE(new)
	if err != nil {
		return "", err
	}
	return string(diff.Diff(oldname, []byte(olds), newname, []byte(news))), nil
}

// HasPrefix tests whether the input s begins with prefix.
func (ns *Namespace) HasPrefix(s, prefix any) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sx, err := cast.ToStringE(prefix)
	if err != nil {
		return false, err
	}

	return strings.HasPrefix(ss, sx), nil
}

// HasSuffix tests whether the input s begins with suffix.
func (ns *Namespace) HasSuffix(s, suffix any) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sx, err := cast.ToStringE(suffix)
	if err != nil {
		return false, err
	}

	return strings.HasSuffix(ss, sx), nil
}

// Replace returns a copy of the string s with all occurrences of old replaced
// with new.  The number of replacements can be limited with an optional fourth
// parameter.
func (ns *Namespace) Replace(s, old, new any, limit ...any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	so, err := cast.ToStringE(old)
	if err != nil {
		return "", err
	}

	sn, err := cast.ToStringE(new)
	if err != nil {
		return "", err
	}

	if len(limit) == 0 {
		return strings.ReplaceAll(ss, so, sn), nil
	}

	lim, err := cast.ToIntE(limit[0])
	if err != nil {
		return "", err
	}

	return strings.Replace(ss, so, sn, lim), nil
}

// SliceString slices a string by specifying a half-open range with
// two indices, start and end. 1 and 4 creates a slice including elements 1 through 3.
// The end index can be omitted, it defaults to the string's length.
func (ns *Namespace) SliceString(a any, startEnd ...any) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	var argStart, argEnd int

	argNum := len(startEnd)

	if argNum > 0 {
		if argStart, err = cast.ToIntE(startEnd[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
	}
	if argNum > 1 {
		if argEnd, err = cast.ToIntE(startEnd[1]); err != nil {
			return "", errors.New("end argument must be integer")
		}
	}

	if argNum > 2 {
		return "", errors.New("too many arguments")
	}

	asRunes := []rune(aStr)

	if argNum > 0 && (argStart < 0 || argStart >= len(asRunes)) {
		return "", errors.New("slice bounds out of range")
	}

	if argNum == 2 {
		if argEnd < 0 || argEnd > len(asRunes) {
			return "", errors.New("slice bounds out of range")
		}
		return string(asRunes[argStart:argEnd]), nil
	} else if argNum == 1 {
		return string(asRunes[argStart:]), nil
	} else {
		return string(asRunes[:]), nil
	}
}

// Split slices an input string into all substrings separated by delimiter.
func (ns *Namespace) Split(a any, delimiter string) ([]string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(aStr, delimiter), nil
}

// Substr extracts parts of a string, beginning at the character at the specified
// position, and returns the specified number of characters.
//
// It normally takes two parameters: start and length.
// It can also take one parameter: start, i.e. length is omitted, in which case
// the substring starting from start until the end of the string will be returned.
//
// To extract characters from the end of the string, use a negative start number.
//
// In addition, borrowing from the extended behavior described at http://php.net/substr,
// if length is given and is negative, then that many characters will be omitted from
// the end of string.
func (ns *Namespace) Substr(a any, nums ...any) (string, error) {
	s, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	asRunes := []rune(s)
	rlen := len(asRunes)

	var start, length int

	switch len(nums) {
	case 0:
		return "", errors.New("too few arguments")
	case 1:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be an integer")
		}
		length = rlen
	case 2:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be an integer")
		}
		if length, err = cast.ToIntE(nums[1]); err != nil {
			return "", errors.New("length argument must be an integer")
		}
	default:
		return "", errors.New("too many arguments")
	}

	if rlen == 0 {
		return "", nil
	}

	if start < 0 {
		start += rlen
	}

	// start was originally negative beyond rlen
	if start < 0 {
		start = 0
	}

	if start > rlen-1 {
		return "", nil
	}

	end := rlen

	switch {
	case length == 0:
		return "", nil
	case length < 0:
		end += length
	case length > 0:
		end = start + length
	}

	if start >= end {
		return "", nil
	}

	if end < 0 {
		return "", nil
	}

	if end > rlen {
		end = rlen
	}

	return string(asRunes[start:end]), nil
}

// Title returns a copy of the input s with all Unicode letters that begin words
// mapped to their title case.
func (ns *Namespace) Title(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}
	return ns.deps.Conf.CreateTitle(ss), nil
}

// FirstUpper converts s making  the first character upper case.
func (ns *Namespace) FirstUpper(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return helpers.FirstUpper(ss), nil
}

// ToLower returns a copy of the input s with all Unicode letters mapped to their
// lower case.
func (ns *Namespace) ToLower(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return strings.ToLower(ss), nil
}

// ToUpper returns a copy of the input s with all Unicode letters mapped to their
// upper case.
func (ns *Namespace) ToUpper(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(ss), nil
}

// Trim returns converts the strings s removing all leading and trailing characters defined
// contained.
func (ns *Namespace) Trim(s, cutset any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return strings.Trim(ss, sc), nil
}

// TrimLeft returns a slice of the string s with all leading characters
// contained in cutset removed.
func (ns *Namespace) TrimLeft(cutset, s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return strings.TrimLeft(ss, sc), nil
}

// TrimPrefix returns s without the provided leading prefix string. If s doesn't
// start with prefix, s is returned unchanged.
func (ns *Namespace) TrimPrefix(prefix, s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sx, err := cast.ToStringE(prefix)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(ss, sx), nil
}

// TrimRight returns a slice of the string s with all trailing characters
// contained in cutset removed.
func (ns *Namespace) TrimRight(cutset, s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(ss, sc), nil
}

// TrimSuffix returns s without the provided trailing suffix string. If s
// doesn't end with suffix, s is returned unchanged.
func (ns *Namespace) TrimSuffix(suffix, s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sx, err := cast.ToStringE(suffix)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(ss, sx), nil
}

// Repeat returns a new string consisting of n copies of the string s.
func (ns *Namespace) Repeat(n, s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sn, err := cast.ToIntE(n)
	if err != nil {
		return "", err
	}

	if sn < 0 {
		return "", errors.New("strings: negative Repeat count")
	}

	return strings.Repeat(ss, sn), nil
}
