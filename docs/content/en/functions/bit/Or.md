---
title: bit.Or
description: Bitwise ORs two or more numbers together.
categories: []
keywords: []
action:
  aliases: [bor]
  related:
    - functions/bit/And
    - functions/bit/Clear
    - functions/bit/Xnor
    - functions/bit/Xor
  returnType: int64
  signatures: [bit.Or NUMBER NUMBER...]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#04b" (bit.Or 0b1100 0b0110) }} → 0b1110
{{ fmt.Printf "%#x" (bit.Or 0x1C2 0x7F) }} → 0x1ff
```
