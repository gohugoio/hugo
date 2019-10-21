---
title: jsonify
linktitle: jsonify
description: Encodes a given object to JSON.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings,json]
signature: ["jsonify INPUT"]
workson: []
hugoversion:
relatedfuncs: [plainify]
deprecated: false
aliases: []
---

```
{{ dict "title" .Title "content" .Plain | jsonify }}
```

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/
