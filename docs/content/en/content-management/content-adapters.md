---
title: Content adapters
description: Create content adapters to dynamically add content when building your site.
categories: []
keywords: []
---

{{< new-in 0.126.0 />}}

## Overview

A content adapter is a template that dynamically creates pages when building a site. For example, use a content adapter to create pages from a remote data source such as JSON, TOML, YAML, or XML.

Unlike templates that reside in the `layouts` directory, content adapters reside in the `content` directory, no more than one per directory per language. When a content adapter creates a page, the page's [logical path](g) will be relative to the content adapter.

```text
content/
├── articles/
│   ├── _index.md
│   ├── article-1.md
│   └── article-2.md
├── books/
│   ├── _content.gotmpl  <-- content adapter
│   └── _index.md
└── films/
    ├── _content.gotmpl  <-- content adapter
    └── _index.md
```

Each content adapter is named `_content.gotmpl` and uses the same [syntax] as templates in the `layouts` directory. You can use any of the [template functions] within a content adapter, as well as the methods described below.

## Methods

Use these methods within a content adapter.

### AddPage

Adds a page to the site.

```go-html-template {file="content/books/_content.gotmpl"}
{{ $content := dict
  "mediaType" "text/markdown"
  "value" "The _Hunchback of Notre Dame_ was written by Victor Hugo."
}}
{{ $page := dict
  "content" $content
  "kind" "page"
  "path" "the-hunchback-of-notre-dame"
  "title" "The Hunchback of Notre Dame"
}}
{{ .AddPage $page }}
```

### AddResource

Adds a page resource to the site.

```go-html-template {file="content/books/_content.gotmpl"}
{{ with resources.Get "images/a.jpg" }}
  {{ $content := dict
    "mediaType" .MediaType.Type
    "value" .
  }}
  {{ $resource := dict
    "content" $content
    "path" "the-hunchback-of-notre-dame/cover.jpg"
  }}
  {{ $.AddResource $resource }}
{{ end }}
```

Then retrieve the new page resource with something like:

