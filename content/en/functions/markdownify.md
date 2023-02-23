---
title: markdownify
linktitle: markdownify
description: Runs the provided string through the Markdown processor.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2023-02-09
keywords: [markdown,content]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["markdownify INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---


```
{{ .Title | markdownify }}
```

{{< new-in "0.93.0" >}} **Note**: `markdownify` now supports [Render Hooks] just like [`.Page.RenderString`]. However, if you use more complicated [Render Hooks] relying on page context, use [`.Page.RenderString`] instead. See [GitHub issue #9692](https://github.com/gohugoio/hugo/issues/9692) for more details.

[Render Hooks]: /templates/render-hooks/
[`.Page.RenderString`]: /functions/renderstring/
