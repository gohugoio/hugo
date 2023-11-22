---
title: math.Sum
description: Returns the sum of all numbers. Accepts scalars, slices, or both.
categories: []
action:
  aliases: []
  related:
    - functions/math/Add
    - functions/math/Div
    - functions/math/Mul
    - functions/math/Product
    - functions/math/Sub
  returnType: float64
  signatures: [math.Sum VALUE...]
---

{{< new-in 0.114.0 >}}

```go-html-template
{{ math.Sum 1 (slice 2 3) 4 }} → 10
```
