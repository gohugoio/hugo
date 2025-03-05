---
title: Content
description: Returns the rendered content of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: template.HTML
    signatures: [PAGE.Content]
---

The `Content` method on a `Page` object renders Markdown and shortcodes to HTML.

```go-html-template
{{ .Content }}
```
