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

package metainject

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/transform"
)

func TestHugoGeneratorInject(t *testing.T) {
	hugoGeneratorTag = "META"
	for i, this := range []struct {
		in     string
		expect string
	}{
		{`<head>
	<foo />
</head>`, `<head>
	META
	<foo />
</head>`},
		{`<HEAD>
	<foo />
</HEAD>`, `<HEAD>
	META
	<foo />
</HEAD>`},
		{`<head><meta name="generator" content="Jekyll"></head>`, `<head><meta name="generator" content="Jekyll"></head>`},
		{`<head><meta name='generator' content='Jekyll'></head>`, `<head><meta name='generator' content='Jekyll'></head>`},
		{`<head><meta name=generator content=Jekyll></head>`, `<head><meta name=generator content=Jekyll></head>`},
		{`<head><META     NAME="GENERATOR" content="Jekyll"></head>`, `<head><META     NAME="GENERATOR" content="Jekyll"></head>`},
		{"", ""},
		{"</head>", "</head>"},
		{"<head>", "<head>\n\tMETA"},
	} {
		in := strings.NewReader(this.in)
		out := new(bytes.Buffer)

		tr := transform.New(HugoGenerator)
		tr.Apply(out, in)

		if out.String() != this.expect {
			t.Errorf("[%d] Expected \n%q got \n%q", i, this.expect, out.String())
		}
	}
}
