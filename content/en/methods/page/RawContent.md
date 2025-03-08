---
title: RawContent
description: Returns the raw content of the given page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [PAGE.RawContent]
---

The `RawContent` method on a `Page` object returns the raw content. The raw content does not include front matter.

```go-html-template
{{ .RawContent }}
```

This is useful when rendering a page in a plain text [output format](g).

> [!note]
> [Shortcodes](g) within the content are not rendered. To get the raw content with shortcodes rendered, use the [`RenderShortcodes`] method on a `Page` object.

[`RenderShortcodes`]: /methods/page/rendershortcodes/
