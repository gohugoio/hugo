---
title: or
description: Returns the first truthy argument. If all arguments are falsy, returns the last argument.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: any
    signatures: [or VALUE...]
---

{{% include "/_common/functions/truthy-falsy.md" %}}

The `or` function evaluates the arguments from left to right, and returns when the result is determined.

```go-html-template
{{ or 0 1 2 }} → 1 (int)
{{ or false "a" 1 }} → a (string)
{{ or 0 true "a" }} → true (bool)

{{ or false "" 0 }} → 0 (int)
{{ or 0 "" false }} → false (bool)

{{ or true (math.Div 1 0) }} → true (bool)
```

{{% include "/_common/functions/go-template/text-template.md" %}}
