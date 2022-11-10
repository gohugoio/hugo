---
title: urlquery
linktitle: urlquery
description: Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query.
date: 2022-01-18
publishdate: 2022-01-18
lastmod: 2022-01-18
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls]
signature: ["urlquery INPUT [INPUT]..."]
hugoversion:
deprecated: false
workson: []
relatedfuncs: []
aliases: []
---


This template code:

```go-html-template
{{ $u := urlquery "https://" "example.com" | safeURL }}
<a href="https://example.org?url={{ $u }}">Link</a>
```

Is rendered to:

```html
<a href="https://example.org?url=https%3A%2F%2Fexample.com">Link</a>
```
