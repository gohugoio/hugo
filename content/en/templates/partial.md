---
title: Partial templates
description: Partials are smaller, context-aware components in your list and page templates that can be used economically to keep your templating DRY.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 110
weight: 110
toc: true
aliases: [/templates/partials/,/layout/chrome/]
---

{{< youtube pjS4pOLyB7c >}}

## Use partials in your templates

All partials for your Hugo project are located in a single `layouts/partials` directory. For better organization, you can create multiple subdirectories within `partials` as well:

```txt
layouts/
└── partials/
    ├── footer/
    │   ├── scripts.html
    │   └── site-footer.html
    ├── head/
    │   ├── favicons.html
    │   ├── metadata.html
    │   ├── prerender.html
    │   └── twitter.html
    └── header/
        ├── site-header.html
        └── site-nav.html
```

All partials are called within your templates using the following pattern:

```go-html-template
{{ partial "<PATH>/<PARTIAL>.html" . }}
```

{{% note %}}
One of the most common mistakes with new Hugo users is failing to pass a context to the partial call. In the pattern above, note how "the dot" (`.`) is required as the second argument to give the partial context. You can read more about "the dot" in the [Hugo templating introduction](/templates/introduction/#context).
{{% /note %}}

{{% note %}}
`<PARTIAL>` including `baseof` is reserved. ([#5373](https://github.com/gohugoio/hugo/issues/5373))
{{% /note %}}

As shown in the above example directory structure, you can nest your directories within `partials` for better source organization. You only need to call the nested partial's path relative to the `partials` directory:

```go-html-template
{{ partial "header/site-header.html" . }}
{{ partial "footer/scripts.html" . }}
```

### Variable scoping

The second argument in a partial call is the variable being passed down. The above examples are passing the `.`, which tells the template receiving the partial to apply the current [context][context].

This means the partial will *only* be able to access those variables. The partial is isolated and cannot access the outer scope. From within the partial, `$.Var` is equivalent to `.Var`.

## Returning a value from a partial

In addition to outputting markup, partials can be used to return a value of any type. In order to return a value, a partial must include a lone `return` statement *at the end of the partial*.

### Example GetFeatured

```go-html-template
{{/* layouts/partials/GetFeatured.html */}}
{{ return first . (where site.RegularPages "Params.featured" true) }}
```

```go-html-template
{{/* layouts/index.html */}}
{{ range partial "GetFeatured.html" 5 }}
  [...]
{{ end }}
```

### Example GetImage

```go-html-template
{{/* layouts/partials/GetImage.html */}}
{{ $image := false }}
{{ with .Params.gallery }}
  {{ $image = index . 0 }}
{{ end }}
{{ with .Params.image }}
  {{ $image = . }}
{{ end }}
{{ return $image }}
```

```go-html-template
{{/* layouts/_default/single.html */}}
{{ with partial "GetImage.html" . }}
  [...]
{{ end }}
```

{{% note %}}
Only one `return` statement is allowed per partial file.
{{% /note %}}

## Inline partials

You can also define partials inline in the template. But remember that template namespace is global, so you need to make sure that the names are unique to avoid conflicts.

```go-html-template
Value: {{ partial "my-inline-partial.html" . }}

{{ define "partials/my-inline-partial.html" }}
{{ $value := 32 }}
{{ return $value }}
{{ end }}
```

## Cached partials

The `partialCached` template function provides significant performance gains for complex templates that don't need to be re-rendered on every invocation. See&nbsp;[details][partialcached].

## Examples

### `header.html`

The following `header.html` partial template is used for [spf13.com](https://spf13.com/):

{{< code file=layouts/partials/header.html >}}
<!DOCTYPE html>
<html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
<head>
    <meta charset="utf-8">

    {{ partial "meta.html" . }}

    <base href="{{ .Site.BaseURL }}">
    <title> {{ .Title }} : spf13.com </title>
    <link rel="canonical" href="{{ .Permalink }}">
    {{ if .RSSLink }}<link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

    {{ partial "head_includes.html" . }}
</head>
{{< /code >}}

{{% note %}}
The `header.html` example partial was built before the introduction of block templates to Hugo. Read more on [base templates and blocks](/templates/base/) for defining the outer chrome or shell of your master templates (i.e., your site's head, header, and footer). You can even combine blocks and partials for added flexibility.
{{% /note %}}

### `footer.html`

The following `footer.html` partial template is used for [spf13.com](https://spf13.com/):

{{< code file=layouts/partials/footer.html >}}
<footer>
  <div>
    <p>
    &copy; 2013-14 Steve Francia.
    <a href="https://creativecommons.org/licenses/by/3.0/" title="Creative Commons Attribution">Some rights reserved</a>;
    please attribute properly and link back.
    </p>
  </div>
</footer>
{{< /code >}}

[context]: /templates/introduction/
[customize]: /hugo-modules/theme-components/
[lookup order]: /templates/lookup-order/
[partialcached]: /functions/partials/includecached/
[themes]: /themes/
