// Copyright 2015 The Hugo Authors. All rights reserved.
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

package parser

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gohugoio/hugo/parser/metadecoders"
)

func TestInterfaceToConfig(t *testing.T) {
	cases := []struct {
		input  interface{}
		format metadecoders.Format
		want   []byte
		isErr  bool
	}{
		// TOML
		{map[string]interface{}{}, metadecoders.TOML, nil, false},
		{
			map[string]interface{}{"title": "test 1"},
			metadecoders.TOML,
			[]byte("title = \"test 1\"\n"),
			false,
		},

		// YAML
		{map[string]interface{}{}, metadecoders.YAML, []byte("{}\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			metadecoders.YAML,
			[]byte("title: test 1\n"),
			false,
		},

		// JSON
		{map[string]interface{}{}, metadecoders.JSON, []byte("{}\n"), false},
		{
			map[string]interface{}{"title": "test 1"},
			metadecoders.JSON,
			[]byte("{\n   \"title\": \"test 1\"\n}\n"),
			false,
		},

		// Errors
		{nil, metadecoders.TOML, nil, true},
		{map[string]interface{}{}, "foo", nil, true},
	}

	for i, c := range cases {
		var buf bytes.Buffer

		err := InterfaceToConfig(c.input, c.format, &buf)
		if err != nil {
			if c.isErr {
				continue
			}
			t.Fatalf("[%d] unexpected error value: %v", i, err)
		}

		if !reflect.DeepEqual(buf.Bytes(), c.want) {
			t.Errorf("[%d] not equal:\nwant %q,\n got %q", i, c.want, buf.Bytes())
		}
	}
}
