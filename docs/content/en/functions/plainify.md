---
title: plainify
description: Returns a string with all HTML tags removed.
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
