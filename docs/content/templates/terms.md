---
aliases:
- /indexes/lists/
- /doc/indexes/
- /extras/indexes
lastmod: 2015-09-15
date: 2014-05-21
linktitle: Taxonomies
menu:
  main:
    parent: layout
next: /templates/views
prev: /templates/homepage
title: Taxonomy Templates
weight: 60
toc: true
---

Two templates are needed to create lists of the terms for a given taxonomy. 
1. A Taxonomy terms template
2. A Taxonomy list template [list template](/templates/list/)

## Taxonomy terms template
A Taxonomy Terms template will be rendered at /`PLURAL`/
(e.g. http://spf13.com/topics/)

Hugo searches for Taxonomy terms templates in the following folders:
* /layouts/taxonomy/`SINGULAR`.terms.html (e.g. `/layouts/taxonomy/topic.terms.html`)
* /theme/layouts/taxonomy/`SINGULAR`.terms.html
* /layouts/\_default/terms.html

If a template isn’t present, then the next listed will be used. This enables you to craft
specific layouts when you want to without creating more templates than necessary. 
For most sites, only the `_default` file at the end of the list will be needed.

If no Taxonomy terms template is provided then Hugo will not automatically generate terms pages for it. It is also
common for people to render taxonomy terms lists on other pages such as
the homepage or the sidebar (such as a tag cloud) and not have a
dedicated page for the terms.

## Taxonomy list template 
A Taxonomy list template will be rendered at /`PLURAL/TERM/`/
(e.g. http://spf13.com/topics/development/)

Hugo searches for Taxonomy list templates in the following folders:
* /layouts/taxonomy/`SINGULAR`.html (e.g. `/layouts/taxonomy/topic.html`)

Taxonomy list templates share the same facilities as [list template](/templates/list/)

## Variables

Taxonomy Terms and List pages are of the type "node" and have all the
[node variables](/templates/variables/) and
[site variables](/templates/variables/)
available to use in the templates.

Taxonomy Terms pages will additionally have:

* **.Data.Singular** The singular name of the taxonomy
* **.Data.Plural** The plural name of the taxonomy
* **.Data.Terms** The taxonomy itself
* **.Data.Terms.Alphabetical** The Terms alphabetized
* **.Data.Terms.ByCount** The Terms ordered by popularity

The last two can also be reversed: **.Data.Terms.Alphabetical.Reverse**, **.Data.Terms.ByCount.Reverse**.

### Example terms.html files

List pages are of the type "node" and have all the
[node variables](/templates/variables/) and
[site variables](/templates/variables/)
available to use in the templates.

This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates](/templates/partials/). The list of taxonomy
templates cannot use a [content view](/templates/views/) as they don't display the content, but
rather information about the content.

This particular template lists all of the Tags used on
[spf13.com](http://spf13.com/) and provides a count for the number of pieces of
content tagged with each tag.

`.Data.Terms` is a map of terms ⇒ [contents]

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
        <h1 id="title">{{ .Title }}</h1>

        <ul>
        {{ $data := .Data }}
        {{ range $key, $value := .Data.Terms }}
          <li><a href="{{ $data.Plural }}/{{ $key | urlize }}">{{ $key }}</a> {{ len $value }}</li>
        {{ end }}
       </ul>
      </div>
    </section>

    {{ partial "footer.html" . }}


Another example listing the content for each term (ordered by Date):

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
        <h1 id="title">{{ .Title }}</h1>

        {{ $data := .Data }}
        {{ range $key,$value := .Data.Terms.ByCount }}
        <h2><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}">{{ $value.Name }}</a> {{ $value.Count }}</h2>
        <ul>
        {{ range $value.Pages.ByDate }}
          <li><a href="{{ .Permalink }}">{{ .Title }}</a></li>
        {{ end }}
        </ul>
        {{ end }}
      </div>
    </section>

    {{ partial "footer.html" . }}


## Ordering

Hugo can order the meta data in two different ways. It can be ordered:

* by the number of contents assigned to that key, or
* alphabetically.

### Example terms.html file (alphabetical)

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
        <h1 id="title">{{ .Title }}</h1>
        <ul>
        {{ $data := .Data }}
        {{ range $key, $value := .Data.Terms.Alphabetical }}
          <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}">{{ $value.Name }}</a> {{ $value.Count }}</li>
        {{ end }}
        </ul>
      </div>
    </section>
    {{ partial "footer.html" . }}

### Example terms.html file (ordered by popularity)

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
        <h1 id="title">{{ .Title }}</h1>
        <ul>
        {{ $data := .Data }}
        {{ range $key, $value := .Data.Terms.ByCount }}
          <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}">{{ $value.Name }}</a> {{ $value.Count }}</li>
        {{ end }}
        </ul>
      </div>
    </section>

    {{ partial "footer.html" . }}
