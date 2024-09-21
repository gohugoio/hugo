---
title: bit.LeadingZeros
description: Counts leading 0 bits in a number (when cast to uint64).
categories: []
keywords: []
action:
  aliases: [clz]
  related:
    - functions/bit/OnesCount
    - functions/bit/TrailingZeros
  returnType: int64
  signatures: [bit.LeadingZeros NUMBER]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ bit.LeadingZeros 0 }} → 64
{{ bit.LeadingZeros -1 }} → 0
{{ bit.LeadingZeros 0x11223344 }} → 32
```
