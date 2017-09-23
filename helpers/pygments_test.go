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
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestParsePygmentsArgs(t *testing.T) {
	assert := require.New(t)

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
		spec, err := NewContentSpec(v)
		assert.NoError(err)

		result1, err := spec.createPygmentsOptionsString(this.in)
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
	assert := require.New(t)

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

		spec, err := NewContentSpec(v)
		assert.NoError(err)

		result, err := spec.createPygmentsOptionsString(this.in)
		if err != nil {
			t.Errorf("[%d] parsePygmentArgs failed: %s", i, err)
			continue
		}
		if result != expect {
			t.Errorf("[%d] parsePygmentArgs got %v but expected %v", i, result, expect)
		}
	}
}

func TestHlLinesToRanges(t *testing.T) {
	var zero [][2]int

	for _, this := range []struct {
		in       string
		expected interface{}
	}{
		{"", zero},
		{"1 4", [][2]int{[2]int{1, 1}, [2]int{4, 4}}},
		{"1-4 5-8", [][2]int{[2]int{1, 4}, [2]int{5, 8}}},
		{" 1   4 ", [][2]int{[2]int{1, 1}, [2]int{4, 4}}},
		{"1-4    5-8 ", [][2]int{[2]int{1, 4}, [2]int{5, 8}}},
		{"1-4 5", [][2]int{[2]int{1, 4}, [2]int{5, 5}}},
		{"4 5-9", [][2]int{[2]int{4, 4}, [2]int{5, 9}}},
		{" 1 -4 5 - 8  ", true},
		{"a b", true},
	} {
		got, err := hlLinesToRanges(this.in)

		if expectErr, ok := this.expected.(bool); ok && expectErr {
			if err == nil {
				t.Fatal("No error")
			}
		} else if err != nil {
			t.Fatalf("Got error: %s", err)
		} else if !reflect.DeepEqual(this.expected, got) {
			t.Fatalf("Expected\n%v but got\n%v", this.expected, got)
		}
	}
}
