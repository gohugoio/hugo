// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"reflect"
	"sort"
)

// Slicer defines a very generic way to create a typed slice. This is used
// in collections.Slice template func to get types such as Pages, PageGroups etc.
// instead of the less useful []interface{}.
type Slicer interface {
	Slice(items any) (any, error)
}

// Slice returns a slice of all passed arguments.
func Slice(args ...any) any {
	if len(args) == 0 {
		return args
	}

	first := args[0]
	firstType := reflect.TypeOf(first)

	if firstType == nil {
		return args
	}

	if g, ok := first.(Slicer); ok {
		v, err := g.Slice(args)
		if err == nil {
			return v
		}

		// If Slice fails, the items are not of the same type and
		// []interface{} is the best we can do.
		return args
	}

	if len(args) > 1 {
		// This can be a mix of types.
		for i := 1; i < len(args); i++ {
			if firstType != reflect.TypeOf(args[i]) {
				// []interface{} is the best we can do
				return args
			}
		}
	}

	slice := reflect.MakeSlice(reflect.SliceOf(firstType), len(args), len(args))
	for i, arg := range args {
		slice.Index(i).Set(reflect.ValueOf(arg))
	}
	return slice.Interface()
}

// StringSliceToInterfaceSlice converts ss to []interface{}.
func StringSliceToInterfaceSlice(ss []string) []any {
	result := make([]any, len(ss))
	for i, s := range ss {
		result[i] = s
	}
	return result
}

type SortedStringSlice []string

// Contains returns true if s is in ss.
func (ss SortedStringSlice) Contains(s string) bool {
	i := sort.SearchStrings(ss, s)
	return i < len(ss) && ss[i] == s
}

// Count returns the number of times s is in ss.
func (ss SortedStringSlice) Count(s string) int {
	var count int
	i := sort.SearchStrings(ss, s)
	for i < len(ss) && ss[i] == s {
		count++
		i++
	}
	return count
}
