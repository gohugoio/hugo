---
title: not
description: Returns the boolean negation of its single argument.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/and
    - functions/go-template/or
  returnType: bool
  signatures: [not VALUE]
---

Unlike the `and` and `or` operators, the `not` operator always returns a boolean value.

```go-html-template
{{ not true }} → false
{{ not false }} → true

{{ not 1 }} → false
{{ not 0 }} → true

{{ not "x" }} → false
{{ not "" }} → true
```

Use the `not` operator, twice in succession, to cast any value to a boolean value. For example:

```go-html-template
{{ 42 | not | not }} → true
{{ "" | not | not }} → false
```

{{% include "functions/go-template/_common/text-template.md" %}}
