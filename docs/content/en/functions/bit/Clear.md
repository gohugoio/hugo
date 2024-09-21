---
title: bit.Clear
description: Clears bits of the first number if they're set to 1 in any of the other numbers respectively (aka. ANDN).
categories: []
keywords: []
action:
  aliases: [bandn, bclear]
  related:
    - functions/bit/And
    - functions/bit/Not
    - functions/bit/Or
    - functions/bit/Xnor
    - functions/bit/Xor
  returnType: int64
  signatures: [bit.Clear NUMBER NUMBER...]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#04b" (bit.Clear 0b1100 0b0110) }} → 0b1000
{{ fmt.Printf "%#x" (bit.Clear 0x1C2 0x7F) }} → 0x180
```
