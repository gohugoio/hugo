---
title: Template types
description: Create templates of different types to render your content, resources, and data.
categories: []
keywords: []
weight: 30
aliases: [
  '/templates/base/',
  '/templates/content-view/',
  '/templates/home/',
  '/templates/lists/',
  '/templates/partial/',
  '/templates/section/',
  '/templates/single/',
  '/templates/taxonomy/',
  '/templates/term/',
]
---

## Structure

Create templates in the `layouts` directory in the root of your project.

Although your site may not require each of these templates, the example below is typical for a site of medium complexity.

```text
layouts/
├── _markup/
│   ├── render-image.html   <-- render hook
│   └── render-link.html    <-- render hook
├── _partials/
│   ├── footer.html
│   └── header.html
├── _shortcodes/
│   ├── audio.html
│   └── video.html
├── books/
│   ├── page.html
│   └── section.html
├── films/
│   ├── card.html           <-- content view
│   ├── page.html
│   └── section.html
├── baseof.html
├── home.html
├── page.html
├── section.html
├── taxonomy.html
└── term.html
```

Hugo's [template lookup order] determines the template path, allowing you to create unique templates for any page.

> [!note]
> You must have thorough understanding of the template lookup order when creating templates. Template selection is based on template type, page kind, content type, section, language, and output format.

The purpose of each template type is described below.

## Base

A base template reduces duplicate code by wrapping other templates within a shell.

For example, the base template below calls the [`partial`] function to include partial templates for the `head`, `header`, and `footer` elements of each page, and it calls the [`block`] function to include `home`, `page`, `section`, `taxonomy`, and `term` templates within the `main` element of each page.

```go-html-template {file="layouts/baseof.html"}
<!DOCTYPE html>
<html lang="{{ or site.Language.LanguageCode }}" dir="{{ or site.Language.LanguageDirection `ltr` }}">
<head>
  {{ partial "head.html" . }}
</head>
<body>
  <header>
    {{ partial "header.html" . }}
  </header>
  <main>
    {{ block "main" . }}{{ end }}
  </main>
  <footer>
    {{ partial "footer.html" . }}
  </footer>
</body>
</html>
```

The `block` construct above is used to define a set of root templates that are then customized by redefining the block templates within. See&nbsp;[details](/functions/go-template/block/)

## Home

A home template renders your site's home page. For example, the home template below inherits the site's shell from the [base template] and renders the home page content, such as a list of other pages.

