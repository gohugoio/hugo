// Copyright 2023 The Hugo Authors. All rights reserved.
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
package debug

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	Name    string
	Address any
	foo     string
}

func (u *User) M1() string         { return "" }
func (u *User) M2(v string) string { return "" }
func (u *User) m3(v string) string { return "" }

// Non Pointer type methods
func (u User) M4(v string) string { return "" }
func (u User) m5(v string) string { return "" }

func TestList(t *testing.T) {
	t.Parallel()

	namespace := new(Namespace)

	for i, test := range []struct {
		val    any
		expect []string
	}{
		// Map
		{map[string]any{"key1": 1, "key2": 2, "key3": 3}, []string{"key1", "key2", "key3"}},
		// Map non string keys
		{map[int]any{1: 1, 2: 2, 3: 3}, []string{"<int Value>", "<int Value>", "<int Value>"}},
		// Struct
		{User{}, []string{"Name", "Address", "M1", "M2", "M4"}},
		// Pointer
		{&User{}, []string{"Name", "Address", "M1", "M2", "M4"}},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			result := namespace.List(test.val)

			if !reflect.DeepEqual(result, test.expect) {
				t.Fatalf("List called with value: %#v got\n%#v but expected\n%#v", test.val, result, test.expect)
			}
		})
	}
}
