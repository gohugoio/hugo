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

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/stretchr/testify/require"
)

func TestMethods(t *testing.T) {

	var (
		zeroIE     = reflect.TypeOf((*IEmbed)(nil)).Elem()
		zeroIEOnly = reflect.TypeOf((*IEOnly)(nil)).Elem()
		zeroI      = reflect.TypeOf((*I)(nil)).Elem()
	)

	dir, _ := os.Getwd()
	c := NewInspector(dir)

	t.Run("MethodsFromTypes", func(t *testing.T) {
		assert := require.New(t)

		methods := c.MethodsFromTypes([]reflect.Type{zeroI}, nil)

		methodsStr := fmt.Sprint(methods)

		assert.Contains(methodsStr, "Method1(arg0 herrors.ErrorContext)")
		assert.Contains(methodsStr, "Method7() interface {}")
		assert.Contains(methodsStr, "Method0() string\n Method4() string")
		assert.Contains(methodsStr, "MethodEmbed3(arg0 string) string\n MethodEmbed1() string")

		assert.Contains(methods.Imports(), "github.com/gohugoio/hugo/common/herrors")
	})

	t.Run("EmbedOnly", func(t *testing.T) {
		assert := require.New(t)

		methods := c.MethodsFromTypes([]reflect.Type{zeroIEOnly}, nil)

		methodsStr := fmt.Sprint(methods)

		assert.Contains(methodsStr, "MethodEmbed3(arg0 string) string")

	})

	t.Run("ToMarshalJSON", func(t *testing.T) {
		assert := require.New(t)

		m, pkg := c.MethodsFromTypes(
			[]reflect.Type{zeroI},
			[]reflect.Type{zeroIE}).ToMarshalJSON("*page", "page")

		assert.Contains(m, "method6 := p.Method6()")
		assert.Contains(m, "Method0: method0,")
		assert.Contains(m, "return json.Marshal(&s)")

		assert.Contains(pkg, "github.com/gohugoio/hugo/common/herrors")
		assert.Contains(pkg, "encoding/json")

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
	Method7() interface{}
	Method8() herrors.ErrorContext
	method2()
	method9() os.FileInfo
}

type IEOnly interface {
	IEmbed
}
