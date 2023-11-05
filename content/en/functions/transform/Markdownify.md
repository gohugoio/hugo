---
title: transform.Markdownify
description: Renders markdown to HTML.
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
{{ .Title | markdownify }}
```

If the resulting HTML is a single paragraph, Hugo removes the wrapping `p` tags to produce inline HTML as required per the example above.

To keep the wrapping `p` tags for a single paragraph, use the [`RenderString`] method on the `Page` object, setting the `display` option to `block`.

If the resulting HTML is two or more paragraphs, Hugo leaves the wrapping `p` tags in place.

[`RenderString`]: /methods/page/renderstring/

{{% note %}}
Although the `markdownify` function honors [markdown render hooks] when rendering markdown to HTML, use the `RenderString` method instead of `markdownify` if a render hook accesses `.Page` context. See issue [#9692] for details.

[markdown render hooks]: /templates/render-hooks/
[#9692]: https://github.com/gohugoio/hugo/issues/9692
{{% /note %}}
