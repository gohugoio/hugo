// +build gofuzz

// Copyright 2020 The Hugo Authors. All rights reserved.
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

package transform

import (
	"github.com/spf13/viper"
)

func FuzzMarkdownify(data []byte) int {
	v := viper.New()
	v.Set("contentDir", "content")
	ns := New(newDeps(v))

	for _, test := range []struct {
		s interface{}
	}{
		{string(data)},
	} {
		_, err := ns.Markdownify(test.s)
		if err != nil {
			return 0
		}
	}
	return 1
}
