---
title: bit.ShiftRight
description: Shifts a number to the right by a number of bits, extending the sign bit (aka. arithmetic shift).
categories: []
keywords: []
action:
  aliases: [asr]
  related:
    - functions/bit/ShiftLeft
  returnType: int64
  signatures: [bit.ShiftRight NUMBER SHIFT]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#x" (bit.ShiftRight 0x10 1) }} → 0x8
{{ fmt.Printf "%#x" (bit.ShiftRight 0x101 4) }} → 0x10
{{/* Sign extension */}}
{{ fmt.Printf "%#x" (bit.ShiftRight -0x1 1) }} → -0x1
{{ fmt.Printf "%#x" (bit.ShiftRight -0x1 20) }} → -0x1
{{ fmt.Printf "%#x" (bit.ShiftRight -0x8 1) }} → -0x4
```
