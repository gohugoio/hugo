---
title: plainify
linktitle: plainify
description: Strips any HTML and returns the plain text version of the provided string.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [strings]
signature:
workson: []
hugoversion:
relatedfuncs: [jsonify]
deprecated: false
aliases: []
---

`plainify` strips any HTML and returns the plain text version of the provided string.

```
{{ "<b>BatMan</b>" | plainify }} → "BatMan"
```

See also the [`.PlainWords`, `.Plain`, and `.RawContent` page variables][pagevars].


[pagevars]: /variables/page/


