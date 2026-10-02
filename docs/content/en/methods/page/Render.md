---
title: Render
description: Returns the result of rendering a view template with the given page as context, or with an optional context argument.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: template.HTML
    signatures: ['PAGE.Render VIEW [CONTEXT]']
aliases: [/functions/render]
---

The `Render` method on a `Page` object renders a [view template][] with the given page as [context](g), or with an optional context argument.

{{< new-in 0.164.0 >}}
The `VIEW` argument now supports slash-separated directory paths.
{{< /new-in >}}

{{< new-in 0.166.0 >}}
This method now accepts an optional `CONTEXT` argument.
{{< /new-in >}}

The `VIEW` argument is the name of a _view_ template, optionally preceded by a slash-separated directory path. Do not include a file extension. Hugo resolves the template via the [template lookup order][], so the same `VIEW` value may map to different templates depending on the page being rendered.

By default, Hugo passes the `Page` object as the context (the dot) when rendering the template. To pass a different context, provide the optional `CONTEXT` argument.

## Examples

The following examples demonstrate calling this method with and without a custom context argument.

### Default context

When called without a context argument, the `Page` object is the context within the template:

```go-html-template {file="layouts/home.html"}
<ul>
  {{ range site.RegularPages }}
    <li>{{ .Render "_views/summary" }}</li>
  {{ end }}
</ul>
```

```go-html-template {file="layouts/_views/summary.html"}
<a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a>
```

### Custom context

To pass additional data to a view template, provide a custom context argument. This example passes a map as context, combining the `Page` object with an additional key-value pair:

```go-html-template {file="layouts/home.html"}
<div>
  {{ range site.RegularPages }}
    {{ .Render "_views/card" (dict "page" . "class" "featured") }}
  {{ end }}
</div>
```

```go-html-template {file="layouts/_views/card.html"}
<div class="card {{ .class }}">
  <h2><a href="{{ .page.RelPermalink }}">{{ .page.LinkTitle }}</a></h2>
  {{ .page.Summary }}
</div>
```

## Organization

As a best practice, place _view_ templates together in a dedicated subdirectory. Hugo does not reserve a directory name for _view_ templates as it does for `_partials`, `_shortcodes`, and `_markup`. The examples below use `_views`, where the underscore prefix differentiates it from other path segments and conveys its purpose, but a directory named `foo` would work equally well.

The following example uses path segments to organize _view_ templates in a dedicated subdirectory:

```tree
layouts/
├── _views/
│   └── summary.html
├── books/
│   └── _views/
│       └── summary.html
├── baseof.html
├── home.html
├── page.html
├── section.html
├── taxonomy.html
└── term.html
```

And this template:

```go-html-template {file="layouts/home.html"}
<ul>
  {{ range site.RegularPages }}
    {{ .Render "_views/summary" }}
  {{ end }}
</ul>
```

When rendering content of type `books`, the `Render` method calls:

```text
layouts/books/_views/summary.html
```

For all other pages, the `Render` method calls:

```text
layouts/_views/summary.html
```

## Notes

Although similar to the [`partial`][] function, there are key differences.

{{% include "/_common/render-vs-partial.md" %}}

[`partial`]: /functions/partials/include/
[template lookup order]: /templates/lookup-order/
[view template]: /templates/types/#view