```go-html-template {file="layouts/page.html"}
{{ with .Resources.Get "cover.jpg" }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

### Site

Returns the `Site` to which the pages will be added.

```go-html-template {file="content/books/_content.gotmpl"}
{{ .Site.Title }}
```

> [!note]
> Note that the `Site` returned isn't fully built when invoked from the content adapters; if you try to call methods that depends on pages, e.g. `.Site.Pages`, you will get an error saying "this method cannot be called before the site is fully initialized".

### Store

Returns a persistent "scratch pad" to store and manipulate data. The main use case for this is to transfer values between executions when [EnableAllLanguages](#enablealllanguages) is set. See [examples](/methods/page/store/).

```go-html-template {file="content/books/_content.gotmpl"}
{{ .Store.Set "key" "value" }}
{{ .Store.Get "key" }}
```

### EnableAllLanguages

By default, Hugo executes the content adapter for the language defined by the `_content.gotmpl` file. Use this method to activate the content adapter for all languages.

```go-html-template {file="content/books/_content.gotmpl"}
{{ .EnableAllLanguages }}
{{ $content := dict
  "mediaType" "text/markdown"
  "value" "The _Hunchback of Notre Dame_ was written by Victor Hugo."
}}
{{ $page := dict
  "content" $content
  "kind" "page"
  "path" "the-hunchback-of-notre-dame"
  "title" "The Hunchback of Notre Dame"
}}
{{ .AddPage $page }}
```

## Page map

Set any [front matter field] in the map passed to the [`AddPage`](#addpage) method, excluding `markup`. Instead of setting the `markup` field, specify the `content.mediaType` as described below.

This table describes the fields most commonly passed to the `AddPage` method.

Key|Description|Required
:--|:--|:-:
`content.mediaType`|The content [media type]. Default is `text/markdown`. See [content formats] for examples.|&nbsp;
`content.value`|The content value as a string.|&nbsp;
`dates.date`|The page creation date as a `time.Time` value.|&nbsp;
`dates.expiryDate`|The page expiry date as a `time.Time` value.|&nbsp;
`dates.lastmod`|The page last modification date as a `time.Time` value.|&nbsp;
`dates.publishDate`|The page publication date as a `time.Time` value.|&nbsp;
`params`|A map of page parameters.|&nbsp;
`path`|The page's [logical path](g) relative to the content adapter. Do not include a leading slash or file extension.|:heavy_check_mark:
`title`|The page title.|&nbsp;

> [!note]
> While `path` is the only required field, we recommend setting `title` as well.
>
> When setting the `path`, Hugo transforms the given string to a logical path. For example, setting `path` to `A B C` produces a logical path of `/section/a-b-c`.

## Resource map

Construct the map passed to the [`AddResource`](#addresource) method using the fields below.

Key|Description|Required
:--|:--|:-:
`content.mediaType`|The content [media type].|:heavy_check_mark:
`content.value`|The content value as a string or resource.|:heavy_check_mark:
`name`|The resource name.|&nbsp;
`params`|A map of resource parameters.|&nbsp;
`path`|The resources's [logical path](g) relative to the content adapter. Do not include a leading slash.|:heavy_check_mark:
`title`|The resource title.|&nbsp;

> [!note]
> When `content.value` is a string, Hugo generates a new resource with a publication path relative to the page. However, if `content.value` is already a resource, Hugo directly uses its value and publishes it relative to the site root. This latter method is more efficient.
>
> When setting the `path`, Hugo transforms the given string to a logical path. For example, setting `path` to `A B C/cover.jpg` produces a logical path of `/section/a-b-c/cover.jpg`.

## Example

Create pages from remote data, where each page represents a book review.

Step 1
: Create the content structure.

  ```text
  content/
  └── books/
      ├── _content.gotmpl  <-- content adapter
      └── _index.md
  ```

Step 2
: Inspect the remote data to determine how to map key-value pairs to front matter fields.\
  <https://gohugo.io/shared/examples/data/books.json>

Step 3
: Create the content adapter.

  ```go-html-template {file="content/books/_content.gotmpl" copy=true}
  {{/* Get remote data. */}}
  {{ $data := dict }}
  {{ $url := "https://gohugo.io/shared/examples/data/books.json" }}
  {{ with try (resources.GetRemote $url) }}
    {{ with .Err }}
      {{ errorf "Unable to get remote resource %s: %s" $url . }}
    {{ else with .Value }}
      {{ $data = . | transform.Unmarshal }}
    {{ else }}
      {{ errorf "Unable to get remote resource %s" $url }}
    {{ end }}
  {{ end }}

  {{/* Add pages and page resources. */}}
  {{ range $data }}

    {{/* Add page. */}}
    {{ $content := dict "mediaType" "text/markdown" "value" .summary }}
    {{ $dates := dict "date" (time.AsTime .date) }}
    {{ $params := dict "author" .author "isbn" .isbn "rating" .rating "tags" .tags }}
    {{ $page := dict
      "content" $content
      "dates" $dates
      "kind" "page"
      "params" $params
      "path" .title
      "title" .title
    }}
    {{ $.AddPage $page }}

    {{/* Add page resource. */}}
    {{ $item := . }}
    {{ with $url := $item.cover }}
      {{ with try (resources.GetRemote $url) }}
        {{ with .Err }}
          {{ errorf "Unable to get remote resource %s: %s" $url . }}
        {{ else with .Value }}
          {{ $content := dict "mediaType" .MediaType.Type "value" .Content }}
          {{ $params := dict "alt" $item.title }}
          {{ $resource := dict
            "content" $content
            "params" $params
            "path" (printf "%s/cover.%s" $item.title .MediaType.SubType)
          }}
          {{ $.AddResource $resource }}
        {{ else }}
          {{ errorf "Unable to get remote resource %s" $url }}
        {{ end }}
      {{ end }}
    {{ end }}

  {{ end }}
  ```

Step 4
: Create a _page_ template to render each book review.

  ```go-html-template {file="layouts/books/page.html" copy=true}
  {{ define "main" }}
    <h1>{{ .Title }}</h1>

    {{ with .Resources.GetMatch "cover.*" }}
      <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="{{ .Params.alt }}">
    {{ end }}

    <p>Author: {{ .Params.author }}</p>

    <p>
      ISBN: {{ .Params.isbn }}<br>
      Rating: {{ .Params.rating }}<br>
      Review date: {{ .Date | time.Format ":date_long" }}
    </p>

    {{ with .GetTerms "tags" }}
      <p>Tags:</p>
      <ul>
        {{ range . }}
          <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
        {{ end }}
      </ul>
    {{ end }}

    {{ .Content }}
  {{ end }}
  ```

## Multilingual sites

With multilingual sites you can:

1. Create one content adapter for all languages using the [`EnableAllLanguages`](#enablealllanguages) method as described above.
1. Create content adapters unique to each language. See the examples below.

### Translations by file name

With this site configuration:

{{< code-toggle file=hugo >}}
[languages.en]
weight = 1

[languages.de]
weight = 2
{{< /code-toggle >}}

Include a language designator in the content adapter's file name.

```text
content/
└── books/
    ├── _content.de.gotmpl
    ├── _content.en.gotmpl
    ├── _index.de.md
    └── _index.en.md
```

### Translations by content directory

With this site configuration:

{{< code-toggle file=hugo >}}
[languages.en]
contentDir = 'content/en'
weight = 1

[languages.de]
contentDir = 'content/de'
weight = 2
{{< /code-toggle >}}

Create a single content adapter in each directory:

```text
content/
├── de/
│   └── books/
│       ├── _content.gotmpl
│       └── _index.md
└── en/
    └── books/
        ├── _content.gotmpl
        └── _index.md
```

## Page collisions

Two or more pages collide when they have the same publication path. Due to concurrency, the content of the published page is indeterminate. Consider this example:

```text
content/
└── books/
    ├── _content.gotmpl  <-- content adapter
    ├── _index.md
    └── the-hunchback-of-notre-dame.md
```

If the content adapter also creates `books/the-hunchback-of-notre-dame`, the content of the published page is indeterminate. You can not define the processing order.

To detect page collisions, use the `--printPathWarnings` flag when building your site.

[content formats]: /content-management/formats/#classification
[front matter field]: /content-management/front-matter/#fields
[media type]: https://en.wikipedia.org/wiki/Media_type
[syntax]: /templates/introduction/
[template functions]: /functions/
