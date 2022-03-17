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

package codegen

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/common/herrors"
)

func TestMethods(t *testing.T) {
	var (
		zeroIE     = reflect.TypeOf((*IEmbed)(nil)).Elem()
		zeroIEOnly = reflect.TypeOf((*IEOnly)(nil)).Elem()
		zeroI      = reflect.TypeOf((*I)(nil)).Elem()
	)

	dir, _ := os.Getwd()
	insp := NewInspector(dir)

	t.Run("MethodsFromTypes", func(t *testing.T) {
		c := qt.New(t)

		methods := insp.MethodsFromTypes([]reflect.Type{zeroI}, nil)

		methodsStr := fmt.Sprint(methods)

		c.Assert(methodsStr, qt.Contains, "Method1(arg0 herrors.ErrorContext)")
		c.Assert(methodsStr, qt.Contains, "Method7() interface {}")
		c.Assert(methodsStr, qt.Contains, "Method0() string\n Method4() string")
		c.Assert(methodsStr, qt.Contains, "MethodEmbed3(arg0 string) string\n MethodEmbed1() string")

		c.Assert(methods.Imports(), qt.Contains, "github.com/gohugoio/hugo/common/herrors")
	})

	t.Run("EmbedOnly", func(t *testing.T) {
		c := qt.New(t)

		methods := insp.MethodsFromTypes([]reflect.Type{zeroIEOnly}, nil)

		methodsStr := fmt.Sprint(methods)

		c.Assert(methodsStr, qt.Contains, "MethodEmbed3(arg0 string) string")
	})

	t.Run("ToMarshalJSON", func(t *testing.T) {
		c := qt.New(t)

		m, pkg := insp.MethodsFromTypes(
			[]reflect.Type{zeroI},
			[]reflect.Type{zeroIE}).ToMarshalJSON("*page", "page")

		c.Assert(m, qt.Contains, "method6 := p.Method6()")
		c.Assert(m, qt.Contains, "Method0: method0,")
		c.Assert(m, qt.Contains, "return json.Marshal(&s)")

		c.Assert(pkg, qt.Contains, "github.com/gohugoio/hugo/common/herrors")
		c.Assert(pkg, qt.Contains, "encoding/json")

		fmt.Println(pkg)
	})
}

type I interface {
	IEmbed
	Method0() string
	Method4() string
	Method1(myerr herrors.ErrorContext)
	Method3(myint int, mystring string)
	Method5() (string, error)
	Method6() *net.IP
	Method7() any
	Method8() herrors.ErrorContext
	method2()
	method9() os.FileInfo
}

type IEOnly interface {
	IEmbed
}
