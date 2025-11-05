// Copyright 2025 The Hugo Authors. All rights reserved.
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

package page

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestPagePathsBuilder(t *testing.T) {
	c := qt.New(t)

	d := TargetPathDescriptor{}
	b := getPagePathBuilder(d)
	defer putPagePathBuilder(b)
	b.Add("foo", "bar")

	c.Assert(b.Path(0), qt.Equals, "/foo/bar")
}

func BenchmarkPagePathsBuilderPath(b *testing.B) {
	d := TargetPathDescriptor{}
	pb := getPagePathBuilder(d)
	defer putPagePathBuilder(pb)
	pb.Add("foo", "bar")

	for b.Loop() {
		_ = pb.Path(0)
	}
}

func BenchmarkPagePathsBuilderPathDir(b *testing.B) {
	d := TargetPathDescriptor{}
	pb := getPagePathBuilder(d)
	defer putPagePathBuilder(pb)
	pb.Add("foo", "bar")
	pb.prefixPath = "foo/"

	for b.Loop() {
		_ = pb.PathDir()
	}
}
