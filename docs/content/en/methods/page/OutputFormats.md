---
title: OutputFormats
description: Returns a slice of OutputFormat objects, each representing one of the output formats enabled for the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/AlternativeOutputFormats
  returnType: '[]OutputFormat'
  signatures: [PAGE.OutputFormats]
toc: true
---

{{% include "methods/page/_common/output-format-definition.md" %}}

The `OutputFormats` method on a `Page` object returns a slice of `OutputFormat` objects, each representing one of the output formats enabled for the given page. See&nbsp;[details](/templates/output-formats/).

## Methods

{{% include "methods/page/_common/output-format-methods.md" %}}

## Example

To link to the RSS feed for the current page:

```go-html-template
{{ with .OutputFormats.Get "rss" -}}
  <a href="{{ .RelPermalink }}">RSS Feed</a>
{{ end }}
```

On the site's home page, Hugo renders this to:

```html
<a href="/index.xml">RSS Feed</a>
```

Please see the [link to output formats] section to understand the importance of the construct above.

[link to output formats]: /templates/output-formats/#link-to-output-formats
