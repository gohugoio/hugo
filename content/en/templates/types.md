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
│   ├── view_card.html      <-- content view
│   ├── view_li.html        <-- content view
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

A _base_ template serves as a foundational layout that other templates can build upon. It typically defines the common structural components of your HTML, such as the `html`, `head`, and `body` elements. It also often includes recurring features like headers, footers, navigation, and script inclusions that appear across multiple pages of your site. By defining these common aspects once in a _base_ template, you avoid redundancy, ensure consistency, and simplify the maintenance of your website.

Hugo can apply a _base_ template to the following template types: [home](#home), [page](#page), [section](#section), [taxonomy](#taxonomy), [term](#term), [single](#single), [list](#list), and [all](#all). When Hugo parses any of these template types, it will apply a _base_ template only if the template being parsed meets these specific conditions:

- It must include at least one [`define`] [action](g).
- It can only contain `define` actions, whitespace, and [template comments]. No other content is allowed.

> [!note]
> If a template doesn't meet all these criteria, Hugo executes it exactly as provided, without applying a _base_ template.

When Hugo applies a _base_ template, it replaces its [`block`] actions with content from the corresponding `define` actions found in the template to which the base template is applied.

For example, the _base_ template below calls the [`partial`] function to include `head`, `header`, and `footer` elements. The `block` action acts as a placeholder, and its content will be replaced by a matching `define` action  from the template to which it is applied.

```go-html-template {file="layouts/baseof.html"}
<!DOCTYPE html>
<html lang="{{ site.Language.LanguageCode }}" dir="{{ or site.Language.LanguageDirection `ltr` }}">
<head>
  {{ partial "head.html" . }}
</head>
<body>
  <header>
    {{ partial "header.html" . }}
  </header>
  <main>
    {{ block "main" . }}
      This will be replaced with content from the 
      corresponding "define" action found in the template
      to which this base template is applied.
    {{ end }}
  </main>
  <footer>
    {{ partial "footer.html" . }}
  </footer>
</body>
</html>
```

```go-html-template {file="layouts/home.html"}
{{ define "main" }}
  This will replace the content of the "block" action
  found in the base template.
{{ end }}
```

## Home

A _home_ template renders your site's home page.

For example, Hugo applies a _base_ template to the _home_ template below, then renders the page content and a list of the site's regular pages.

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

A _page_ template renders a regular page.

For example, Hugo applies a _base_ template to the _page_ template below, then renders the page title and page content.

```go-html-template {file="layouts/page.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

## Section

A _section_ template renders a list of pages within a [section](g).

For example, Hugo applies a _base_ template to the _section_ template below, then renders the page title, page content, and a list of pages in the current section.

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

A _taxonomy_ template renders a list of terms in a [taxonomy](g).

For example, Hugo applies a _base_ template to the _taxonomy_ template below, then renders the page title, page content, and a list of [terms](g) in the current taxonomy.

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

Within a _taxonomy_ template, the [`Data`] object provides these taxonomy-specific methods:

- [`Singular`][taxonomy-singular]
- [`Plural`][taxonomy-plural]
- [`Terms`]

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

A _term_ template renders a list of pages associated with a [term](g).

For example, Hugo applies a _base_ template to the _term_ template below, then renders the page title, page content, and a list of pages associated with the current term.

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

Within a _term_ template, the [`Data`] object provides these term-specific methods:

- [`Singular`][term-singular]
- [`Plural`][term-plural]
- [`Term`]

## Single

A _single_ template is a fallback for a _page_ template. If a _page_ template does not exist, Hugo will look for a _single_ template instead.

For example, Hugo applies a _base_ template to the _single_ template below, then renders the page title and page content.

```go-html-template {file="layouts/single.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

## List

A _list_ template is a fallback for [home](#home), [section](#section), [taxonomy](#taxonomy), and [term](#term) templates. If one of these template types does not exist, Hugo will look for a _list_ template instead.

For example, Hugo applies a _base_ template to the _list_ template below, then renders the page title, page content, and a list of pages.

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

An _all_ template is a fallback for [home](#home), [page](#page), [section](#section), [taxonomy](#taxonomy), [term](#term), [single](#single), and [list](#list) templates. If one of these template types does not exist, Hugo will look for an _all_ template instead.

For example, Hugo applies a _base_ template to the _all_ template below, then conditionally renders a page based on its page kind.

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

A _partial_ template is typically used to render a component of your site, though you may also create _partial_ templates that return values.

For example, the _partial_ template below renders copyright information:

```go-html-template {file="layouts/_partials/footer.html"}
<p>Copyright {{ now.Year }}. All rights reserved.</p>
```

Execute the _partial_ template by calling the [`partial`] or [`partialCached`] function, optionally passing context as the second argument:

```go-html-template {file="layouts/baseof.html"}
{{ partial "footer.html" . }}
```

<!-- https://github.com/gohugoio/hugo/pull/13614#issuecomment-2805977008 -->
Unlike other template types, Hugo does not consider the current page kind, content type, logical path, language, or output format when searching for a matching _partial_ template. However, it _does_ apply the same name matching logic it uses for other template types. This means it tries to find the most specific match first, then progressively looks for more general versions if the specific one isn't found.

For example, with this call:

```go-html-template {file="layouts/baseof.html"}
{{ partial "footer.section.de.html" . }}
```

Hugo uses this lookup order to find a matching template:

1. `layouts/_partials/footer.section.de.html`
1. `layouts/_partials/footer.section.html`
1. `layouts/_partials/footer.de.html`
1. `layouts/_partials/footer.html`

A _partial_ template can also be defined inline within another template. However, it's important to note that the template namespace is global; ensuring unique names for these _partial_ templates is necessary to prevent conflicts.

```go-html-template
Value: {{ partial "my-inline-partial.html" . }}

{{ define "_partials/my-inline-partial.html" }}
  {{ $value := 32 }}
  {{ return $value }}
{{ end }}
```

## Content view

A _content view_ template is similar to a _partial_ template, invoked by calling the [`Render`] method on a `Page` object. Unlike _partial_ templates, _content view_ templates:

- Inherit the context of the current page
- Can target any page kind, content type, logical path, language, or output format

For example, Hugo applies a _base_ template to the _home_ template below, then renders the page content and a card component for each page within the "films" section of your site.

```go-html-template {file="layouts/home.html"}
{{ define "main" }}
  {{ .Content }}
  <ul>
    {{ range where site.RegularPages "Section" "films" }}
      {{ .Render "view_card" }}
    {{ end }}
  </ul>
{{ end }}
```

```go-html-template {file="layouts/films/view_card.html"}
<div class="card">
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
</div>
```

In the example above, the content view template's name starts with `view_`. While not strictly required, this naming convention helps distinguish content view templates from other templates within the same directory, improving organization and clarity.

## Render hook

A _render hook_ template overrides the conversion of Markdown to HTML.

For example, the _render hook_ template below adds an anchor link to the right of each heading.

```go-html-template {file="layouts/_markup/render-heading.html"}
<h{{ .Level }} id="{{ .Anchor }}" {{- with .Attributes.class }} class="{{ . }}" {{- end }}>
  {{ .Text }}
  <a href="#{{ .Anchor }}">#</a>
</h{{ .Level }}>
```

Learn more about [render hook templates](/render-hooks/).

## Shortcode

A _shortcode_ template is used to render a component of your site. Unlike _partial_ or _content view_ templates, _shortcode_ templates are called from content pages.

For example, the _shortcode_ template below renders an audio element from a [global resource](g).

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
[`define`]: /functions/go-template/define/
[`partial`]: /functions/partials/include/
[`partialCached`]: /functions/partials/includeCached/
[`Render`]: /methods/page/render/
[`Term`]: /methods/page/data/#term
[`Terms`]: /methods/page/data/#terms
[taxonomy-plural]: /methods/page/data/#plural
[taxonomy-singular]: /methods/page/data/#singular
[template comments]: /templates/introduction/#comments
[template lookup order]: /templates/lookup-order/
[term-plural]: /methods/page/data/#plural-1
[term-singular]: /methods/page/data/#singular-1
