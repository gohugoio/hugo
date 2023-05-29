---
title: .RenderString
description: "Renders markup to HTML."
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [markdown,goldmark,render]
signature: [".RenderString MARKUP"]
---

`.RenderString` is a method on `Page` that renders some markup to HTML using the content renderer defined for that page (if not set in the options).

The method takes an optional map argument with these options:

display ("inline")
: `inline` or `block`. If `inline` (default), surrounding `<p></p>` on short snippets will be trimmed.

markup (defaults to the Page's markup)
: See identifiers in [List of content formats](/content-management/formats/#list-of-content-formats).

Some examples:

```go-html-template
{{ $optBlock := dict "display" "block" }}
{{ $optOrg := dict "markup" "org" }}
{{ "**Bold Markdown**" | $p.RenderString }}
{{ "**Bold Block Markdown**" | $p.RenderString  $optBlock }}
{{ "/italic org mode/" | $p.RenderString  $optOrg }}
```

{{< new-in "0.93.0" >}} **Note**: [markdownify](/functions/markdownify/) uses this function in order to support [Render Hooks](/getting-started/configuration-markup/#markdown-render-hooks).
