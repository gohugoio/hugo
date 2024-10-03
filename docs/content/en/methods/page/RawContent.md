---
title: RawContent
description: Returns the raw content of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/Content
    - methods/page/Summary
    - methods/page/ContentWithoutSummary
    - methods/page/Plain
    - methods/page/PlainWords
    - methods/page/RenderShortcodes
  returnType: string
  signatures: [PAGE.RawContent]
---

The `RawContent` method on a `Page` object returns the raw content. The raw content does not include front matter.

```go-html-template
{{ .RawContent }}
```

This is useful when rendering a page in a plain text [output format].

{{% note %}}
[Shortcodes] within the content are not rendered. To get the raw content with shortcodes rendered, use the [`RenderShortcodes`] method on a `Page` object.

[shortcodes]: /getting-started/glossary/#shortcode
[`RenderShortcodes`]: /methods/page/rendershortcodes/
{{% /note %}}

[output format]: /templates/output-formats/
