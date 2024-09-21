---
title: bit.OnesCount
description: Counts the number of 1 bits in a number (aka. population count).
categories: []
keywords: []
action:
  aliases: [popcnt]
  related: []
  returnType: int64
  signatures: [bit.OnesCount NUMBER]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ bit.OnesCount 0 }} → 0
{{ bit.OnesCount -1 }} → 64
{{ bit.OnesCount 0x1111 }} → 4
{{ bit.OnesCount 0xFF }} → 8
```
