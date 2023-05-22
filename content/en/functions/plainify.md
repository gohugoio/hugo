---
title: plainify
description: Strips any HTML and returns the plain text version of the provided string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["plainify INPUT"]
relatedfuncs: [jsonify]
---

```go-html-template
{{ "<b>BatMan</b>" | plainify }} → "BatMan"
```

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/
