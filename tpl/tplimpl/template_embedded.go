// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package tplimpl

func (t *templateHandler) embedShortcodes() {
	t.addInternalShortcode("ref.html", `{{ if len .Params | eq 2 }}{{ ref .Page (.Get 0) (.Get 1) }}{{ else }}{{ ref .Page (.Get 0) }}{{ end }}`)
	t.addInternalShortcode("relref.html", `{{ if len .Params | eq 2 }}{{ relref .Page (.Get 0) (.Get 1) }}{{ else }}{{ relref .Page (.Get 0) }}{{ end }}`)
	t.addInternalShortcode("highlight.html", `{{ if len .Params | eq 2 }}{{ highlight .Inner (.Get 0) (.Get 1) }}{{ else }}{{ highlight .Inner (.Get 0) "" }}{{ end }}`)
	t.addInternalShortcode("test.html", `This is a simple Test`)
	t.addInternalShortcode("figure.html", `<!-- image -->
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
	t.addInternalShortcode("speakerdeck.html", "<script async class='speakerdeck-embed' data-id='{{ index .Params 0 }}' data-ratio='1.33333333333333' src='//speakerdeck.com/assets/embed.js'></script>")
	t.addInternalShortcode("youtube.html", `{{ if .IsNamedParams }}
<div {{ if .Get "class" }}class="{{ .Get "class" }}"{{ else }}style="position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;"{{ end }}>
  <iframe src="//www.youtube.com/embed/{{ .Get "id" }}?{{ with .Get "autoplay" }}{{ if eq . "true" }}autoplay=1{{ end }}{{ end }}" 
  {{ if not (.Get "class") }}style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" {{ end }}allowfullscreen frameborder="0"></iframe>
</div>{{ else }}
<div {{ if len .Params | eq 2 }}class="{{ .Get 1 }}"{{ else }}style="position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;"{{ end }}>
  <iframe src="//www.youtube.com/embed/{{ .Get 0 }}" {{ if len .Params | eq 1 }}style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" {{ end }}allowfullscreen frameborder="0"></iframe>
 </div>
{{ end }}`)
	t.addInternalShortcode("vimeo.html", `{{ if .IsNamedParams }}<div {{ if .Get "class" }}class="{{ .Get "class" }}"{{ else }}style="position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;"{{ end }}>
  <iframe src="//player.vimeo.com/video/{{ .Get "id" }}" {{ if not (.Get "class") }}style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" {{ end }}webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
 </div>{{ else }}
<div {{ if len .Params | eq 2 }}class="{{ .Get 1 }}"{{ else }}style="position: relative; padding-bottom: 56.25%; padding-top: 30px; height: 0; overflow: hidden;"{{ end }}>
  <iframe src="//player.vimeo.com/video/{{ .Get 0 }}" {{ if len .Params | eq 1 }}style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;" {{ end }}webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
 </div>
{{ end }}`)
	t.addInternalShortcode("gist.html", `<script src="//gist.github.com/{{ index .Params 0 }}/{{ index .Params 1 }}.js{{if len .Params | eq 3 }}?file={{ index .Params 2 }}{{end}}"></script>`)
	t.addInternalShortcode("tweet.html", `{{ (getJSON "https://api.twitter.com/1/statuses/oembed.json?id=" (index .Params 0)).html | safeHTML }}`)
	t.addInternalShortcode("instagram.html", `{{ if len .Params | eq 2 }}{{ if eq (.Get 1) "hidecaption" }}{{ with getJSON "https://api.instagram.com/oembed/?url=https://instagram.com/p/" (index .Params 0) "/&hidecaption=1" }}{{ .html | safeHTML }}{{ end }}{{ end }}{{ else }}{{ with getJSON "https://api.instagram.com/oembed/?url=https://instagram.com/p/" (index .Params 0) "/&hidecaption=0" }}{{ .html | safeHTML }}{{ end }}{{ end }}`)
}

func (t *templateHandler) embedTemplates() {

	t.addInternalTemplate("_default", "rss.xml", `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ if eq  .Title  .Site.Title }}{{ .Site.Title }}{{ else }}{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}{{ end }}</title>
    <link>{{ .Permalink }}</link>
    <description>Recent content {{ if ne  .Title  .Site.Title }}{{ with .Title }}in {{.}} {{ end }}{{ end }}on {{ .Site.Title }}</description>
    <generator>Hugo -- gohugo.io</generator>{{ with .Site.LanguageCode }}
    <language>{{.}}</language>{{end}}{{ with .Site.Author.email }}
    <managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}{{ with .Site.Author.email }}
    <webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}{{ with .Site.Copyright }}
    <copyright>{{.}}</copyright>{{end}}{{ if not .Date.IsZero }}
    <lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</lastBuildDate>{{ end }}
    {{ with .OutputFormats.Get "RSS" }}
	{{ printf "<atom:link href=%q rel=\"self\" type=%q />" .Permalink .MediaType | safeHTML }}
    {{ end }}
    {{ range .Data.Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</pubDate>
      {{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Summary | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>`)

	t.addInternalTemplate("_default", "sitemap.xml", `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml">
  {{ range .Data.Pages }}
  <url>
    <loc>{{ .Permalink }}</loc>{{ if not .Lastmod.IsZero }}
    <lastmod>{{ safeHTML ( .Lastmod.Format "2006-01-02T15:04:05-07:00" ) }}</lastmod>{{ end }}{{ with .Sitemap.ChangeFreq }}
    <changefreq>{{ . }}</changefreq>{{ end }}{{ if ge .Sitemap.Priority 0.0 }}
    <priority>{{ .Sitemap.Priority }}</priority>{{ end }}{{ if .IsTranslated }}{{ range .Translations }}
    <xhtml:link
                rel="alternate"
                hreflang="{{ .Lang }}"
                href="{{ .Permalink }}"
                />{{ end }}
    <xhtml:link
                rel="alternate"
                hreflang="{{ .Lang }}"
                href="{{ .Permalink }}"
                />{{ end }}
  </url>
  {{ end }}
</urlset>`)

	// For multilanguage sites
	t.addInternalTemplate("_default", "sitemapindex.xml", `<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	{{ range . }}
	<sitemap>
	   	<loc>{{ .SitemapAbsURL }}</loc>
		{{ if not .LastChange.IsZero }}
	   	<lastmod>{{ .LastChange.Format "2006-01-02T15:04:05-07:00" | safeHTML }}</lastmod>
		{{ end }}
	</sitemap>
	{{ end }}
</sitemapindex>
`)

	t.addInternalTemplate("", "pagination.html", `{{ $pag := $.Paginator }}
{{ if gt $pag.TotalPages 1 }}
<ul class="pagination">
    {{ with $pag.First }}
    <li>
        <a href="{{ .URL }}" aria-label="First"><span aria-hidden="true">&laquo;&laquo;</span></a>
    </li>
    {{ end }}
    <li
    {{ if not $pag.HasPrev }}class="disabled"{{ end }}>
    <a href="{{ if $pag.HasPrev }}{{ $pag.Prev.URL }}{{ end }}" aria-label="Previous"><span aria-hidden="true">&laquo;</span></a>
    </li>
    {{ $.Scratch.Set "__paginator.ellipsed" false }}
    {{ range $pag.Pagers }}
    {{ $right := sub .TotalPages .PageNumber }}
    {{ $showNumber := or (le .PageNumber 3) (eq $right 0) }}
    {{ $showNumber := or $showNumber (and (gt .PageNumber (sub $pag.PageNumber 2)) (lt .PageNumber (add $pag.PageNumber 2)))  }}
    {{ if $showNumber }} 
        {{ $.Scratch.Set "__paginator.ellipsed" false }}
        {{ $.Scratch.Set "__paginator.shouldEllipse" false }}
    {{ else }}
        {{ $.Scratch.Set "__paginator.shouldEllipse" (not ($.Scratch.Get "__paginator.ellipsed") ) }}
        {{ $.Scratch.Set "__paginator.ellipsed" true }}
    {{ end }}
    {{ if $showNumber }}
    <li
    {{ if eq . $pag }}class="active"{{ end }}><a href="{{ .URL }}">{{ .PageNumber }}</a></li>
    {{ else if ($.Scratch.Get "__paginator.shouldEllipse") }}
    <li class="disabled"><span aria-hidden="true">&hellip;</span></li>
    {{ end }}
    {{ end }}
    <li
    {{ if not $pag.HasNext }}class="disabled"{{ end }}>
    <a href="{{ if $pag.HasNext }}{{ $pag.Next.URL }}{{ end }}" aria-label="Next"><span aria-hidden="true">&raquo;</span></a>
    </li>
    {{ with $pag.Last }}
    <li>
        <a href="{{ .URL }}" aria-label="Last"><span aria-hidden="true">&raquo;&raquo;</span></a>
    </li>
    {{ end }}
</ul>
{{ end }}`)

	t.addInternalTemplate("", "disqus.html", `{{ if .Site.DisqusShortname }}<div id="disqus_thread"></div>
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
	t.addInternalTemplate("", "opengraph.html", `<meta property="og:title" content="{{ .Title }}" />
<meta property="og:description" content="{{ with .Description }}{{ . }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ else }}{{ with .Site.Params.description }}{{ . }}{{ end }}{{ end }}{{ end }}" />
<meta property="og:type" content="{{ if .IsPage }}article{{ else }}website{{ end }}" />
<meta property="og:url" content="{{ .Permalink }}" />
{{ with .Params.images }}{{ range first 6 . }}
  <meta property="og:image" content="{{ . | absURL }}" />
{{ end }}{{ end }}

{{ if .IsPage }}
{{ if not .PublishDate.IsZero }}<meta property="article:published_time" content="{{ .PublishDate.Format "2006-01-02T15:04:05-07:00" | safeHTML }}"/>
{{ else if not .Date.IsZero }}<meta property="article:published_time" content="{{ .Date.Format "2006-01-02T15:04:05-07:00" | safeHTML }}"/>{{ end }}
{{ if not .Lastmod.IsZero }}<meta property="article:modified_time" content="{{ .Lastmod.Format "2006-01-02T15:04:05-07:00" | safeHTML }}"/>{{ end }}
{{ else }}
{{ if not .Date.IsZero }}<meta property="og:updated_time" content="{{ .Date.Format "2006-01-02T15:04:05-07:00" | safeHTML }}"/>{{ end }}
{{ end }}{{ with .Params.audio }}
<meta property="og:audio" content="{{ . }}" />{{ end }}{{ with .Params.locale }}
<meta property="og:locale" content="{{ . }}" />{{ end }}{{ with .Site.Params.title }}
<meta property="og:site_name" content="{{ . }}" />{{ end }}{{ with .Params.videos }}
{{ range .Params.videos }}
  <meta property="og:video" content="{{ . | absURL }}" />
{{ end }}{{ end }}

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
<meta property="article:section" content="{{ .Section }}" />
{{ with .Params.tags }}{{ range first 6 . }}
  <meta property="article:tag" content="{{ . }}" />{{ end }}{{ end }}
{{ end }}{{ end }}

<!-- Facebook Page Admin ID for Domain Insights -->
{{ with .Site.Social.facebook_admin }}<meta property="fb:admins" content="{{ . }}" />{{ end }}`)

	t.addInternalTemplate("", "twitter_cards.html", `{{ if .IsPage }}
{{ with .Params.images }}
<!-- Twitter summary card with large image must be at least 280x150px -->
  <meta name="twitter:card" content="summary_large_image"/>
  <meta name="twitter:image:src" content="{{ index . 0 | absURL }}"/>
{{ else }}
  <meta name="twitter:card" content="summary"/>
{{ end }}

<!-- Twitter Card data -->
<meta name="twitter:text:title" content="{{ .Title }}"/>
<meta name="twitter:title" content="{{ .Title }}"/>
<meta name="twitter:description" content="{{ with .Description }}{{ . }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ else }}{{ with .Site.Params.description }}{{ . }}{{ end }}{{ end }}{{ end }}"/>
{{ with .Site.Social.twitter }}<meta name="twitter:site" content="@{{ . }}"/>{{ end }}
{{ range .Site.Authors }}
  {{ with .twitter }}<meta name="twitter:creator" content="@{{ . }}"/>{{ end }}
{{ end }}{{ end }}`)

	t.addInternalTemplate("", "google_news.html", `{{ if .IsPage }}{{ with .Params.news_keywords }}
  <meta name="news_keywords" content="{{ range $i, $kw := first 10 . }}{{ if $i }},{{ end }}{{ $kw }}{{ end }}" />
{{ end }}{{ end }}`)

	t.addInternalTemplate("", "schema.html", `{{ with .Site.Social.GooglePlus }}<link rel="publisher" href="{{ . }}"/>{{ end }}
<meta itemprop="name" content="{{ .Title }}">
<meta itemprop="description" content="{{ with .Description }}{{ . }}{{ else }}{{if .IsPage}}{{ .Summary }}{{ else }}{{ with .Site.Params.description }}{{ . }}{{ end }}{{ end }}{{ end }}">

{{if .IsPage}}{{ $ISO8601 := "2006-01-02T15:04:05-07:00" }}{{ if not .PublishDate.IsZero }}
<meta itemprop="datePublished" content="{{ .PublishDate.Format $ISO8601 | safeHTML }}" />{{ end }}
{{ if not .Date.IsZero }}<meta itemprop="dateModified" content="{{ .Date.Format $ISO8601 | safeHTML }}" />{{ end }}
<meta itemprop="wordCount" content="{{ .WordCount }}">
{{ with .Params.images }}{{ range first 6 . }}
  <meta itemprop="image" content="{{ . | absURL }}">
{{ end }}{{ end }}

<!-- Output all taxonomies as schema.org keywords -->
<meta itemprop="keywords" content="{{ range $plural, $terms := .Site.Taxonomies }}{{ range $term, $val := $terms }}{{ printf "%s," $term }}{{ end }}{{ end }}" />
{{ end }}`)

	t.addInternalTemplate("", "google_analytics.html", `{{ with .Site.GoogleAnalytics }}
<script>
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

ga('create', '{{ . }}', 'auto');
ga('send', 'pageview');
</script>
{{ end }}`)

	t.addInternalTemplate("", "google_analytics_async.html", `{{ with .Site.GoogleAnalytics }}
<script>
window.ga=window.ga||function(){(ga.q=ga.q||[]).push(arguments)};ga.l=+new Date;
ga('create', '{{ . }}', 'auto');
ga('send', 'pageview');
</script>
<script async src='//www.google-analytics.com/analytics.js'></script>
{{ end }}`)

	t.addInternalTemplate("_default", "robots.txt", "User-agent: *")
}
