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

package hqt

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cast"
)

// IsSameString asserts that two strings are equal. The two strings
// are normalized (whitespace removed) before doing a ==.
// Also note that two strings can be the same even if they're of different
// types.
var IsSameString qt.Checker = &stringChecker{
	argNames: []string{"got", "want"},
}

// IsSameType asserts that got is the same type as want.
var IsSameType qt.Checker = &typeChecker{
	argNames: []string{"got", "want"},
}

type argNames []string

func (a argNames) ArgNames() []string {
	return a
}

type typeChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got and args[0] is of the same type.
func (c *typeChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	if want := args[0]; reflect.TypeOf(got) != reflect.TypeOf(want) {
		if _, ok := got.(error); ok && want == nil {
			return errors.New("got non-nil error")
		}
		return errors.New("values are not of same type")
	}
	return nil
}

type stringChecker struct {
	argNames
}

// Check implements Checker.Check by checking that got and args[0] represents the same normalized text (whitespace etc. trimmed).
func (c *stringChecker) Check(got interface{}, args []interface{}, note func(key string, value interface{})) (err error) {
	s1, s2 := cast.ToString(got), cast.ToString(args[0])

	if s1 == s2 {
		return nil
	}

	s1, s2 = normalizeString(s1), normalizeString(s2)

	if s1 == s2 {
		return nil
	}

	return fmt.Errorf("values are not the same text: %s", htesting.DiffStrings(s1, s2))
}

func normalizeString(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	return strings.Join(lines, "\n")
}

// DeepAllowUnexported creates an option to allow compare of unexported types
// in the given list of types.
// see https://github.com/google/go-cmp/issues/40#issuecomment-328615283
func DeepAllowUnexported(vs ...interface{}) cmp.Option {
	m := make(map[reflect.Type]struct{})
	for _, v := range vs {
		structTypes(reflect.ValueOf(v), m)
	}
	var typs []interface{}
	for t := range m {
		typs = append(typs, reflect.New(t).Elem().Interface())
	}
	return cmp.AllowUnexported(typs...)
}

func structTypes(v reflect.Value, m map[reflect.Type]struct{}) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Interface:
		if !v.IsNil() {
			structTypes(v.Elem(), m)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			structTypes(v.Index(i), m)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			structTypes(v.MapIndex(k), m)
		}
	case reflect.Struct:
		m[v.Type()] = struct{}{}
		for i := 0; i < v.NumField(); i++ {
			structTypes(v.Field(i), m)
		}
	}
}
