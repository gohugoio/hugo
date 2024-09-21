---
title: bit.TrailingZeros
description: Counts trailing 0 bits in a number (when cast to uint64).
categories: []
keywords: []
action:
  aliases: [ctz]
  related:
    - functions/bit/OnesCount
    - functions/bit/LeadingZeros
  returnType: int64
  signatures: [bit.TrailingZeros NUMBER]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ bit.TrailingZeros 0 }} → 64
{{ bit.TrailingZeros -1 }} → 0
{{ bit.TrailingZeros 0x10000 }} → 16
```
