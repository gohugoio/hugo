---
title: math.Div
description: Divides the first number by one or more numbers.
categories: []
keywords: []
action:
  aliases: [div]
  related:
    - functions/math/Add
    - functions/math/Mul
    - functions/math/Product
    - functions/math/Sub
    - functions/math/Sum
  returnType: any
  signatures: [math.Div VALUE VALUE...]
---

If one of the numbers is a [`float`], the result is a `float`.

```go-html-template
{{ div 12 3 2 }} → 2
```

[`float`]: /getting-started/glossary/#float
