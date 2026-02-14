---
title: AlternativeOutputFormats
description: Returns a slice of OutputFormat objects, excluding the current output format, each representing one of the output formats enabled for the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.OutputFormats
    signatures: [PAGE.AlternativeOutputFormats]
---

{{% glossary-term "output format" %}}

The `AlternativeOutputFormats` method on a `Page` object returns a slice of `OutputFormat` objects, excluding the current output format, each representing one of the output formats enabled for the given page. See&nbsp;[details](/configuration/output-formats/).

For example, to generate a `link` element for each of the alternative output formats:

```go-html-template
{{ range .AlternativeOutputFormats }}
  {{ printf "<link rel=%q type=%q href=%q>" .Rel .MediaType.Type .Permalink | safeHTML }}
{{ end }}
```

Hugo renders this to something like:

```html
<link rel="alternate" type="application/rss+xml" href="https://example.org/index.xml">
<link rel="alternate" type="application/json" href="https://example.org/index.json">
```
