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

	_strings "strings"
	"unicode/utf8"

	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
)

// New returns a new instance of the strings-namespaced template functions.
func New(d *deps.Deps) *Namespace {
	titleCaseStyle := d.Cfg.GetString("titleCaseStyle")
	titleFunc := helpers.GetTitleFunc(titleCaseStyle)
	return &Namespace{deps: d, titleFunc: titleFunc}
}

// Namespace provides template functions for the "strings" namespace.
// Most functions mimic the Go stdlib, but the order of the parameters may be
// different to ease their use in the Go template system.
type Namespace struct {
	titleFunc func(s string) string
	deps      *deps.Deps
}

// CountRunes returns the number of runes in s, excluding whitepace.
func (ns *Namespace) CountRunes(s interface{}) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, _errors.Wrap(err, "Failed to convert content to string")
	}

	counter := 0
	for _, r := range helpers.StripHTML(ss) {
		if !helpers.IsWhitespace(r) {
			counter++
		}
	}

	return counter, nil
}

// RuneCount returns the number of runes in s.
func (ns *Namespace) RuneCount(s interface{}) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, _errors.Wrap(err, "Failed to convert content to string")
	}
	return utf8.RuneCountInString(ss), nil
}

// CountWords returns the approximate word count in s.
func (ns *Namespace) CountWords(s interface{}) (int, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, _errors.Wrap(err, "Failed to convert content to string")
	}

	counter := 0
	for _, word := range _strings.Fields(helpers.StripHTML(ss)) {
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			counter++
		} else {
			counter += runeCount
		}
	}

	return counter, nil
}

// Chomp returns a copy of s with all trailing newline characters removed.
func (ns *Namespace) Chomp(s interface{}) (interface{}, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	res := _strings.TrimRight(ss, "\r\n")
	switch s.(type) {
	case template.HTML:
		return template.HTML(res), nil
	default:
		return res, nil
	}

}

// Contains reports whether substr is in s.
func (ns *Namespace) Contains(s, substr interface{}) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	su, err := cast.ToStringE(substr)
	if err != nil {
		return false, err
	}

	return _strings.Contains(ss, su), nil
}

// ContainsAny reports whether any Unicode code points in chars are within s.
func (ns *Namespace) ContainsAny(s, chars interface{}) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sc, err := cast.ToStringE(chars)
	if err != nil {
		return false, err
	}

	return _strings.ContainsAny(ss, sc), nil
}

// HasPrefix tests whether the input s begins with prefix.
func (ns *Namespace) HasPrefix(s, prefix interface{}) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sx, err := cast.ToStringE(prefix)
	if err != nil {
		return false, err
	}

	return _strings.HasPrefix(ss, sx), nil
}

// HasSuffix tests whether the input s begins with suffix.
func (ns *Namespace) HasSuffix(s, suffix interface{}) (bool, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return false, err
	}

	sx, err := cast.ToStringE(suffix)
	if err != nil {
		return false, err
	}

	return _strings.HasSuffix(ss, sx), nil
}

// Replace returns a copy of the string s with all occurrences of old replaced
// with new.
func (ns *Namespace) Replace(s, old, new interface{}) (string, error) {
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

	return _strings.Replace(ss, so, sn, -1), nil
}

// SliceString slices a string by specifying a half-open range with
// two indices, start and end. 1 and 4 creates a slice including elements 1 through 3.
// The end index can be omitted, it defaults to the string's length.
func (ns *Namespace) SliceString(a interface{}, startEnd ...interface{}) (string, error) {
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
func (ns *Namespace) Split(a interface{}, delimiter string) ([]string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return []string{}, err
	}

	return _strings.Split(aStr, delimiter), nil
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
func (ns *Namespace) Substr(a interface{}, nums ...interface{}) (string, error) {
	aStr, err := cast.ToStringE(a)
	if err != nil {
		return "", err
	}

	var start, length int

	asRunes := []rune(aStr)

	switch len(nums) {
	case 0:
		return "", errors.New("too less arguments")
	case 1:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
		length = len(asRunes)
	case 2:
		if start, err = cast.ToIntE(nums[0]); err != nil {
			return "", errors.New("start argument must be integer")
		}
		if length, err = cast.ToIntE(nums[1]); err != nil {
			return "", errors.New("length argument must be integer")
		}
	default:
		return "", errors.New("too many arguments")
	}

	if start < -len(asRunes) {
		start = 0
	}
	if start > len(asRunes) {
		return "", fmt.Errorf("start position out of bounds for %d-byte string", len(aStr))
	}

	var s, e int
	if start >= 0 && length >= 0 {
		s = start
		e = start + length
	} else if start < 0 && length >= 0 {
		s = len(asRunes) + start - length + 1
		e = len(asRunes) + start + 1
	} else if start >= 0 && length < 0 {
		s = start
		e = len(asRunes) + length
	} else {
		s = len(asRunes) + start
		e = len(asRunes) + length
	}

	if s > e {
		return "", fmt.Errorf("calculated start position greater than end position: %d > %d", s, e)
	}
	if e > len(asRunes) {
		e = len(asRunes)
	}

	return string(asRunes[s:e]), nil
}

// Title returns a copy of the input s with all Unicode letters that begin words
// mapped to their title case.
func (ns *Namespace) Title(s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return ns.titleFunc(ss), nil
}

// FirstUpper returns a string with the first character as upper case.
func (ns *Namespace) FirstUpper(s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return helpers.FirstUpper(ss), nil
}

// ToLower returns a copy of the input s with all Unicode letters mapped to their
// lower case.
func (ns *Namespace) ToLower(s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return _strings.ToLower(ss), nil
}

// ToUpper returns a copy of the input s with all Unicode letters mapped to their
// upper case.
func (ns *Namespace) ToUpper(s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return _strings.ToUpper(ss), nil
}

// Trim returns a string with all leading and trailing characters defined
// contained in cutset removed.
func (ns *Namespace) Trim(s, cutset interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return _strings.Trim(ss, sc), nil
}

// TrimLeft returns a slice of the string s with all leading characters
// contained in cutset removed.
func (ns *Namespace) TrimLeft(cutset, s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return _strings.TrimLeft(ss, sc), nil
}

// TrimPrefix returns s without the provided leading prefix string. If s doesn't
// start with prefix, s is returned unchanged.
func (ns *Namespace) TrimPrefix(prefix, s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sx, err := cast.ToStringE(prefix)
	if err != nil {
		return "", err
	}

	return _strings.TrimPrefix(ss, sx), nil
}

// TrimRight returns a slice of the string s with all trailing characters
// contained in cutset removed.
func (ns *Namespace) TrimRight(cutset, s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sc, err := cast.ToStringE(cutset)
	if err != nil {
		return "", err
	}

	return _strings.TrimRight(ss, sc), nil
}

// TrimSuffix returns s without the provided trailing suffix string. If s
// doesn't end with suffix, s is returned unchanged.
func (ns *Namespace) TrimSuffix(suffix, s interface{}) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	sx, err := cast.ToStringE(suffix)
	if err != nil {
		return "", err
	}

	return _strings.TrimSuffix(ss, sx), nil
}

// Repeat returns a new string consisting of count copies of the string s.
func (ns *Namespace) Repeat(n, s interface{}) (string, error) {
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

	return _strings.Repeat(ss, sn), nil
}
