---
title: .RenderString
description: "Renders markup to HTML."
godocref:
date: 2019-12-18
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [markdown,goldmark,render]
signature: [".RenderString MARKUP"]
---

{{< new-in "0.62.0" >}} 

`.RenderString` is a method on `Page` that renders some markup to HTML using the content renderer defined for that page (if not set in the options).

*Note* that this method does not parse and render shortcodes.

The method takes an optional map argument with these options:

display ("inline")
: `inline` or `block`. If `inline` (default), surrounding ´<p></p>` on short snippets will be trimmed.

markup (defaults to the Page's markup)
: See identifiers in [List of content formats](/content-management/formats/#list-of-content-formats).

Some examples:

```go-html-template
{{ $optBlock := dict "display" "block" }}
{{ $optOrg := dict "markup" "org" }}
{{ "**Bold Markdown**" | $p.RenderString }}
{{  "**Bold Block Markdown**" | $p.RenderString  $optBlock }}
{{  "/italic org mode/" | $p.RenderString  $optOrg }}
```


**Note** that this method is more powerful than the similar [markdownify](/functions/markdownify/) function as it also supports [Render Hooks](/getting-started/configuration-markup/#markdown-render-hooks) and it has options to render other markup formats.
