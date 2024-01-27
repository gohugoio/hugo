---
title: math.Sub
description: Subtracts one or more numbers from the first number. 
categories: []
keywords: []
action:
  aliases: [sub]
  related:
    - functions/math/Add
    - functions/math/Div
    - functions/math/Mul
    - functions/math/Product
    - functions/math/Sum
  returnType: any
  signatures: [math.Sub VALUE VALUE...]
---

If one of the numbers is a [`float`], the result is a `float`.

```go-html-template
{{ sub 12 3 2 }} → 7
```

[`float`]: /getting-started/glossary/#float
