---
title: math.Mul
description: Multiplies two or more numbers.
categories: []
keywords: []
action:
  aliases: [mul]
  related:
    - functions/math/Add
    - functions/math/Div
    - functions/math/Product
    - functions/math/Sub
    - functions/math/Sum
  returnType: any
  signatures: [math.Mul VALUE VALUE...]
---

If one of the numbers is a [`float`], the result is a `float`.

```go-html-template
{{ mul 12 3 2 }} → 72
```

[`float`]: /getting-started/glossary/#float
