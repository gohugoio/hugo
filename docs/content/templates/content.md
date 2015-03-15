---
aliases:
- /layout/content/
date: 2013-07-01
linktitle: Single Content
menu:
  main:
    parent: layout
next: /templates/list
prev: /templates/variables
title: Single Content Template
weight: 30
---

The primary view of content in Hugo is the single view. Hugo, for every
Markdown file provided, will render it with a single template.


## Which Template will be rendered?
Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present,
then the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
than necessary. For most sites, only the `_default` file at the end of
the list will be needed.

Users can specify the `type` and `layout` in the [front-matter](/content/front-matter/). `Section`
is determined based on the content file’s location. If `type` is provide,
it will be used instead of `section`.

### Single

* /layouts/`TYPE`-or-`SECTION`/`LAYOUT`.html
* /layouts/`TYPE`-or-`SECTION`/single.html
* /layouts/\_default/single.html
* /themes/`THEME`/layouts/`TYPE`-or-`SECTION`/`LAYOUT`.html
* /themes/`THEME`/layouts/`TYPE`-or-`SECTION`/single.html
* /themes/`THEME`/layouts/\_default/single.html

## Example Single Template File

Content pages are of the type "page" and have all the [page
variables](/layout/variables/) and [site
variables](/templates/variables/) available to use in the templates.

In the following examples we have created two different content types as well as
a default content type.

The default content template to be used in the event that a specific
template has not been provided for that type. The default type works the
same as the other types, but the directory must be called "\_default".

    ▾ layouts/
      ▾ _default/
          single.html
      ▾ post/
          single.html
      ▾ project/
          single.html


## post/single.html
This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates](/templates/partials/)

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}
    {{ $baseurl := .Site.BaseUrl }}

    <section id="main">
      <h1 id="title">{{ .Title }}</h1>
      <div>
            <article id="content">
               {{ .Content }}
            </article>
      </div>
    </section>

    <aside id="meta">
        <div>
        <section>
          <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
          <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
        </section>
        <ul id="categories">
          {{ range .Params.topics }}
            <li><a href="{{ $baseurl }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        <ul id="tags">
          {{ range .Params.tags }}
            <li> <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        </div>
        <div>
            {{ if .Prev }}
              <a class="previous" href="{{.Prev.Permalink}}"> {{.Prev.Title}}</a>
            {{ end }}
            {{ if .Next }}
              <a class="next" href="{{.Next.Permalink}}"> {{.Next.Title}}</a>
            {{ end }}
        </div>
    </aside>

    {{ partial "disqus.html" . }}
    {{ partial "footer.html" . }}


## project/single.html
This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates](/templates/partials/)


    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}
    {{ $baseurl := .Site.BaseUrl }}

    <section id="main">
      <h1 id="title">{{ .Title }}</h1>
      <div>
            <article id="content">
               {{ .Content }}
            </article>
      </div>
    </section>

    <aside id="meta">
        <div>
        <section>
          <h4 id="date"> {{ .Date.Format "Mon Jan 2, 2006" }} </h4>
          <h5 id="wc"> {{ .FuzzyWordCount }} Words </h5>
        </section>
        <ul id="categories">
          {{ range .Params.topics }}
          <li><a href="{{ $baseurl }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        <ul id="tags">
          {{ range .Params.tags }}
            <li> <a href="{{ $baseurl }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
          {{ end }}
        </ul>
        </div>
    </aside>

    {{if isset .Params "project_url" }}
    <div id="ribbon">
        <a href="{{ index .Params "project_url" }}" rel="me">Fork me on GitHub</a>
    </div>
    {{ end }}

    {{ partial "footer.html" . }}

Notice how the project/single.html template uses an additional parameter unique
to this template. This doesn't need to be defined ahead of time. If the key is
present in the front matter than it can be used in the template. To
easily generate new content of this type with these keys ready use
[content archetypes](/content/archetypes/).
