// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"bytes"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestChainZeroTransformers(t *testing.T) {
	tr := New()
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	if err := tr.Apply(in, out); err != nil {
		t.Errorf("A zero transformer chain returned an error.")
	}
}

func TestChainingMultipleTransformers(t *testing.T) {
	f1 := func(ct FromTo) error {
		_, err := ct.To().Write(bytes.Replace(ct.From().Bytes(), []byte("f1"), []byte("f1r"), -1))
		return err
	}
	f2 := func(ct FromTo) error {
		_, err := ct.To().Write(bytes.Replace(ct.From().Bytes(), []byte("f2"), []byte("f2r"), -1))
		return err
	}
	f3 := func(ct FromTo) error {
		_, err := ct.To().Write(bytes.Replace(ct.From().Bytes(), []byte("f3"), []byte("f3r"), -1))
		return err
	}

	f4 := func(ct FromTo) error {
		_, err := ct.To().Write(bytes.Replace(ct.From().Bytes(), []byte("f4"), []byte("f4r"), -1))
		return err
	}

	tr := New(f1, f2, f3, f4)

	out := new(bytes.Buffer)
	if err := tr.Apply(out, strings.NewReader("Test: f4 f3 f1 f2 f1 The End.")); err != nil {
		t.Errorf("Multi transformer chain returned an error: %s", err)
	}

	expected := "Test: f4r f3r f1r f2r f1r The End."

	if out.String() != expected {
		t.Errorf("Expected %s got %s", expected, out.String())
	}
}

func TestNewEmptyTransforms(t *testing.T) {
	c := qt.New(t)
	transforms := NewEmpty()
	c.Assert(cap(transforms), qt.Equals, 20)
}
