---
title: compare.Conditional
description: Returns one of two arguments depending on the value of the control argument.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [cond]
    returnType: any
    signatures: [compare.Conditional CONTROL ARG1 ARG2]
aliases: [/functions/cond]
---

## Usage

The `compare.Conditional` function returns one of two arguments depending on the value of the control argument. If `CONTROL` is truthy the function returns `ARG1`, otherwise it returns `ARG2`.

Unlike [ternary operators][] in other languages, the `compare.Conditional` function does not perform [short-circuit evaluation][]. It evaluates both `ARG1` and `ARG2` regardless of the `CONTROL` value.

## Examples

```go-html-template
{{ $qty := 42 }}
{{ compare.Conditional (compare.Le $qty 3) "few" "many" }} → many
```

Due to the absence of short-circuit evaluation, these examples throw an error:

```go-html-template
{{ compare.Conditional true "true" (div 1 0) }}
{{ compare.Conditional false (div 1 0) "false" }}
```

[short-circuit evaluation]: https://en.wikipedia.org/wiki/Short-circuit_evaluation
[ternary operators]: https://en.wikipedia.org/wiki/Ternary_conditional_operator
