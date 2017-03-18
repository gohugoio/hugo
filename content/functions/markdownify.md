---
title: markdownify
linktitle: markdownify
description: Runs the provided string through the Markdown processor.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [markdown,content]
categories: [functions]
ns:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`markdownify` runs the provided string through the Markdown processor. The result will be declared as "safe" so Go/html templates do not filter it.

```
{{ .Title | markdownify }}
```