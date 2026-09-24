---
title: math.Add
description: Returns the result of adding one or more numbers to the first number.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [add]
    returnType: any
    signatures: [math.Add VALUE VALUE...]
---

If one of the numbers is a [`float`](g), the result is a `float`.

```go-html-template
{{ add 12 3 2 }} → 17
```

You can also use the `add` function to concatenate strings.

```go-html-template
{{ add "hu" "go" }} → hugo
```
