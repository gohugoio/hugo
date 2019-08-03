---
title: markdownify
linktitle: markdownify
description: Runs the provided string through the Markdown processor. Takes an optional options map specifying the markup format.
godocref:
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
{{ .Title | markdownify (dict "format" "org") }}
{{ .Title | markdownify (dict "format" "asciidoc") }}
```
