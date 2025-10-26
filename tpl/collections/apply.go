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

package collections

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/common/hreflect"
)

// Apply takes an array or slice c and returns a new slice with the function fname applied over it.
func (ns *Namespace) Apply(ctx context.Context, c any, fname string, args ...any) (any, error) {
	if c == nil {
		return make([]any, 0), nil
	}

	if fname == "apply" {
		return nil, errors.New("can't apply myself (no turtles allowed)")
	}

	seqv := reflect.ValueOf(c)
	seqv, isNil := hreflect.Indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	fnv, found := ns.lookupFunc(ctx, fname)
	if !found {
		return nil, errors.New("can't find function " + fname)
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		r := make([]any, seqv.Len())
		for i := range seqv.Len() {
			vv := seqv.Index(i)

			vvv, err := applyFnToThis(ctx, fnv, vv, args...)
			if err != nil {
				return nil, err
			}

			r[i] = vvv.Interface()
		}

		return r, nil
	default:
		return nil, fmt.Errorf("can't apply over %v", c)
	}
}

func applyFnToThis(ctx context.Context, fn, this reflect.Value, args ...any) (reflect.Value, error) {
	num := fn.Type().NumIn()
	if num > 0 && hreflect.IsContextType(fn.Type().In(0)) {
		args = append([]any{ctx}, args...)
	}

	n := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == "." {
			n[i] = this
		} else {
			n[i] = reflect.ValueOf(arg)
		}
	}

	if fn.Type().IsVariadic() {
		num--
	}

	// TODO(bep) see #1098 - also see template_tests.go
	/*if len(args) < num {
		return reflect.ValueOf(nil), errors.New("Too few arguments")
	} else if len(args) > num {
		return reflect.ValueOf(nil), errors.New("Too many arguments")
	}*/

	for i := range num {
		// AssignableTo reports whether xt is assignable to type targ.
		if xt, targ := n[i].Type(), fn.Type().In(i); !xt.AssignableTo(targ) {
			return reflect.ValueOf(nil), errors.New("called apply using " + xt.String() + " as type " + targ.String())
		}
	}

	res := fn.Call(n)

	if len(res) == 1 || res[1].IsNil() {
		return res[0], nil
	}
	return reflect.ValueOf(nil), res[1].Interface().(error)
}

func (ns *Namespace) lookupFunc(ctx context.Context, fname string) (reflect.Value, bool) {
	namespace, methodName, ok := strings.Cut(fname, ".")
	if !ok {
		return ns.deps.GetTemplateStore().GetFunc(fname)
	}

	// Namespace
	nv, found := ns.lookupFunc(ctx, namespace)
	if !found {
		return reflect.Value{}, false
	}

	fn, ok := nv.Interface().(func(context.Context, ...any) (any, error))
	if !ok {
		return reflect.Value{}, false
	}
	v, err := fn(ctx)
	if err != nil {
		panic(err)
	}
	nv = reflect.ValueOf(v)

	// method
	m := hreflect.GetMethodByName(nv, methodName)

	if m.Kind() == reflect.Invalid {
		return reflect.Value{}, false
	}
	return m, true
}
