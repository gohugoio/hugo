---
title: Rotate
description: Returns a collection of pages that vary along the specified dimension while sharing the current page's values for the other dimensions, including the current page, sorted by the dimension's default sort order.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pages
    signatures: [PAGE.Rotate DIMENSION]
---

{{< new-in 0.153.0 />}}

The rotate method on a page object returns a collection of pages that vary along the specified [dimension](g), while holding the other dimensions constant. The result includes the current page and is sorted according to the rules of the specified dimension. For example, rotating along [language](g) returns all language variants that share the current page's [version](g) and [role](g).

The `DIMENSION` argument must be one of `language`, `version`, or `role`.

## Sort order

Use the following rules to understand how Hugo sorts the collection returned by the `Rotate` method.

| Dimension | Primary Sort | Secondary Sort |
| :--- | :--- | :--- |
| Language | Weight ascending | Lexicographical ascending |
| Version | Weight ascending | Semantic version descending |
| Role | Weight ascending | Lexicographical ascending |

## Examples

To render a list of the current page's language variants, including the current page, while sharing its current version and role:

```go-html-template
{{/* Returns languages sorted by weight ascending, then lexicographically ascending */}}
{{ range .Rotate "language" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To render a list of the current page's version variants, including the current page, while sharing its current language and role:

```go-html-template
{{/* Returns versions sorted by weight ascending, then semantic version descending */}}
{{ range .Rotate "version" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

To render a list of the current page's role variants, including the current page, while sharing its current language and version:

```go-html-template
{{/* Returns roles sorted by weight ascending, then lexicographically ascending */}}
{{ range .Rotate "role" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```
