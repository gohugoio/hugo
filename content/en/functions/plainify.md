---
title: plainify
description: Returns a string with all HTML tags removed.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: transform
relatedFuncs: []
signature:
  - transform.Plainify
  - plainify INPUT
---

```go-html-template
{{ "<b>BatMan</b>" | plainify }} → "BatMan"
```

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/
