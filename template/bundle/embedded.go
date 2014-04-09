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
	t.AddInternalShortcode("highlight.html", `{{ .Get 0 | highlight .Inner  }}`)
	t.AddInternalShortcode("test.html", `This is a simple Test`)
	t.AddInternalShortcode("figure.html", `<!-- image -->
<figure {{ with .Get "class" }}class="{{.}}"{{ end }}>
    {{ with .Get "link"}}<a href="{{.}}">{{ end }}
        <img src="{{ .Get "src" }}" {{ if or (.Get "alt") (.Get "caption") }}alt="{{ with .Get "alt"}}{{.}}{{else}}{{ .Get "caption" }}{{ end }}"{{ end }} />
    {{ if .Get "link"}}</a>{{ end }}
    {{ if or (or (.Get "title") (.Get "caption")) (.Get "attr")}}
    <figcaption>{{ if isset .Params "title" }}
        <h4>{{ .Get "title" }}</h4>{{ end }}
        {{ if or (.Get "caption") (.Get "attr")}}<p>
        {{ .Get "caption" }}
        {{ with .Get "attrlink"}}<a href="{{.}}"> {{ end }}
            {{ .Get "attr" }}
        {{ if .Get "attrlink"}}</a> {{ end }}
        </p> {{ end }}
    </figcaption>
    {{ end }}
</figure>
<!-- image -->`)
}

func (t *GoHtmlTemplate) EmbedTemplates() {

	t.AddInternalTemplate("_default", "rss.xml", `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
      <title>{{ .Title }} on {{ .Site.Title }} </title>
      <generator uri="https://hugo.spf13.com">Hugo</generator>
    <link>{{ .Permalink }}</link>
    {{ with .Site.LanguageCode }}<language>{{.}}</language>{{end}}
    {{ with .Site.Author }}<author>{{.}}</author>{{end}}
    {{ with .Site.Copyright }}<copyright>{{.}}</copyright>{{end}}
    <updated>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</updated>
    {{ range first 15 .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</pubDate>
      {{with .Site.Author}}<author>{{.}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Content | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>`)

}
