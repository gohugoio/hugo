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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigBasename(t *testing.T) {
	for _, this := range []struct {
		kind   string
		expect string
	}{
		{"toml", "config.toml"},
		{"json", "config.json"},
		{"yaml", "config.yml"},
		{"other", "config.other"},
	} {
		result := ConfigBasename(this.kind)
		assert.Equal(t, this.expect, result)
	}
}

func TestFormatToLeadRune(t *testing.T) {
	for i, this := range []struct {
		kind   string
		expect rune
	}{
		{"yaml", '-'},
		{"yml", '-'},
		{"toml", '+'},
		{"json", '{'},
		{"js", '{'},
		{"unknown", '+'},
	} {
		result := FormatToLeadRune(this.kind)

		if result != this.expect {
			t.Errorf("[%d] got %q but expected %q", i, result, this.expect)
		}
	}
}
