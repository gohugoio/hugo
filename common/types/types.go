// Copyright 2017-present The Hugo Authors. All rights reserved.
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

// Package types contains types shared between packages in Hugo.
package types

import (
	"fmt"

	"github.com/spf13/cast"
)

// KeyValues holds an key and a slice of values.
type KeyValues struct {
	Key    interface{}
	Values []interface{}
}

// KeyString returns the key as a string, an empty string if conversion fails.
func (k KeyValues) KeyString() string {
	return cast.ToString(k.Key)
}

func (k KeyValues) String() string {
	return fmt.Sprintf("%v: %v", k.Key, k.Values)
}

func NewKeyValuesStrings(key string, values ...string) KeyValues {
	iv := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		iv[i] = values[i]
	}
	return KeyValues{Key: key, Values: iv}
}
