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

// Package lang provides template functions for content internationalization.
package lang

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

// New returns a new instance of the lang-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "lang" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Translate returns a translated string for id.
func (ns *Namespace) Translate(id interface{}, args ...interface{}) (string, error) {
	sid, err := cast.ToStringE(id)
	if err != nil {
		return "", nil
	}

	return ns.deps.Translate(sid, args...), nil
}

// NumFmt formats a number with the given precision using the
// negative, decimal, and grouping options.  The `options`
// parameter is a string consisting of `<negative> <decimal> <grouping>`.  The
// default `options` value is `- . ,`.
//
// Note that numbers are rounded up at 5 or greater.
// So, with precision set to 0, 1.5 becomes `2`, and 1.4 becomes `1`.
func (ns *Namespace) NumFmt(precision, number interface{}, options ...interface{}) (string, error) {
	prec, err := cast.ToIntE(precision)
	if err != nil {
		return "", err
	}

	n, err := cast.ToFloat64E(number)
	if err != nil {
		return "", err
	}

	var neg, dec, grp string

	if len(options) == 0 {
		// defaults
		neg, dec, grp = "-", ".", ","
	} else {
		delim := " "

		if len(options) == 2 {
			// custom delimiter
			s, err := cast.ToStringE(options[1])
			if err != nil {
				return "", nil
			}

			delim = s
		}

		s, err := cast.ToStringE(options[0])
		if err != nil {
			return "", nil
		}

		rs := strings.Split(s, delim)
		switch len(rs) {
		case 0:
		case 1:
			neg = rs[0]
		case 2:
			neg, dec = rs[0], rs[1]
		case 3:
			neg, dec, grp = rs[0], rs[1], rs[2]
		default:
			return "", errors.New("too many fields in options parameter to NumFmt")
		}
	}

	// Logic from MIT Licensed github.com/go-playground/locales/
	// Original Copyright (c) 2016 Go Playground

	s := strconv.FormatFloat(math.Abs(n), 'f', prec, 64)
	L := len(s) + 2 + len(s[:len(s)-1-prec])/3

	var count int
	inWhole := prec == 0
	b := make([]byte, 0, L)

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			for j := len(dec) - 1; j >= 0; j-- {
				b = append(b, dec[j])
			}
			inWhole = true
			continue
		}

		if inWhole {
			if count == 3 {
				for j := len(grp) - 1; j >= 0; j-- {
					b = append(b, grp[j])
				}
				count = 1
			} else {
				count++
			}
		}

		b = append(b, s[i])
	}

	if n < 0 {
		for j := len(neg) - 1; j >= 0; j-- {
			b = append(b, neg[j])
		}
	}

	// reverse
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b), nil
}

type pagesLanguageMerger interface {
	MergeByLanguageInterface(other interface{}) (interface{}, error)
}

// Merge creates a union of pages from two languages.
func (ns *Namespace) Merge(p2, p1 interface{}) (interface{}, error) {
	merger, ok := p1.(pagesLanguageMerger)
	if !ok {
		return nil, fmt.Errorf("language merge not supported for %T", p1)
	}
	return merger.MergeByLanguageInterface(p2)
}
