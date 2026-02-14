---
title: OutputFormats
description: Returns a slice of OutputFormat objects, each representing one of the output formats enabled for the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '[]OutputFormat'
    signatures: [PAGE.OutputFormats]
---

{{% glossary-term "output format" %}}

The `OutputFormats` method on a `Page` object returns a slice of `OutputFormat` objects, each representing one of the output formats enabled for the given page. See&nbsp;[details](/configuration/output-formats/).

## Methods

### Canonical

{{< new-in "0.154.4" />}}

(`page.OutputFormat`) Returns the [canonical output format](g) for the current page, if defined. Once you have captured the object, use any of its [associated methods][].

```go-html-template
{{ with .Site.Home.OutputFormats.Canonical }}
  {{ .MediaType.Type }} → text/html
  {{ .MediaType.MainType }} → text
  {{ .MediaType.SubType }} → html
  {{ .Name }} → html
  {{ .Permalink }} → https://example.org/
  {{ .Rel }} → canonical
  {{ .RelPermalink }} → /
{{ end }}
```

### Get

(`page.OutputFormat`) Returns the `OutputFormat` object with the given identifier. Once you have captured the object, use any of its [associated methods][].

```go-html-template
{{ with .Site.Home.OutputFormats.Get "rss" }}
  {{ .MediaType.Type }} → application/rss+xml
  {{ .MediaType.MainType }} → application
  {{ .MediaType.SubType }} → rss
  {{ .Name }} → rss
  {{ .Permalink }} → https://example.org/index.xml
  {{ .Rel }} → alternate
  {{ .RelPermalink }} → /index.xml
{{ end }}
```

## Examples

To render a `link` element pointing to the [canonical output format](g) for the current page:

```go-html-template
{{ with .OutputFormats.Canonical }}
  {{ printf "<link rel=%q type=%q href=%q>" .Rel .MediaType.Type .Permalink | safeHTML }}
{{ end }}
```

To render an anchor element pointing to the `rss` output format for the current page:

```go-html-template
{{ with .OutputFormats.Get "rss" }}
  <a href="{{ .RelPermalink }}">RSS Feed</a>
{{ end }}
```

Please see the [link to output formats] section to understand the importance of the construct above.

[associated methods]: /methods/output-format/
[link to output formats]: /configuration/output-formats/#link-to-output-formats
