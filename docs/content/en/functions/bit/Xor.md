---
title: bit.Xor
description: Bitwise XORs two or more numbers together.
categories: []
keywords: []
action:
  aliases: [bxor]
  related:
    - functions/bit/And
    - functions/bit/Clear
    - functions/bit/Or
    - functions/bit/Xnor
  returnType: int64
  signatures: [bit.Xor NUMBER NUMBER...]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#04b" (bit.Xor 0b1100 0b0110) }} → 0b1010
{{ fmt.Printf "%#x" (bit.Xor 0x1C2 0x7F) }} → 0x1bd
```
