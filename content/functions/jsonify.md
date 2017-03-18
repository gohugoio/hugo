---
title: jsonify
linktitle: jsonify
description: Encodes a given object to JSON.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [strings,json]
ns:
signature:
workson: []
hugoversion:
relatedfuncs: [plainify]
deprecated: false
aliases: []
---

`jsonify` encodes a given object to JSON and converts it to HTML-safe content.

```
{{ dict "title" .Title "content" .Plain | jsonify }}
```

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/