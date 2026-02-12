// Copyright 2025 The Hugo Authors. All rights reserved.
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
	"iter"
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/common/hmaps"
	"github.com/gohugoio/hugo/common/hreflect"
)

// WhereIter returns an iterator that yields matching elements from collection c.
// For slices/arrays it yields (index, element), for maps it yields (key, value).
func (ns *Namespace) WhereIter(ctx context.Context, c, key any, args ...any) (iter.Seq2[any, any], error) {
	seqv, isNil := hreflect.Indirect(reflect.ValueOf(c))
	if isNil {
		return nil, errors.New("can't iterate over a nil value of type " + reflect.ValueOf(c).Type().String())
	}

	mv, op, err := parseWhereArgs(args...)
	if err != nil {
		return nil, err
	}

	ctxv := reflect.ValueOf(ctx)

	var path []string
	kv := reflect.ValueOf(key)
	if kv.Kind() == reflect.String {
		path = strings.Split(strings.Trim(kv.String(), "."), ".")
	}

	switch seqv.Kind() {
	case reflect.Array, reflect.Slice:
		return ns.whereArrayIter(ctxv, seqv, kv, mv, path, op), nil
	case reflect.Map:
		return ns.whereMapIter(ctxv, seqv, kv, mv, path, op), nil
	default:
		return nil, fmt.Errorf("can't iterate over %T", c)
	}
}

func (ns *Namespace) collectWhereArray(seqv reflect.Value, it iter.Seq2[any, any]) any {
	rv := reflect.MakeSlice(seqv.Type(), 0, 0)
	for _, v := range it {
		rv = reflect.Append(rv, reflect.ValueOf(v))
	}
	return rv.Interface()
}

func (ns *Namespace) collectWhereMap(seqv reflect.Value, it iter.Seq2[any, any]) any {
	rv := reflect.MakeMap(seqv.Type())
	for k, v := range it {
		rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
	}
	return rv.Interface()
}

func (ns *Namespace) whereArrayIter(ctxv, seqv, kv, mv reflect.Value, path []string, op string) iter.Seq2[any, any] {
	return ns.whereArrayIterWithErr(ctxv, seqv, kv, mv, path, op, nil)
}

// elemResolver resolves a sub-element from a reflect.Value.
// Built once before the loop to avoid repeated type checks and reflect.ValueOf allocations.
type elemResolver func(reflect.Value) (reflect.Value, error)

// newElemResolver returns a resolver optimized for the given element type and path.
// Returns nil if the element type is an interface or optimization isn't possible,
// in which case the caller should fall back to evaluateSubElem.
func (ns *Namespace) newElemResolver(ctxv reflect.Value, elemType reflect.Type, path []string) elemResolver {
	if elemType.Kind() == reflect.Interface {
		return nil
	}

	// Unwrap pointer.
	baseType := elemType
	isPtr := baseType.Kind() == reflect.Pointer
	if isPtr {
		baseType = baseType.Elem()
	}

	if baseType == reflect.TypeFor[hmaps.Params]() {
		return func(v reflect.Value) (reflect.Value, error) {
			params := v.Interface().(hmaps.Params)
			return reflect.ValueOf(params.GetNested(path...)), nil
		}
	}

	if len(path) != 1 {
		return nil
	}
	name := path[0]

	// Check for method first, matching evaluateSubElem order.
	ptrType := baseType
	if !hreflect.IsInterfaceOrPointer(ptrType.Kind()) {
		ptrType = reflect.PointerTo(baseType)
	}
	mt := hreflect.GetMethodByNameForType(ptrType, name)
	if mt.Func.IsValid() {
		return ns.newMethodResolver(ctxv, elemType, mt)
	}

	switch baseType.Kind() {
	case reflect.Map:
		if baseType.Key().Kind() == reflect.String {
			kv := reflect.ValueOf(name)
			return func(v reflect.Value) (reflect.Value, error) {
				if isPtr {
					v = v.Elem()
				}
				return v.MapIndex(kv), nil
			}
		}
	case reflect.Struct:
		ft, ok := baseType.FieldByName(name)
		if ok {
			if ft.PkgPath != "" && !ft.Anonymous {
				return func(v reflect.Value) (reflect.Value, error) {
					return zero, fmt.Errorf("%s is an unexported field of struct type %s", name, elemType)
				}
			}
			idx := ft.Index
			return func(v reflect.Value) (reflect.Value, error) {
				if isPtr {
					v = v.Elem()
				}
				return v.FieldByIndex(idx), nil
			}
		}
	}

	return nil
}

