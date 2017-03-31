---
title: Single Page Templates
linktitle:
description: The primary view of content in Hugo is the single view. Hugo will render every Markdown file provided with a corresponding single template.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [page]
menu:
  main:
    parent: "Templates"
    weight: 60
weight: 60	#rem
draft: false
aliases: [/layout/content/]
toc: true
---

The primary view of content in Hugo is the single view. Hugo's default behavior is to render every Markdown file provided with a corresponding single template.

## Single Page Template Lookup Order

You can specify `type` (i.e., [content type][]) and `layout` in a single content file's [front matter][]. However, you cannot specify `section` because this is determined based on file location (see [content section][section]).

Hugo assumes your content section and content type are the same unless you tell Hugo otherwise by providing a `type` directly in the front matter of a content file. This is why #1 and #3 come before #2 and #4, respectively, in the following lookup order. Values in angle brackets (`<>`) are variable.

1. `/layouts/<TYPE>/<LAYOUT>.html`
2. `/layouts/<SECTION>>/<LAYOUT>.html`
3. `/layouts/<TYPE>/single.html`
4. `/layouts/<SECTION>/single.html`
5. `/layouts/_default/single.html`
6. `/themes/<THEME>/layouts/<TYPE>/<LAYOUT>.html`
7. `/themes/<THEME>/layouts/<SECTION>/<LAYOUT>.html`
8. `/themes/<THEME>/layouts/<TYPE>/single.html`
9. `/themes/<THEME>/layouts/<SECTION>/single.html`
10. `/themes/<THEME>/layouts/_default/single.html`

{{% note %}}
`my-first-event.md` is significant because it demonstrates the role of the lookup order in Hugo themes. Both the root project directory *and* the `mytheme` themes directory have a file at `_default/single.html`. Understanding this order allows you to [customize Hugo themes](/themes/customizing/) by creating template files with identical names in your project directory that step in front of theme template files in the lookup. This allows you to customize the look and feel of your website while maintaining compatibility with the theme's upstream.
{{% /note %}}

## Example Single Page Templates

Content pages are of the type `page` and will therefore have all the [page variables][] and [site variables][] available to use in their templates.

### `post/single.html`

This content template is used for [spf13.com][spf13]. It makes use of [partial templates][partials]:

{{% code file="layouts/post/single.html" download="single.html" %}}
```html
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}
{{ $baseURL := .Site.BaseURL }}
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
        <li><a href="{{ $baseURL }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
      {{ end }}
    </ul>
    <ul id="tags">
      {{ range .Params.tags }}
        <li> <a href="{{ $baseURL }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
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
```
{{% /code %}}

### `project/single.html`

This content template is also used for [spf13.com][spf13] and makes use of [partial templates][partials]:

{{% code file="project/single.html" download="single.html" %}}
```html
  {{ partial "header.html" . }}
  {{ partial "subheader.html" . }}
  {{ $baseURL := .Site.BaseURL }}

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
        <li><a href="{{ $baseURL }}/topics/{{ . | urlize }}">{{ . }}</a> </li>
        {{ end }}
      </ul>
      <ul id="tags">
        {{ range .Params.tags }}
          <li> <a href="{{ $baseURL }}/tags/{{ . | urlize }}">{{ . }}</a> </li>
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
```
{{% /code %}}

Notice how `project/single.html` uses an additional parameter unique to this template. This doesn't need to be defined ahead of time. The key can wait to be used in the template if present in the content file's front matter.

To easily generate new instances of this content type (e.g., new `.md` files in `project/`) with preconfigured front matter, use [content archetypes][archetypes].

[archetypes]: /content-management/archetypes/
[config]: /getting-started/configuration/
[content type]: /content-management/types/
[directory structure]: /getting-started/directory-structure/
[dry]: https://en.wikipedia.org/wiki/Don%27t_repeat_yourself
[front matter]: /content-management/front-matter/
[page variables]: /variables/page/
[partials]: /templates/partials/
[section]: /content-management/sections/
[site variables]: /variables/site/
[spf13]: http://spf13.com/