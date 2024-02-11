---
title: transform.Markdownify
description: Renders Markdown to HTML.
categories: []
keywords: []
action:
  aliases: [markdownify]
  related:
    - methods/page/RenderString
    - methods/page/RenderShortcodes
  returnType: template.HTML
  signatures: [transform.Markdownify INPUT]
aliases: [/functions/markdownify]
---

```go-html-template
<h2>{{ .Title | markdownify }}</h2>
```

If the resulting HTML is a single paragraph, Hugo removes the wrapping `p` tags to produce inline HTML as required per the example above.

To keep the wrapping `p` tags for a single paragraph, use the [`RenderString`] method on the `Page` object, setting the `display` option to `block`.

[`RenderString`]: /methods/page/renderstring/

{{% note %}}
Although the `markdownify` function honors [Markdown render hooks] when rendering Markdown to HTML, use the `RenderString` method instead of `markdownify` if a render hook accesses `.Page` context. See issue [#9692] for details.

[Markdown render hooks]: /render-hooks/
[#9692]: https://github.com/gohugoio/hugo/issues/9692
{{% /note %}}