func (ns *Namespace) newMethodResolver(ctxv reflect.Value, elemType reflect.Type, mt reflect.Method) elemResolver {
	if mt.PkgPath != "" {
		return func(v reflect.Value) (reflect.Value, error) {
			return zero, fmt.Errorf("%s is an unexported method of type %s", mt.Name, elemType)
		}
	}

	numIn := mt.Type.NumIn()
	maxNumIn := 1
	needsCtx := numIn > 1 && hreflect.IsContextType(mt.Type.In(1))
	if needsCtx {
		maxNumIn = 2
	}

	switch {
	case mt.Type.NumIn() > maxNumIn:
		return nil
	case mt.Type.NumOut() == 0:
		return nil
	case mt.Type.NumOut() > 2:
		return nil
	case mt.Type.NumOut() == 1 && mt.Type.Out(0).Implements(errorType):
		return nil
	case mt.Type.NumOut() == 2 && !mt.Type.Out(1).Implements(errorType):
		return nil
	}

	fn := mt.Func
	hasErrOut := mt.Type.NumOut() == 2

	isPtr := elemType.Kind() == reflect.Pointer

	return func(v reflect.Value) (reflect.Value, error) {
		recv := v
		if !isPtr && !hreflect.IsInterfaceOrPointer(recv.Kind()) && recv.CanAddr() {
			recv = recv.Addr()
		}
		var args []reflect.Value
		if needsCtx {
			args = []reflect.Value{recv, ctxv}
		} else {
			args = []reflect.Value{recv}
		}
		res := fn.Call(args)
		if hasErrOut && !res[1].IsNil() {
			return zero, nil
		}
		return res[0], nil
	}
}

// whereArrayIterWithErr returns an iterator over matching slice/array elements.
// If errp is non-nil, any error from checkCondition is written to it and iteration stops.
func (ns *Namespace) whereArrayIterWithErr(ctxv, seqv, kv, mv reflect.Value, path []string, op string, errp *error) iter.Seq2[any, any] {
	var resolve elemResolver
	if kv.Kind() == reflect.String && len(path) > 0 {
		resolve = ns.newElemResolver(ctxv, seqv.Type().Elem(), path)
	}

	return func(yield func(any, any) bool) {
		for i := range seqv.Len() {
			var vvv reflect.Value
			rvv := seqv.Index(i)

			if resolve != nil {
				var err error
				vvv, err = resolve(rvv)
				if err != nil {
					if errp != nil {
						*errp = err
					}
					return
				}
			} else if kv.Kind() == reflect.String {
				if params, ok := rvv.Interface().(hmaps.Params); ok {
					vvv = reflect.ValueOf(params.GetNested(path...))
				} else {
					vvv = rvv
					for j, elemName := range path {
						var err error
						vvv, err = evaluateSubElem(ctxv, vvv, elemName)
						if err != nil {
							continue
						}
						if j < len(path)-1 && vvv.IsValid() {
							if params, ok := vvv.Interface().(hmaps.Params); ok {
								vvv = reflect.ValueOf(params.GetNested(path[j+1:]...))
								break
							}
						}
					}
				}
			} else {
				vv, _ := hreflect.Indirect(rvv)
				if vv.Kind() == reflect.Map && kv.Type().AssignableTo(vv.Type().Key()) {
					vvv = vv.MapIndex(kv)
				}
			}

			ok, err := ns.checkCondition(vvv, mv, op)
			if err != nil {
				if errp != nil {
					*errp = err
				}
				return
			}
			if ok {
				if !yield(i, rvv.Interface()) {
					return
				}
			}
		}
	}
}

func (ns *Namespace) whereMapIter(ctxv, seqv, kv, mv reflect.Value, path []string, op string) iter.Seq2[any, any] {
	return ns.whereMapIterWithErr(ctxv, seqv, kv, mv, path, op, nil)
}

func (ns *Namespace) whereMapIterWithErr(ctxv, seqv, kv, mv reflect.Value, path []string, op string, errp *error) iter.Seq2[any, any] {
	return func(yield func(any, any) bool) {
		k := reflect.New(seqv.Type().Key()).Elem()
		elemv := reflect.New(seqv.Type().Elem()).Elem()
		miter := seqv.MapRange()
		for miter.Next() {
			k.SetIterKey(miter)
			elemv.SetIterValue(miter)

			matched := false
			switch elemv.Kind() {
			case reflect.Array, reflect.Slice:
				r, err := ns.checkWhereArray(ctxv, elemv, kv, mv, path, op)
				if err != nil {
					if errp != nil {
						*errp = err
					}
					return
				}
				if rr := reflect.ValueOf(r); rr.Kind() == reflect.Slice && rr.Len() > 0 {
					matched = true
				}
			case reflect.Interface:
				elemvv, isNil := hreflect.Indirect(elemv)
				if !isNil {
					switch elemvv.Kind() {
					case reflect.Array, reflect.Slice:
						r, err := ns.checkWhereArray(ctxv, elemvv, kv, mv, path, op)
						if err != nil {
							if errp != nil {
								*errp = err
							}
							return
						}
						if rr := reflect.ValueOf(r); rr.Kind() == reflect.Slice && rr.Len() > 0 {
							matched = true
						}
					}
				}
			}

			if matched {
				if !yield(k.Interface(), elemv.Interface()) {
					return
				}
			}
		}
	}
}
