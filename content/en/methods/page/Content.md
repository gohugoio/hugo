---
title: Content
description: Returns the rendered content of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/RawContent
    - methods/page/Plain
    - methods/page/PlainWords
    - methods/page/RenderShortcodes
  returnType: template.HTML
  signatures: [PAGE.Content]
---

The `Content` method on a `Page` object renders markdown and shortcodes to HTML. The content does not include front matter.

[shortcodes]: /getting-started/glossary/#shortcode

```go-html-template
{{ .Content }}
```
