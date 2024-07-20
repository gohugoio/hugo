---
title: math.Add
description: Adds two or more numbers.
categories: []
keywords: []
action:
  aliases: [add]
  related:
    - functions/math/Div
    - functions/math/Mul
    - functions/math/Product
    - functions/math/Sub
    - functions/math/Sum
  returnType: any
  signatures: [math.Add VALUE VALUE...]
---

If one of the numbers is a [`float`], the result is a `float`.

```go-html-template
{{ add 12 3 2 }} → 17
```

[`float`]: /getting-started/glossary/#float

You can also use the `add` function to concatenate strings.

```go-html-template
{{ add "hu" "go" }} → hugo
```
