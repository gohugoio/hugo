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

// whereArrayIterWithErr returns an iterator over matching slice/array elements.
// If errp is non-nil, any error from checkCondition is written to it and iteration stops.
func (ns *Namespace) whereArrayIterWithErr(ctxv, seqv, kv, mv reflect.Value, path []string, op string, errp *error) iter.Seq2[any, any] {
	return func(yield func(any, any) bool) {
		for i := range seqv.Len() {
			var vvv reflect.Value
			rvv := seqv.Index(i)

			if kv.Kind() == reflect.String {
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
