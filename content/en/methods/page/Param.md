---
title: Param
description: Returns a page parameter with the given key, falling back to a site parameter if present.
categories: []
keywords: []
action:
  related: []
  returnType: any
  signatures: [PAGE.Param KEY]
aliases: [/functions/param]
---

The `Param` method on a `Page` object looks for the given `KEY` in page parameters, and returns the corresponding value. If it cannot find the `KEY` in page parameters, it looks for the `KEY` in site parameters. If it cannot find the `KEY` in either location, the `Param` method returns `nil`.

Site and theme developers commonly set parameters at the site level, allowing content authors to override those parameters at the page level.

For example, to show a table of contents on every page, but allow authors to hide the table of contents as needed:

Configuration:

{{< code-toggle file=hugo >}}
[params]
display_toc = true
{{< /code-toggle >}}

Content:

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
date = 2023-01-01
draft = false
[params]
display_toc = false
{{< /code-toggle >}}

Template:

```go-html-template
{{ if .Param "display_toc" }}
  {{ .TableOfContents }}
{{ end }}
```

The `Param` method returns the value associated with the given `KEY`, regardless of whether the value is truthy or falsy. If you need to ignore falsy values, use this construct instead:

```go-html-template
{{ or .Params.foo site.Params.foo }}
```
