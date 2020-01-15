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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/tpl"
)

// Apply takes a map, array, or slice and returns a new slice with the function fname applied over it.
func (ns *Namespace) Apply(seq interface{}, fname string, args ...interface{}) (interface{}, error) {
	if seq == nil {
		return make([]interface{}, 0), nil
	}

	if fname == "apply" {
		return nil, errors.New("can't apply myself (no turtles allowed)")
	}

	seqv := reflect.ValueOf(seq)
	seqv, isNil := indirect(seqv)
	if isNil {
		return nil, errors.New("can't iterate over a nil value")
	}

	fnv, found := ns.lookupFunc(fname)
	if !found {
		return nil, errors.New("can't find function " + fname)
	}

	// fnv := reflect.ValueOf(fn)

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		r := make([]interface{}, seqv.Len())
		for i := 0; i < seqv.Len(); i++ {
			vv := seqv.Index(i)

			vvv, err := applyFnToThis(fnv, vv, args...)

			if err != nil {
				return nil, err
			}

			r[i] = vvv.Interface()
		}

		return r, nil
	default:
		return nil, fmt.Errorf("can't apply over %v", seq)
	}
}

func applyFnToThis(fn, this reflect.Value, args ...interface{}) (reflect.Value, error) {
	n := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == "." {
			n[i] = this
		} else {
			n[i] = reflect.ValueOf(arg)
		}
	}

	num := fn.Type().NumIn()

	if fn.Type().IsVariadic() {
		num--
	}

	// TODO(bep) see #1098 - also see template_tests.go
	/*if len(args) < num {
		return reflect.ValueOf(nil), errors.New("Too few arguments")
	} else if len(args) > num {
		return reflect.ValueOf(nil), errors.New("Too many arguments")
	}*/

	for i := 0; i < num; i++ {
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

func (ns *Namespace) lookupFunc(fname string) (reflect.Value, bool) {
	if !strings.ContainsRune(fname, '.') {
		templ := ns.deps.Tmpl().(tpl.TemplateFuncGetter)
		return templ.GetFunc(fname)
	}

	ss := strings.SplitN(fname, ".", 2)

	// namespace
	nv, found := ns.lookupFunc(ss[0])
	if !found {
		return reflect.Value{}, false
	}

	// method
	m := nv.MethodByName(ss[1])
	// if reflect.DeepEqual(m, reflect.Value{}) {
	if m.Kind() == reflect.Invalid {
		return reflect.Value{}, false
	}
	return m, true
}

// indirect is borrowed from the Go stdlib: 'text/template/exec.go'
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}

func indirectInterface(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
		if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
			break
		}
	}
	return v, false
}
