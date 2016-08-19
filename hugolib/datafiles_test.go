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

package hugolib

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/hugo/parser"
	"github.com/spf13/hugo/source"
)

func TestDataDirJSON(t *testing.T) {
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("test/foo.json"), Content: []byte(`{ "bar": "foofoo"  }`)},
		{Name: filepath.FromSlash("test.json"), Content: []byte(`{ "hello": [ { "world": "foo" } ] }`)},
	}

	expected, err := parser.HandleJSONMetaData([]byte(`{ "test": { "hello": [{ "world": "foo"  }] , "foo": { "bar":"foofoo" } } }`))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

func TestDataDirToml(t *testing.T) {
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("test/kung.toml"), Content: []byte("[foo]\nbar = 1")},
	}

	expected, err := parser.HandleTOMLMetaData([]byte("[test]\n[test.kung]\n[test.kung.foo]\nbar = 1"))

	if err != nil {
		t.Fatalf("Error %s", err)
	}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

func TestDataDirYAMLWithOverridenValue(t *testing.T) {
	sources := []source.ByteSource{
		// filepath.Walk walks the files in lexical order, '/' comes before '.'. Simulate this:
		{Name: filepath.FromSlash("a.yaml"), Content: []byte("a: 1")},
		{Name: filepath.FromSlash("test/v1.yaml"), Content: []byte("v1-2: 2")},
		{Name: filepath.FromSlash("test/v2.yaml"), Content: []byte("v2:\n- 2\n- 3")},
		{Name: filepath.FromSlash("test.yaml"), Content: []byte("v1: 1")},
	}

	expected := map[string]interface{}{"a": map[string]interface{}{"a": 1},
		"test": map[string]interface{}{"v1": map[string]interface{}{"v1-2": 2}, "v2": map[string]interface{}{"v2": []interface{}{2, 3}}}}

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: sources}})
}

// issue 892
func TestDataDirMultipleSources(t *testing.T) {
	s1 := []source.ByteSource{
		{Name: filepath.FromSlash("test/first.toml"), Content: []byte("bar = 1")},
	}

	s2 := []source.ByteSource{
		{Name: filepath.FromSlash("test/first.toml"), Content: []byte("bar = 2")},
		{Name: filepath.FromSlash("test/second.toml"), Content: []byte("tender = 2")},
	}

	expected, _ := parser.HandleTOMLMetaData([]byte("[test.first]\nbar = 1\n[test.second]\ntender=2"))

	doTestDataDir(t, expected, []source.Input{&source.InMemorySource{ByteSource: s1}, &source.InMemorySource{ByteSource: s2}})

}

func TestDataDirUnknownFormat(t *testing.T) {
	sources := []source.ByteSource{
		{Name: filepath.FromSlash("test.roml"), Content: []byte("boo")},
	}
	s := &Site{}
	err := s.loadData([]source.Input{&source.InMemorySource{ByteSource: sources}})
	if err != nil {
		t.Fatalf("Should not return an error")
	}
}

func doTestDataDir(t *testing.T, expected interface{}, sources []source.Input) {
	s := &Site{}
	err := s.loadData(sources)
	if err != nil {
		t.Fatalf("Error loading data: %s", err)
	}
	if !reflect.DeepEqual(expected, s.Data) {
		t.Errorf("Expected structure\n%#v got\n%#v", expected, s.Data)
	}
}
