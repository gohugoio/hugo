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

{{< new-in "0.93.0" >}} **Note**: `markdownify` now supports [Render Hooks] just like [.RenderString](/functions/renderstring/).

[Render Hooks]: /templates/render-hooks/
