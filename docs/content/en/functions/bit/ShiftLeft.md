---
title: bit.ShiftLeft
description: Shifts a number to the left by a number of bits.
categories: []
keywords: []
action:
  aliases: [lsl]
  related:
    - functions/bit/ShiftRight
  returnType: int64
  signatures: [bit.ShiftLeft NUMBER SHIFT]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#x" (bit.ShiftLeft 0x10 1) }} → 0x20
{{ fmt.Printf "%#x" (bit.ShiftLeft 0x101 4) }} → 0x1010
{{ fmt.Printf "%#x" (bit.ShiftLeft 0x1 63) }} → -0x8000000000000000 
{{ fmt.Printf "%#x" (bit.ShiftLeft 0x1 64) }} → 0x0
```
