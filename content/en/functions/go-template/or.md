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

```go-html-template
{{ or 0 1 2 }} → 1
{{ or false "a" 1 }} → a
{{ or 0 true "a" }} → true

{{ or false "" 0 }} → 0
{{ or 0 "" false }} → false
```

{{% include "/_common/functions/go-template/text-template.md" %}}
