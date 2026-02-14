---
title: collections.Apply
description: Returns a slice by transforming each element of the given slice using a specific function and parameters.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [apply]
    returnType: '[]any'
    signatures: [collections.Apply SLICE FUNCTION PARAM...]
aliases: [/functions/apply]
---

The `apply` function takes three or more arguments, depending on the function being applied to the slice elements.

The first argument is the slice itself, the second argument is the function name, and the remaining arguments are passed to the function, with the string `"."` representing the slice element.

```go-html-template
{{ $s := slice "hello" "world" }}

{{ $s = apply $s "strings.FirstUpper" "." }}
{{ $s }} → [Hello World]

{{ $s = apply $s "strings.Replace" "." "l" "_" }}
{{ $s }} →  [He__o Wor_d]
```
