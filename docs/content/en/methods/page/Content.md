---
title: Content
description: Returns the rendered content of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Summary
    - methods/page/ContentWithoutSummary
    - methods/page/RawContent
    - methods/page/Plain
    - methods/page/PlainWords
    - methods/page/RenderShortcodes
  returnType: template.HTML
  signatures: [PAGE.Content]
---

The `Content` method on a `Page` object renders Markdown and shortcodes to HTML.

```go-html-template
{{ .Content }}
```
