---
title: compare.Conditional
description: Returns one of two arguments depending on the value of the control argument.
categories: []
keywords: []
action:
  aliases: [cond]
  related:
    - functions/compare/Default
  returnType: any
  signatures: [compare.Conditional CONTROL ARG1 ARG2]
aliases: [/functions/cond]
---

The CONTROL argument is a boolean value that indicates whether the function should return ARG1 or ARG2. If CONTROL is `true`, the function returns ARG1. Otherwise, the function returns ARG2.

```go-html-template
{{ $qty := 42 }}
{{ cond (le $qty 3) "few" "many" }} → many
```

The CONTROL argument must be either `true` or `false`. To cast a non-boolean value to boolean, pass it through the `not` operator twice.

```go-html-template
{{ cond (42 | not | not) "truthy" "falsy" }} → truthy
{{ cond ("" | not | not) "truthy" "falsy" }} → falsy
```

{{% note %}}
Unlike [ternary operators] in other languages, the `cond` function does not perform [short-circuit evaluation]. The function evaluates both ARG1 and ARG2, regardless of the CONTROL value.

[short-circuit evaluation]: https://en.wikipedia.org/wiki/Short-circuit_evaluation
[ternary operators]: https://en.wikipedia.org/wiki/Ternary_conditional_operator
{{% /note %}}

Due to the absence of short-circuit evaluation, these examples throw an error:

```go-html-template
{{ cond true "true" (div 1 0) }}
{{ cond false (div 1 0) "false" }}
```
