// Copyright 2019 The Hugo Authors. All rights reserved.
// The functions in this file is based on the Go source code, copyright
// The Go Authors and  governed by a BSD-style license.
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

// Package compare provides template functions for comparing values.
package compare

import (
	"fmt"
	"reflect"

	"github.com/gohugoio/hugo/common/hreflect"
)

// Boolean logic, based on:
// https://github.com/golang/go/blob/178a2c42254166cffed1b25fb1d3c7a5727cada6/src/text/template/funcs.go#L302

func truth(arg reflect.Value) bool {
	return hreflect.IsTruthfulValue(arg)
}

// getIf will return the given arg if it is considered truthful, else an empty string.
func (*Namespace) getIf(arg reflect.Value) reflect.Value {
	if truth(arg) {
		return arg
	}
	return reflect.ValueOf("")
}

func (*Namespace) invokeDot(args ...interface{}) interface{} {
	fmt.Println("invokeDot:", args)
	return "FOO"
}

// And computes the Boolean AND of its arguments, returning
// the first false argument it encounters, or the last argument.
func (*Namespace) And(arg0 reflect.Value, args ...reflect.Value) reflect.Value {
	if !truth(arg0) {
		return arg0
	}
	for i := range args {
		arg0 = args[i]
		if !truth(arg0) {
			break
		}
	}
	return arg0
}

// Or computes the Boolean OR of its arguments, returning
// the first true argument it encounters, or the last argument.
func (*Namespace) Or(arg0 reflect.Value, args ...reflect.Value) reflect.Value {
	if truth(arg0) {
		return arg0
	}
	for i := range args {
		arg0 = args[i]
		if truth(arg0) {
			break
		}
	}
	return arg0
}

// Not returns the Boolean negation of its argument.
func (*Namespace) Not(arg reflect.Value) bool {
	return !truth(arg)
}
