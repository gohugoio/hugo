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

package helpers

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParsePygmentsArgs(t *testing.T) {
	for i, this := range []struct {
		in                 string
		pygmentsStyle      string
		pygmentsUseClasses bool
		expect1            interface{}
	}{
		{"", "foo", true, "encoding=utf8,noclasses=false,style=foo"},
		{"style=boo,noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"Style=boo, noClasses=true", "foo", true, "encoding=utf8,noclasses=true,style=boo"},
		{"noclasses=true", "foo", true, "encoding=utf8,noclasses=true,style=foo"},
		{"style=boo", "foo", true, "encoding=utf8,noclasses=false,style=boo"},
		{"boo=invalid", "foo", false, false},
		{"style", "foo", false, false},
	} {

		v := viper.New()
		v.Set("pygmentsStyle", this.pygmentsStyle)
		v.Set("pygmentsUseClasses", this.pygmentsUseClasses)

		result1, err := parsePygmentsOpts(v, this.in)
		if b, ok := this.expect1.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parsePygmentArgs didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
				continue
			}
			if result1 != this.expect1 {
				t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result1, this.expect1)
			}

		}
	}
}

func TestParseDefaultPygmentsArgs(t *testing.T) {
	expect := "encoding=utf8,noclasses=false,style=foo"

	for i, this := range []struct {
		in                 string
		pygmentsStyle      interface{}
		pygmentsUseClasses interface{}
		pygmentsOptions    string
	}{
		{"", "foo", true, "style=override,noclasses=override"},
		{"", nil, nil, "style=foo,noclasses=false"},
		{"style=foo,noclasses=false", nil, nil, "style=override,noclasses=override"},
		{"style=foo,noclasses=false", "override", false, "style=override,noclasses=override"},
	} {
		v := viper.New()

		v.Set("pygmentsOptions", this.pygmentsOptions)

		if s, ok := this.pygmentsStyle.(string); ok {
			v.Set("pygmentsStyle", s)
		}

		if b, ok := this.pygmentsUseClasses.(bool); ok {
			v.Set("pygmentsUseClasses", b)
		}

		result, err := parsePygmentsOpts(v, this.in)
		if err != nil {
			t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
			continue
		}
		if result != expect {
			t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result, expect)
		}
	}
}