```go-html-template {file="layouts/home.html"}
{{ define "main" }}
  {{ .Content }}
  {{ range .Site.RegularPages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

## Page

A page template renders a regular page.

For example, the page template below inherits the site's shell from the [base template] and renders the page title and page content.

```go-html-template {file="layouts/page.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

## Section

A section template renders a list of pages within a section.

For example, the section template below inherits the site's shell from the [base template] and renders a list of pages in the current section.

```go-html-template {file="layouts/section.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

## Taxonomy

A taxonomy template renders a list of terms in a [taxonomy](g).

For example, the taxonomy template below inherits the site's shell from the [base template] and renders a list of terms in the current taxonomy.

```go-html-template {file="layouts/taxonomy.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Within a taxonomy template, the [`Data`] object provides these taxonomy-specific methods:

- [`Singular`][taxonomy-singular]
- [`Plural`][taxonomy-plural]
- [`Terms`].

The `Terms` method returns a [taxonomy object](g), allowing you to call any of its methods including [`Alphabetical`] and [`ByCount`]. For example, use the `ByCount` method to render a list of terms sorted by the number of pages associated with each term:

```go-html-template {file="layouts/taxonomy.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Data.Terms.ByCount }}
    <h2><a href="{{ .Page.RelPermalink }}">{{ .Page.LinkTitle }}</a> ({{ .Count }})</h2>
  {{ end }}
{{ end }}
```

## Term

A term template renders a list of pages associated with a [term](g).

For example, the term template below inherits the site's shell from the [base template] and renders a list of pages associated with the current term.

```go-html-template {file="layouts/term.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Within a term template, the [`Data`] object provides these term-specific methods:

- [`Singular`][term-singular]
- [`Plural`][term-plural]
- [`Term`].

## Single

A single template is a fallback for [page templates](#page). If a page template does not exist, Hugo will look for a single template instead.

Like a page template, a single template renders a regular page.

For example, the single template below inherits the site's shell from the [base template] and renders the page title and page content.

```go-html-template {file="layouts/single.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

## List

A list template is a fallback for these template types: [home](#home), [section](#section), [taxonomy](#taxonomy), and [term](#term). If one of these template types does not exist, Hugo will look for a list template instead.

For example, the list template below inherits the site's shell from the [base template] and renders a list of pages:

```go-html-template {file="layouts/list.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

## All

An "all" template is a fallback for these template types: [home](#home), [page](#page), [section](#section), [taxonomy](#taxonomy), [term](#term), [single](#single), and [list](#list). If one of these template types does not exist, Hugo will look for an "all" template instead.

For example, the contrived "all" template below inherits the site's shell from the [base template] and conditionally renders a page based on its page kind:

```go-html-template {file="layouts/all.html"}
{{ define "main" }}
  {{ if eq .Kind "home" }}
    {{ .Content }}
    {{ range .Site.RegularPages }}
      <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
    {{ end }}
  {{ else if eq .Kind "page" }}
    <h1>{{ .Title }}</h1>
    {{ .Content }}
  {{ else if in (slice "section" "taxonomy" "term") .Kind }}
    <h1>{{ .Title }}</h1>
    {{ .Content }}
    {{ range .Pages }}
      <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
    {{ end }}
  {{ else }}
    {{ errorf "Unsupported page kind: %s" .Kind }}
  {{ end }}
{{ end }}
```

## Partial

A partial template is typically used to render a component of your site, though you may also create partial templates that return values.


For example, the partial template below renders copyright information:

```go-html-template {file="layouts/_partials/footer.html"}
<p>Copyright {{ now.Year }}. All rights reserved.</p>
```

Execute the partial template by calling the [`partial`] or [`partialCached`] function, optionally passing context as the second argument:

```go-html-template {file="layouts/baseof.html"}
{{ partial "footer.html" . }}
```

Unlike other template types, partial template selection is based on the file name passed in the partial call. Hugo does not consider the current page kind, content type, logical path, language, or output format when searching for a matching partial template. However, Hugo _does_ apply the same name matching logic it uses for other templates. This means it tries to find the most specific match first, then progressively looks for more general versions if the specific one isn't found.

For example, with this partial call:

```go-html-template {file="layouts/baseof.html"}
{{ partial "footer.section.de.html" . }}
```

Hugo uses this lookup order to find a matching template:

1. `layouts/_partials/footer.section.de.html`
1. `layouts/_partials/footer.section.html`
1. `layouts/_partials/footer.de.html`
1. `layouts/_partials/footer.html`

Partials can also be defined inline within a template. However, it's important to note that the template namespace is global; ensuring unique names for these partials is necessary to prevent conflicts.

```go-html-template
Value: {{ partial "my-inline-partial.html" . }}

{{ define "_partials/my-inline-partial.html" }}
  {{ $value := 32 }}
  {{ return $value }}
{{ end }}
```

## Content view

A content view template is similar to a partial template, invoked by calling the [`Render`] method on a `Page` object. Unlike partial templates, content view templates:

- Inherit the context of the current page
- Can target any page kind, content type, logical path, language, or output format

For example, the home template below inherits the site's shell from the [base template], and renders a card component for each page within the "films" section of your site.

```go-html-template {file="layouts/home.html"}
{{ define "main" }}
  {{ .Content }}
  <ul>
    {{ range where site.RegularPages "Section" "films" }}
      {{ .Render "card" }}
    {{ end }}
  </ul>
{{ end }}
```

```go-html-template {file="layouts/films/card.html"}
<div class="card">
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
</div>
```

## Render hook

A render hook template overrides the conversion of Markdown to HTML.

For example, the render hook template below adds an anchor link to the right of each heading.

```go-html-template {file="layouts/_markup/heading.html"}
<h{{ .Level }} id="{{ .Anchor }}" {{- with .Attributes.class }} class="{{ . }}" {{- end }}>
  {{ .Text }}
  <a href="#{{ .Anchor }}">#</a>
</h{{ .Level }}>
```

Learn more about [render hook templates](/render-hooks/).

## Shortcode

A shortcode template is used to render a component of your site. Unlike [partial templates](#partial) or [content view templates](#content-view), shortcode templates are called from content pages.

For example, the shortcode template below renders an audio element from a [global resource](g).

```go-html-template {file="layouts/_shortcodes/audio.html"}
{{ with resources.Get (.Get "src") }}
  <audio controls preload="auto" src="{{ .RelPermalink }}"></audio>
{{ end }}
```

Then call the shortcode from within markup:

```text {file="content/example.md"}
{{</* audio src=/audio/test.mp3 */>}}
```

Learn more about [shortcode templates](/templates/shortcode/).

## Other

Use other specialized templates to create:

- [Sitemaps](/templates/sitemap)
- [RSS feeds](/templates/rss/)
- [404 error pages](/templates/404/)
- [robots.txt files](/templates/robots/)

[`Alphabetical`]: /methods/taxonomy/alphabetical/
[`block`]: /functions/go-template/block/
[`ByCount`]: /methods/taxonomy/bycount/
[`Data`]: /methods/page/data/
[`partial`]: /functions/partials/include/
[`partialCached`]: /functions/partials/includeCached/
[`Render`]: /methods/page/render/
[`Taxonomy`]: /methods/taxonomy/
[`Terms`]: /methods/page/data/#terms
[`Term`]: /methods/page/data/#term
[taxonomy-plural]: /methods/page/data/#plural
[taxonomy-singular]: /methods/page/data/#singular
[template lookup order]: /templates/lookup-order/
[term-plural]: /methods/page/data/#plural-1
[term-singular]: /methods/page/data/#singular-1
[base template]: #base
