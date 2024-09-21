---
title: bit.Xnor
description: Bitwise XNORs two or more numbers together.
categories: []
keywords: []
action:
  aliases: [bxnor]
  related:
    - functions/bit/And
    - functions/bit/Clear
    - functions/bit/Not
    - functions/bit/Or
    - functions/bit/Xor
  returnType: int64
  signatures: [bit.Xnor NUMBER NUMBER...]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#04b" (bit.Xnor 0b1100 0b0110) }} → -0b1011
{{ fmt.Printf "%#x" (bit.Xnor 0x1C2 0x7F) }} → -0x1be
```
