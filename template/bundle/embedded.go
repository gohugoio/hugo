// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bundle

type Tmpl struct {
	Name string
	Data string
}

func (t *GoHtmlTemplate) EmbedShortcodes() {
	const k = "shortcodes"

	t.AddInternalTemplate(k, "highlight.html", `{{ $lang := index .Params 0 }}{{ highlight .Inner $lang }}`)
	t.AddInternalTemplate(k, "test.html", `This is a simple Test`)
	t.AddInternalTemplate(k, "figure.html", `<!-- image -->
<figure {{ if isset .Params "class" }}class="{{ index .Params "class" }}"{{ end }}>
    {{ if isset .Params "link"}}<a href="{{ index .Params "link"}}">{{ end }}
        <img src="{{ index .Params "src" }}" {{ if or (isset .Params "alt") (isset .Params "caption") }}alt="{{ if isset .Params "alt"}}{{ index .Params "alt"}}{{else}}{{ index .Params "caption" }}{{ end }}"{{ end }} />
    {{ if isset .Params "link"}}</a>{{ end }}
    {{ if or (or (isset .Params "title") (isset .Params "caption")) (isset .Params "attr")}}
    <figcaption>{{ if isset .Params "title" }}
        <h4>{{ index .Params "title" }}</h4>{{ end }}
        {{ if or (isset .Params "caption") (isset .Params "attr")}}<p>
        {{ index .Params "caption" }}
        {{ if isset .Params "attrlink"}}<a href="{{ index .Params "attrlink"}}"> {{ end }}
            {{ index .Params "attr" }}
        {{ if isset .Params "attrlink"}}</a> {{ end }}
        </p> {{ end }}
    </figcaption>
    {{ end }}
</figure>
<!-- image -->`)

}
