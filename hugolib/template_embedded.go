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

package hugolib

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
    {{ with .Site.Author.name }}<author>{{.}}</author>{{end}}
    {{ with .Site.Copyright }}<copyright>{{.}}</copyright>{{end}}
    <updated>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</updated>
    {{ range first 15 .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 MST" }}</pubDate>
      {{with .Site.Author.name}}<author>{{.}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Content | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>`)

	t.AddInternalTemplate("_default", "sitemap.xml", `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>
    <lastmod>{{ safeHtml ( .Date.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}
  </url>
  {{ end }}
</urlset>`)

	t.AddInternalTemplate("", "disqus.html", `{{ if .Site.DisqusShortname }}<div id="disqus_thread"></div>
<script type="text/javascript">
    var disqus_shortname = '{{ .Site.DisqusShortname }}';
    var disqus_identifier = '{{with .GetParam "disqus_identifier" }}{{ . }}{{ else }}{{ .Permalink }}{{end}}';
    var disqus_title = '{{with .GetParam "disqus_title" }}{{ . }}{{ else }}{{ .Title }}{{end}}';
    var disqus_url = '{{with .GetParam "disqus_url" }}{{ . | html  }}{{ else }}{{ .Permalink }}{{end}}';

    (function() {
        var dsq = document.createElement('script'); dsq.type = 'text/javascript'; dsq.async = true;
        dsq.src = '//' + disqus_shortname + '.disqus.com/embed.js';
        (document.getElementsByTagName('head')[0] || document.getElementsByTagName('body')[0]).appendChild(dsq);
    })();
</script>
<noscript>Please enable JavaScript to view the <a href="http://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
<a href="http://disqus.com" class="dsq-brlink">comments powered by <span class="logo-disqus">Disqus</span></a>{{end}}`)

}
