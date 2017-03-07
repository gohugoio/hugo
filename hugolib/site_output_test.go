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

package hugolib

import (
	"reflect"
	"testing"

	"github.com/spf13/hugo/output"
)

func TestDefaultOutputDefinitions(t *testing.T) {
	defs := defaultOutputDefinitions

	tests := []struct {
		name string
		kind string
		want []output.Type
	}{
		{"RSS not for regular pages", KindPage, []output.Type{output.HTMLType}},
		{"Home Sweet Home", KindHome, []output.Type{output.HTMLType, output.RSSType}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defs.ForKind(tt.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("siteOutputDefinitions.ForKind(%v) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}
