---
title: strings.Split
linkTitle: split
description: Returns a slice of strings by splitting STRING by DELIM.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [split]
  returnType: string
  signatures: [strings.Split STRING DELIM]
relatedFunctions:
  - collections.Delimit
  - strings.Split
aliases: [/functions/split]
---

Examples:

```go-html-template
{{ split "tag1,tag2,tag3" "," }} → ["tag1", "tag2", "tag3"]
{{ split "abc" "" }} → ["a", "b", "c"]
```


{{% note %}}
`split` essentially does the opposite of [delimit](/functions/collections/delimit). While `split` creates a slice from a string, `delimit` creates a string from a slice.
{{% /note %}}
