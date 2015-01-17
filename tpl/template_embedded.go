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

package tpl

type Tmpl struct {
	Name string
	Data string
}

func (t *GoHtmlTemplate) EmbedShortcodes() {
	t.AddInternalShortcode("ref.html", `{{ .Get 0 | ref .Page }}`)
	t.AddInternalShortcode("relref.html", `{{ .Get 0 | relref .Page }}`)
	t.AddInternalShortcode("highlight.html", `{{ .Get 0 | highlight .Inner  }}`)
	t.AddInternalShortcode("test.html", `This is a simple Test`)
	t.AddInternalShortcode("figure.html", `<!-- image -->
<figure {{ with .Get "class" }}class="{{.}}"{{ end }}>
    {{ with .Get "link"}}<a href="{{.}}">{{ end }}
        <img src="{{ .Get "src" }}" {{ if or (.Get "alt") (.Get "caption") }}alt="{{ with .Get "alt"}}{{.}}{{else}}{{ .Get "caption" }}{{ end }}" {{ end }}{{ with .Get "width" }}width="{{.}}" {{ end }}/>
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
    <title>{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}</title>
    <link>{{ .Permalink }}</link>
    <description>Recent content {{ with .Title }}in {{.}} {{ end }}on {{ .Site.Title }}</description>
    <generator>Hugo -- gohugo.io</generator>
    {{ with .Site.LanguageCode }}<language>{{.}}</language>{{end}}
    {{ with .Site.Author.email }}<managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}
    {{ with .Site.Author.email }}<webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}
    {{ with .Site.Copyright }}<copyright>{{.}}</copyright>{{end}}
    <lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05" }} {{ with .Date.Format "MST" }}{{ if eq . "UTC" }}UT{{else}}{{.}}{{end}}{{end}}</lastBuildDate>
    <atom:link href="{{.Url}}" rel="self" type="application/rss+xml" />
    {{ range first 15 .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05" }} {{ with .Date.Format "MST" }}{{ if eq . "UTC" }}UT{{else}}{{.}}{{end}}{{end}}</pubDate>
      {{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
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

	// Add SEO & Social metadata
	t.AddInternalTemplate("_default", "opengraph.html", `<meta property="og:title" content="{{ .Title }}" />
<meta property="og:description" content="{{ if .Description }}{{ .Description }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ end }}{{ end }}" />
<meta property="og:type" content="{{ if .IsPage }}article{{ else }}website{{ end }}" />
<meta property="og:url" content="{{ .Permalink }}" />
{{ with .Params.images }}{{ range first 6 . }}
  <meta property="og:image" content="{{ . }}" />
{{ end }}{{ end }}

<meta property="og:updated_time" content="{{ .Date }}"/>{{ with .Params.audio }}
<meta property="og:audio" content="{{ . }}" />{{ end }}{{ with .Params.locale }}
<meta property="og:locale" content="{{ . }}" />{{ end }}{{ with .Site.Params.title }}
<meta property="og:site_name" content="{{ . }}" />{{ end }}{{ with .Params.videos }}
{{ range .Params.videos }}
  <meta property="og:video" content="{{ . }}" />
{{ end }}

<!-- If it is part of a series, link to related articles -->
{{ $permalink := .Permalink }}
{{ $siteSeries := .Site.Taxonomies.series }}{{ with .Params.series }}
{{ range $name := . }}
  {{ $series := index $siteSeries $name }}
  {{ range $page := first 6 $series.Pages }}
    {{ if ne $page.Permalink $permalink }}<meta property="og:see_also" content="{{ $page.Permalink }}" />{{ end }}
  {{ end }}
{{ end }}{{ end }}

{{ if .IsPage }}
{{ range .Site.Authors }}{{ with .Social.facebook }}
<meta property="article:author" content="https://www.facebook.com/{{ . }}" />{{ end }}{{ with .Site.Social.facebook }}
<meta property="article:publisher" content="https://www.facebook.com/{{ . }}" />{{ end }}
<meta property="article:published_time" content="{{ .PublishDate }}" />
<meta property="article:modified_time" content="{{ .Date }}" />
<meta property="article:section" content="{{ .Section }}" />
{{ with .Params.tags }}{{ range first 6 . }}
  <meta property="article:tag" content="{{ . }}" />{{ end }}{{ end }}
{{ end }}

<!-- Facebook Page Admin ID for Domain Insights -->
{{ with .Site.Social.facebook_admin }}<meta property="fb:admins" content="{{ . }}" />{{ end }}`)

	t.AddInternalTemplate("_default", "twitter_cards.html", `{{ if .IsPage }}
{{ with .Params.images }}
<!-- Twitter summary card with large image must be at least 280x150px -->
  <meta name="twitter:card" content="summary_large_image"/>
  <meta name="twitter:image:src" content="{{ index . 0 }}"/>
{{ else }}
  <meta name="twitter:card" content="summary"/>
{{ end }}

<!-- Twitter Card data -->
<meta name="twitter:title" content="{{ .Title }}"/>
<meta name="twitter:description" content="{{ if .Description }}{{ .Description }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ end }}{{ end }}"/>
{{ with .Site.Social.twitter }}<meta name="twitter:site" content="@{{ . }}"/>{{ end }}
{{ with .Site.Social.twitter_domain }}<meta name="twitter:domain" content="{{ . }}"/>{{ end }}
{{ range .Site.Authors }}
  {{ with .twitter }}<meta name="twitter:creator" content="@{{ . }}"/>{{ end }}
{{ end }}{{ end }}`)

	t.AddInternalTemplate("_default", "google_news.html", `{{ if .IsPage }}{{ with .Params.news_keywords }}
  <meta name="news_keywords" content="{{ range $i, $kw := first 10 . }}{{ if $i }},{{ end }}{{ $kw }}{{ end }}" />
{{ end }}{{ end }}`)

	t.AddInternalTemplate("_default", "schema.html", `{{ with .Site.Social.GooglePlus }}<link rel="publisher" href="{{ . }}"/>{{ end }}
<meta itemprop="name" content="{{ .Title }}">
<meta itemprop="description" content="{{ if .Description }}{{ .Description }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ end }}{{ end }}">

{{if .IsPage}}
<meta itemprop="datePublished" content="{{ .PublishDate }}" />
<meta itemprop="dateModified" content="{{ .Date }}" />
<meta itemprop="wordCount" content="{{ .WordCount }}">
{{ with .Params.images }}{{ range first 6 . }}
  <meta itemprop="image" content="{{ . }}">
{{ end }}{{ end }}

<!-- Output all taxonomies as schema.org keywords -->
<meta itemprop="keywords" content="{{ range $plural, $terms := .Site.Taxonomies }}{{ range $term, $val := $terms }}{{ printf "%s," $term }}{{ end }}{{ end }}" />
{{ end }}{{ end }}`)

}
