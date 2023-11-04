---
title: transform.Plainify
linkTitle: plainify
description: Returns a string with all HTML tags removed.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [plainify]
  returnType: string
  signatures: [transform.Plainify INPUT]
relatedFunctions: []
aliases: [/functions/plainify]
---

```go-html-template
{{ "<b>BatMan</b>" | plainify }} → "BatMan"
```

See also the `.PlainWords`, `.Plain`, and `.RawContent` [page variables][pagevars].

[pagevars]: /variables/page/
