// Copyright 2019 The Hugo Authors. All rights reserved.
// Some functions in this file (see comments) is based on the Go source code,
// copyright The Go Authors and  governed by a BSD-style license.
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

package hreflect

import (
	"fmt"
	"reflect"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/pkg/errors"
)

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

type Invoker struct {
	funcs func(name string) interface{}
}

func NewInvoker(funcs func(name string) interface{}) *Invoker {
	return &Invoker{funcs: funcs}
}

func (i *Invoker) InvokeFunction(path []string, args ...interface{}) (interface{}, error) {
	name := path[0]
	f := i.funcs(name)
	if f == nil {
		return err("function with name %s not found", name)
	}
	result, err := i.invoke(reflect.ValueOf(f), path, args)
	if err != nil {
		return nil, err
	}

	if !result.IsValid() {
		return nil, nil
	}

	return result.Interface(), nil
}

func (i *Invoker) InvokeMethod(receiver interface{}, path []string, args ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(receiver)
	result, err := i.invoke(v, path, args)
	if err != nil {
		return nil, err
	}

	if !result.IsValid() {
		return nil, nil
	}

	return result.Interface(), nil

}

func argsToValues(args []interface{}, typ reflect.Type) []reflect.Value {
	if len(args) == 0 {
		return nil
	}

	toArg := func(typ reflect.Type, v interface{}) reflect.Value {
		if typ == reflectValueType {
			return reflect.ValueOf(reflect.ValueOf(v))
		} else {
			return reflect.ValueOf(v)
		}
	}

	numFixed := len(args)
	if typ.IsVariadic() {
		numFixed = typ.NumIn() - 1
	}

	argsv := make([]reflect.Value, len(args))
	i := 0
	for ; i < numFixed && i < len(args); i++ {
		argsv[i] = toArg(typ.In(i), args[i])
	}
	if typ.IsVariadic() {
		argType := typ.In(typ.NumIn() - 1).Elem()
		for ; i < len(args); i++ {
			argsv[i] = toArg(argType, args[i])
		}
	}

	return argsv
}

func (i *Invoker) invoke(receiver reflect.Value, path []string, args []interface{}) (reflect.Value, error) {
	if len(path) == 0 {
		return receiver, nil
	}
	name := path[0]

	nextPath := 1
	typ := receiver.Type()
	receiver, isNil := indirect(receiver)
	if receiver.Kind() == reflect.Interface && isNil {
		return err("nil pointer evaluating %s.%s", typ, name)
	}

	ptr := receiver
	if ptr.Kind() != reflect.Interface && ptr.Kind() != reflect.Ptr && ptr.CanAddr() {
		ptr = ptr.Addr()
	}

	var fn reflect.Value
	if typ.Kind() == reflect.Func {
		fn = receiver
	} else {
		fn = ptr.MethodByName(name)
	}

	if fn.IsValid() {
		mt := fn.Type()
		if !isValidFunc(mt) {
			return err("method %s not valid", name)
		}

		var argsv []reflect.Value
		if len(path) == 1 {
			numArgs := len(args)
			if mt.IsVariadic() {
				if numArgs < (mt.NumIn() - 1) {
					return err("methods %s expects at leas %d arguments, got %d", name, mt.NumIn()-1, numArgs)
				}
			} else if numArgs != mt.NumIn() {
				return err("methods %s takes %d arguments, got %d", name, mt.NumIn(), numArgs)
			}
			argsv = argsToValues(args, mt)
		}

		result := fn.Call(argsv)
		if mt.NumOut() == 2 {
			if !result[1].IsZero() {
				return reflect.Value{}, result[1].Interface().(error)
			}
		}

		return i.invoke(result[0], path[nextPath:], args)
	}

	switch receiver.Kind() {
	case reflect.Struct:
		if f := receiver.FieldByName(name); f.IsValid() {
			return i.invoke(f, path[1:], args)
		} else {
			return err("no field with name %s found", name)
		}
	case reflect.Map:
		if p, ok := receiver.Interface().(maps.Params); ok {
			// Do case insensitive map lookup
			v := p.Get(path...)
			return reflect.ValueOf(v), nil
		}
		v := receiver.MapIndex(reflect.ValueOf(name))
		if !v.IsValid() {
			return reflect.Value{}, nil
		}
		return i.invoke(v, path[1:], args)
	}
	return receiver, nil
}

func err(s string, args ...interface{}) (reflect.Value, error) {
	return reflect.Value{}, errors.Errorf(s, args...)
}

func isValidFunc(typ reflect.Type) bool {
	switch {
	case typ.NumOut() == 1:
		return true
	case typ.NumOut() == 2 && typ.Out(1) == errorType:
		return true
	}
	return false
}
