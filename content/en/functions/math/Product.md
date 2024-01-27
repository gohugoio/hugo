---
title: math.Product
description: Returns the product of all numbers. Accepts scalars, slices, or both.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/math/Add
    - functions/math/Div
    - functions/math/Mul
    - functions/math/Sub
    - functions/math/Sum
  returnType: float64
  signatures: [math.Product VALUE...]
---

{{< new-in 0.114.0 >}}

```go-html-template
{{ math.Product 1 (slice 2 3) 4 }} → 24
```
