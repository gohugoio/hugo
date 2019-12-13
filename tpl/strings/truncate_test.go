// Copyright 2016 The Hugo Authors. All rights reserved.
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

package strings

import (
	"html/template"

	"reflect"
	"strings"
	"testing"
)

func TestTruncate(t *testing.T) {
	t.Parallel()

	var err error
	cases := []struct {
		v1    interface{}
		v2    interface{}
		v3    interface{}
		want  interface{}
		isErr bool
	}{
		{10, "I am a test sentence", nil, template.HTML("I am a …"), false},
		{10, "", "I am a test sentence", template.HTML("I am a"), false},
		{10, "", "a b c d e f g h i j k", template.HTML("a b c d e"), false},
		{12, "", "<b>Should be escaped</b>", template.HTML("&lt;b&gt;Should be"), false},
		{10, template.HTML(" <a href='#'>Read more</a>"), "I am a test sentence", template.HTML("I am a <a href='#'>Read more</a>"), false},
		{20, template.HTML("I have a <a href='/markdown'>Markdown link</a> inside."), nil, template.HTML("I have a <a href='/markdown'>Markdown …</a>"), false},
		{10, "IamanextremelylongwordthatjustgoesonandonandonjusttoannoyyoualmostasifIwaswritteninGermanActuallyIbettheresagermanwordforthis", nil, template.HTML("Iamanextre …"), false},
		{10, template.HTML("<p>IamanextremelylongwordthatjustgoesonandonandonjusttoannoyyoualmostasifIwaswritteninGermanActuallyIbettheresagermanwordforthis</p>"), nil, template.HTML("<p>Iamanextre …</p>"), false},
		{13, template.HTML("With <a href=\"/markdown\">Markdown</a> inside."), nil, template.HTML("With <a href=\"/markdown\">Markdown …</a>"), false},
		{14, "Hello中国 Good 好的", nil, template.HTML("Hello中国 Good 好 …"), false},
		{15, "", template.HTML("A <br> tag that's not closed"), template.HTML("A <br> tag that's"), false},
		{14, template.HTML("<p>Hello中国 Good 好的</p>"), nil, template.HTML("<p>Hello中国 Good 好 …</p>"), false},
		{2, template.HTML("<p>P1</p><p>P2</p>"), nil, template.HTML("<p>P1 …</p>"), false},
		{3, template.HTML(strings.Repeat("<p>P</p>", 20)), nil, template.HTML("<p>P</p><p>P</p><p>P …</p>"), false},
		{18, template.HTML("<p>test <b>hello</b> test something</p>"), nil, template.HTML("<p>test <b>hello</b> test …</p>"), false},
		{4, template.HTML("<p>a<b><i>b</b>c d e</p>"), nil, template.HTML("<p>a<b><i>b</b>c …</p>"), false},
		{10, nil, nil, template.HTML(""), true},
		{nil, nil, nil, template.HTML(""), true},
	}
	for i, c := range cases {
		var result template.HTML
		if c.v2 == nil {
			result, err = ns.Truncate(c.v1)
		} else if c.v3 == nil {
			result, err = ns.Truncate(c.v1, c.v2)
		} else {
			result, err = ns.Truncate(c.v1, c.v2, c.v3)
		}

		if c.isErr {
			if err == nil {
				t.Errorf("[%d] Slice didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] failed: %s", i, err)
				continue
			}
			if !reflect.DeepEqual(result, c.want) {
				t.Errorf("[%d] got '%s' but expected '%s'", i, result, c.want)
			}
		}
	}

	// Too many arguments
	_, err = ns.Truncate(10, " ...", "I am a test sentence", "wrong")
	if err == nil {
		t.Errorf("Should have errored")
	}

}
