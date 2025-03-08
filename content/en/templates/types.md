---
title: Template types
description: Create templates of different types to render your content, resources, and data.
categories: []
keywords: []
weight: 30
aliases: ['/templates/lists/']
---

## Structure

Create templates in the `layouts` directory in the root of your project.

Although your site may not require each of these templates, the example below is typical for a site of medium complexity.

```text
layouts/
├── _default/
│   ├── _markup/
│   │   ├── render-image.html   <-- render hook
│   │   └── render-link.html    <-- render hook
│   ├── baseof.html
│   ├── home.html
│   ├── section.html
│   ├── single.html
│   ├── taxonomy.html
│   └── term.html
├── articles/
│   └── card.html               <-- content view
├── partials/
│   ├── footer.html
│   └── header.html
└── shortcodes/
    ├── audio.html
    └── video.html
```

Hugo's [template lookup order] determines the template path, allowing you to create unique templates for any page.

> [!note]
> You must have thorough understanding of the template lookup order when creating templates. Template selection is based on template type, page kind, content type, section, language, and output format.

The purpose of each template type is described below.

## Base

Base templates reduce duplicate code by wrapping other templates within a shell.

For example, the base template below calls the [partial] function to include partial templates for the `head`, `header`, and `footer` elements of each page, and it uses the [block] function to include `home`, `single`, `section`, `taxonomy`, and `term` templates within the `main` element of each page.

```go-html-template {file="layouts/_default/baseof.html"}
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

Learn more about [base templates](/templates/base/).

## Home

A home page template is used to render your site's home page, and is the only template required for a single-page website. For example, the home page template below inherits the site's shell from the base template and renders the home page content, such as a list of other pages.

```go-html-template {file="layouts/_default/home.html"}
{{ define "main" }}
  {{ .Content }}
  {{ range site.RegularPages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Learn more about [home page templates](/templates/home/).

## Single

A single template renders a single page.

For example, the single template below inherits the site's shell from the base template, and renders the title and content of each page.

```go-html-template {file="layouts/_default/single.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

Learn more about [single templates](/templates/single/).

## Section

A section template typically renders a list of pages within a section.

For example, the section template below inherits the site's shell from the base template, and renders a list of pages in the current section.

```go-html-template {file="layouts/_default/section.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Learn more about [section templates](/templates/section/).

## Taxonomy

A taxonomy template renders a list of terms in a [taxonomy](g).

For example, the taxonomy template below inherits the site's shell from the base template, and renders a list of terms in the current taxonomy.

```go-html-template {file="layouts/_default/taxonomy.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Learn more about [taxonomy templates](/templates/taxonomy/).

## Term

A term template renders a list of pages associated with a [term](g).

For example, the term template below inherits the site's shell from the base template, and renders a list of pages associated with the current term.

```go-html-template {file="layouts/_default/term.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/filter-sort-group.md" %}}

Learn more about [term templates](/templates/term/).

## Partial

A partial template is typically used to render a component of your site, though you may also create partial templates that return values.

> [!note]
> Unlike other template types, you cannot create partial templates to target a particular page kind, content type, section, language, or output format. Partial templates do not follow Hugo's [template lookup order].

For example, the partial template below renders copyright information.

```go-html-template {file="layouts/partials/footer.html"}
<p>Copyright {{ now.Year }}. All rights reserved.</p>
```

Learn more about [partial templates](/templates/partial/).

## Content view

A content view template is similar to a partial template, invoked by calling the [`Render`] method on a `Page` object. Unlike partial templates, content view templates:

- Automatically inherit the context of the current page
- Follow a lookup order allowing you to target a given content type or section

For example, the home template below inherits the site's shell from the base template, and renders a card component for each page within the "articles" section of your site.

```go-html-template {file="layouts/_default/home.html"}
{{ define "main" }}
  {{ .Content }}
  <ul>
    {{ range where site.RegularPages "Section" "articles" }}
      {{ .Render "card" }}
    {{ end }}
  </ul>
{{ end }}
```

```go-html-template {file="layouts/articles/card.html"}
<div class="card">
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Summary }}
</div>
```

Learn more about [content view templates](/templates/content-view/).

## Render hook

A render hook template overrides the conversion of Markdown to HTML.

For example, the render hook template below adds a `rel` attribute to external links.

```go-html-template {file="layouts/_default/_markup/render-link.html"}
{{- $u := urls.Parse .Destination -}}
<a href="{{ .Destination | safeURL }}"
  {{- with .Title }} title="{{ . }}"{{ end -}}
  {{- if $u.IsAbs }} rel="external"{{ end -}}
>
  {{- with .Text }}{{ . }}{{ end -}}
</a>
{{- /* chomp trailing newline */ -}}
```

Learn more about [render hook templates](/render-hooks/).

## Shortcode

A shortcode template is used to render a component of your site. Unlike partial templates, shortcode templates are called from content pages.

For example, the shortcode template below renders an audio element from a [global resource](g).

```go-html-template {file="layouts/shortcodes/audio.html"}
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

[`Render`]: /methods/page/render/
[block]: /functions/go-template/block/
[partial]: /functions/partials/include/
[template lookup order]: /templates/lookup-order/
[template lookup order]: /templates/lookup-order/
