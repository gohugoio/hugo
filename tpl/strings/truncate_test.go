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
		v1    any
		v2    any
		v3    any
		want  any
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
		{
			42,
			template.HTML(`With strangely formatted
							<a
							href="#"
								target="_blank"
							>HTML</a
							>
							inside.`),
			nil,
			template.HTML(`With strangely formatted
							<a
							href="#"
								target="_blank"
							>HTML …</a>`),
			false,
		},
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

func BenchmarkTruncate(b *testing.B) {
	b.Run("Plain text", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(10, "I am a test sentence")
		}
	})

	b.Run("With link", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(10, template.HTML("I have a <a href='/markdown'>Markdown link</a> inside"))
		}
	})

	b.Run("Plain text (medium)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(371, `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`)
		}
	})

	b.Run("With HTML (medium)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(371, template.HTML(`Lorem ipsum dolor sit amet, <span>consectetur adipiscing elit</span>, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, <p>quis nostrud exercitation ullamco</p> laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur <code>sint occaecat cupidatat</code> non proident, <a href="#my-text">sunt in culpa qui officia deserunt mollit anim id est</a> laborum.`))
		}
	})

	b.Run("Plain text (long)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(1850, `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.
Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.
Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.
Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur?
Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?
At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint
occaecati cupiditate non provident, similique sunt in culpa qui officia deserunt mollitia animi, id est laborum et dolorum fuga.
Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat
facere possimus, omnis voluptas assumenda est, omnis dolor repellendus.
Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae.
Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.`)
		}
	})

	b.Run("HTML text (long)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ns.Truncate(1850, template.HTML(`Lorem ipsum dolor sit amet, <span>consectetur adipiscing elit</span>, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, <p>quis nostrud exercitation ullamco</p> laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur <code>sint occaecat cupidatat</code> non proident, <a href="#my-text">sunt in culpa qui officia deserunt mollit anim id est</a> laborum.
Sed ut <span>perspiciatis unde omnis iste natus error sit voluptatem</span> accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.
Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.
Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem.
Ut enim ad minima veniam, quis nostrum exercitationem <a href="/home-page">ullam corporis suscipit laboriosam</a>, nisi ut aliquid ex ea commodi consequatur?
Quis autem vel eum iure reprehenderit qui in ea voluptate velit esse quam nihil molestiae consequatur, vel illum qui dolorem eum fugiat quo voluptas nulla pariatur?
At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint
occaecati cupiditate non provident, <div class="my-text-class">similique sunt in culpa qui officia deserunt mollitia animi</div>, id est laborum et dolorum fuga.
Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat
facere possimus, omnis voluptas assumenda est, omnis dolor repellendus.
Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe <h3>eveniet ut et voluptates repudiandae sint et molestiae non recusandae.</h3>
Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.`))
		}
	})

}
