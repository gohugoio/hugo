---
title: split
description: Returns a slice of strings by splitting STRING by DELIM.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["split STRING DELIM"]
relatedfuncs: []
---

Examples:

```go-html-template
{{ split "tag1,tag2,tag3" "," }} → ["tag1", "tag2", "tag3"]
{{ split "abc" "" }} → ["a", "b", "c"]
```


{{% note %}}
`split` essentially does the opposite of [delimit](/functions/delimit). While `split` creates a slice from a string, `delimit` creates a string from a slice.
{{% /note %}}
