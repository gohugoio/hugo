---
title: or
description: Returns the first truthy argument. If all arguments are falsy, returns the last argument.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/and
    - functions/go-template/not
  returnType: any
  signatures: [or VALUE...]
---

{{% include "functions/go-template/_common/truthy-falsy.md" %}}

```go-html-template
{{ or 0 1 2 }} → 1
{{ or false "a" 1 }} → a
{{ or 0 true "a" }} → true

{{ or false "" 0 }} → 0
{{ or 0 "" false }} → false
```

{{% include "functions/go-template/_common/text-template.md" %}}
