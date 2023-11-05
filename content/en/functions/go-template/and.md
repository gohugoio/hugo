---
title: and
description: Returns the first falsy argument. If all arguments are truthy, returns the last argument.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/not
    - functions/go-template/or
  returnType: any
  signatures: [and VALUE...]
---

{{% include "functions/go-template/_common/truthy-falsy.md" %}}

```go-html-template
{{ and 1 0 "" }} → 0 (int)
{{ and 1 false 0 }} → false (bool)

{{ and 1 2 3 }} → 3 (int)
{{ and "a" "b" "c" }} → c (string)
{{ and "a" 1 true }} → true (bool)
```

{{% include "functions/go-template/_common/text-template.md" %}}
