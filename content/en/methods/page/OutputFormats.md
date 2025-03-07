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

{{% include "/_common/methods/page/output-format-methods.md" %}}

## Example

To link to the RSS feed for the current page:

```go-html-template
{{ with .OutputFormats.Get "rss" }}
  <a href="{{ .RelPermalink }}">RSS Feed</a>
{{ end }}
```

On the site's home page, Hugo renders this to:

```html
<a href="/index.xml">RSS Feed</a>
```

Please see the [link to output formats] section to understand the importance of the construct above.

[link to output formats]: /configuration/output-formats/#link-to-output-formats
