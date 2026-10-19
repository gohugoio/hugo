---
title: transform.Plainify
description: Returns a string with all HTML tags removed.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [plainify]
    returnType: any
    signatures: [transform.Plainify INPUT]
---

If the argument is of type `template.HTML`, returns `template.HTML`, else returns a `string`.

```go-html-template
{{ "<b>BatMan</b>" | plainify }} → BatMan
```
