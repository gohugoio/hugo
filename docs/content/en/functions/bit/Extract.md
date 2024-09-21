---
title: bit.Extract
description: Extracts a portion of the bits from a number.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/bit/And
    - functions/bit/ShiftRight
  returnType: int64
  signatures: [bit.Extract NUMBER LENGTH SHIFT]
---

{{< new-in 0.135.0 >}}

Extract returns the last LENGTH bits of (NUMBER >> SHIFT).

```go-html-template
{{ fmt.Printf "%#x" (bit.Extract 0xABCDEF 4 8) }} → 0xd
{{ fmt.Printf "%#b" (bit.Extract 0b10011011 3 1) }} → 0b101
```
