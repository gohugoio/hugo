---
title: markdownify
linktitle: markdownify
description: Runs the provided string through the Markdown processor.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
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

*Note*: if you need [Render Hooks][], which `markdownify` doesn't currently
support, use [.RenderString](/functions/renderstring/) instead.

[Render Hooks]: /getting-started/configuration-markup/#markdown-render-hooks
