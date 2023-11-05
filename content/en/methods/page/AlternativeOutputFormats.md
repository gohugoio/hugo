---
title: AlternativeOutputFormats
description: Returns a slice of OutputFormat objects, excluding the current output format, each representing one of the output formats enabled for the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/OutputFormats
  returnType: page.OutputFormats
  signatures: [PAGE.AlternativeOutputFormats]
---

{{% include "methods/page/_common/output-format-definition.md" %}}

The `AlternativeOutputFormats` method on a `Page` object returns a slice of `OutputFormat` objects, excluding the current output format, each representing one of the output formats enabled for the given page.. See&nbsp;[details](/templates/output-formats/).

## Methods

{{% include "methods/page/_common/output-format-methods.md" %}}

## Example

Generate a `link` element in the `<head>` of each page for each of the alternative output formats:

```go-html-template
<head>
  ...
  {{ $title := printf "%s | %s" .Title site.Title }}
  {{ if .IsHome }}
    {{ $title = site.Title }}
  {{ end }}
  {{ range .AlternativeOutputFormats -}}
    {{ printf `<link rel=%q type=%q href=%q title=%q>` .Rel .MediaType.Type .Permalink $title | safeHTML }}
  {{ end }}
  ...
</head>
```

On the site's home page, Hugo renders this to:

```html
<link rel="alternate" type="application/rss+xml" href="https://example.org/index.xml" title="ABC Widgets, Inc.">
```
