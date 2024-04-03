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
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gohugoio/locales"
	translators "github.com/gohugoio/localescompressed"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/deps"
	"github.com/spf13/cast"
)

// New returns a new instance of the lang-namespaced template functions.
func New(deps *deps.Deps, translator locales.Translator) *Namespace {
	return &Namespace{
		translator: translator,
		deps:       deps,
	}
}

// Namespace provides template functions for the "lang" namespace.
type Namespace struct {
	translator locales.Translator
	deps       *deps.Deps
}

// Translate returns a translated string for id.
func (ns *Namespace) Translate(ctx context.Context, id any, args ...any) (string, error) {
	var templateData any

	if len(args) > 0 {
		if len(args) > 1 {
			return "", fmt.Errorf("wrong number of arguments, expecting at most 2, got %d", len(args)+1)
		}
		templateData = args[0]
	}

	sid, err := cast.ToStringE(id)
	if err != nil {
		return "", err
	}

	return ns.deps.Translate(ctx, sid, templateData), nil
}

// FormatNumber formats number with the given precision for the current language.
func (ns *Namespace) FormatNumber(precision, number any) (string, error) {
	p, n, err := ns.castPrecisionNumber(precision, number)
	if err != nil {
		return "", err
	}
	return ns.translator.FmtNumber(n, p), nil
}

// FormatPercent formats number with the given precision for the current language.
// Note that the number is assumed to be a percentage.
func (ns *Namespace) FormatPercent(precision, number any) (string, error) {
	p, n, err := ns.castPrecisionNumber(precision, number)
	if err != nil {
		return "", err
	}
	return ns.translator.FmtPercent(n, p), nil
}

// FormatCurrency returns the currency representation of number for the given currency and precision
// for the current language.
//
// The return value is formatted with at least two decimal places.
func (ns *Namespace) FormatCurrency(precision, currency, number any) (string, error) {
	p, n, err := ns.castPrecisionNumber(precision, number)
	if err != nil {
		return "", err
	}
	c := translators.GetCurrency(cast.ToString(currency))
	if c < 0 {
		return "", fmt.Errorf("unknown currency code: %q", currency)
	}
	return ns.translator.FmtCurrency(n, p, c), nil
}

// FormatAccounting returns the currency representation of number for the given currency and precision
// for the current language in accounting notation.
//
// The return value is formatted with at least two decimal places.
func (ns *Namespace) FormatAccounting(precision, currency, number any) (string, error) {
	p, n, err := ns.castPrecisionNumber(precision, number)
	if err != nil {
		return "", err
	}
	c := translators.GetCurrency(cast.ToString(currency))
	if c < 0 {
		return "", fmt.Errorf("unknown currency code: %q", currency)
	}
	return ns.translator.FmtAccounting(n, p, c), nil
}

func (ns *Namespace) castPrecisionNumber(precision, number any) (uint64, float64, error) {
	p, err := cast.ToUint64E(precision)
	if err != nil {
		return 0, 0, err
	}

	// Sanity check.
	if p > 20 {
		return 0, 0, fmt.Errorf("invalid precision: %d", precision)
	}

	n, err := cast.ToFloat64E(number)
	if err != nil {
		return 0, 0, err
	}
	return p, n, nil
}

// FormatNumberCustom formats a number with the given precision. The first
// options parameter is a space-delimited string of characters to represent
// negativity, the decimal point, and grouping. The default value is `- . ,`.
// The second options parameter defines an alternate delimiting character.
//
// Note that numbers are rounded up at 5 or greater.
// So, with precision set to 0, 1.5 becomes `2`, and 1.4 becomes `1`.
//
// For a simpler function that adapts to the current language, see FormatNumber.
func (ns *Namespace) FormatNumberCustom(precision, number any, options ...any) (string, error) {
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
				return "", err
			}

			delim = s
		}

		s, err := cast.ToStringE(options[0])
		if err != nil {
			return "", err
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

	exp := math.Pow(10.0, float64(prec))
	r := math.Round(n*exp) / exp

	// Logic from MIT Licensed github.com/gohugoio/locales/
	// Original Copyright (c) 2016 Go Playground

	s := strconv.FormatFloat(math.Abs(r), 'f', prec, 64)
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

// Deprecated: Use lang.FormatNumberCustom instead.
func (ns *Namespace) NumFmt(precision, number any, options ...any) (string, error) {
	hugo.Deprecate("lang.NumFmt", "Use lang.FormatNumberCustom instead.", "v0.120.0")
	return ns.FormatNumberCustom(precision, number, options...)
}

type pagesLanguageMerger interface {
	MergeByLanguageInterface(other any) (any, error)
}

// Merge creates a union of pages from two languages.
func (ns *Namespace) Merge(p2, p1 any) (any, error) {
	if !hreflect.IsTruthful(p1) {
		return p2, nil
	}
	if !hreflect.IsTruthful(p2) {
		return p1, nil
	}
	merger, ok := p1.(pagesLanguageMerger)
	if !ok {
		return nil, fmt.Errorf("language merge not supported for %T", p1)
	}
	return merger.MergeByLanguageInterface(p2)
}
