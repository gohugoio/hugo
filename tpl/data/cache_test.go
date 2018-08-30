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

package data

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Parallel()

	fs := new(afero.MemMapFs)

	for i, test := range []struct {
		path    string
		content []byte
		ignore  bool
	}{
		{"http://Foo.Bar/foo_Bar-Foo", []byte(`T€st Content 123`), false},
		{"fOO,bar:foo%bAR", []byte(`T€st Content 123 fOO,bar:foo%bAR`), false},
		{"FOo/BaR.html", []byte(`FOo/BaR.html T€st Content 123`), false},
		{"трям/трям", []byte(`T€st трям/трям Content 123`), false},
		{"은행", []byte(`T€st C은행ontent 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 123`), false},
		{"Банковский кассир", []byte(`Банковский кассир T€st Content 456`), true},
	} {
		msg := fmt.Sprintf("Test #%d: %v", i, test)

		cfg := viper.New()

		c, err := getCache(test.path, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)
		assert.Nil(t, c, msg)

		err = writeCache(test.path, test.content, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)

		c, err = getCache(test.path, fs, cfg, test.ignore)
		assert.NoError(t, err, msg)

		if test.ignore {
			assert.Nil(t, c, msg)
		} else {
			assert.Equal(t, string(test.content), string(c))
		}
	}
}
