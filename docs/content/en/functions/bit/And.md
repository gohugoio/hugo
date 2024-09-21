---
title: bit.And
description: Bitwise ANDs two or more numbers together.
categories: []
keywords: []
action:
  aliases: [band]
  related:
    - functions/bit/Clear
    - functions/bit/Or
    - functions/bit/Xnor
    - functions/bit/Xor
  returnType: int64
  signatures: [bit.And NUMBER NUMBER...]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#04b" (bit.And 0b1100 0b0110) }} → 0b0100
{{ fmt.Printf "%#x" (bit.And 0x1C2 0x7F) }} → 0x42
```
