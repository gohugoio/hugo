---
title: bit.Not
description: Performs a bitwise NOT on a number.
categories: []
keywords: []
action:
  aliases: [bnot]
  related: []
  returnType: int64
  signatures: [bit.Not NUMBER]
---

{{< new-in 0.135.0 >}}

```go-html-template
{{ fmt.Printf "%#x" (bit.Not 0) }} → -0x1
{{ fmt.Printf "%#x" (bit.Not 0xAA55AA55) }} → -0xaa55aa56
```
