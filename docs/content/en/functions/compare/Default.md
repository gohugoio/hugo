---
title: compare.Default
description: Returns the second argument if set, else the first argument.
keywords: []
params:
  functions_and_methods:
    aliases: [default]
    returnType: any
    signatures: [compare.Default DEFAULT INPUT]
aliases: [/functions/default]
---

The `default` function returns the second argument if set, else the first argument.

> [!note]
> When the second argument is the boolean `false` value, the `default` function returns `false`. All _other_ falsy values are considered unset.
>
> The falsy values are `false`, `0`, any `nil` pointer or interface value, any array, slice, map, or string of length zero, and zero `time.Time` values.
>
> Everything else is truthy.
>
> To set a default value based on truthiness, use the [`or`] operator instead.

The `default` function returns the second argument if set:

```go-html-template
{{ default 42 1 }} → 1
{{ default 42 "foo" }} → foo
{{ default 42 (dict "k" "v") }} → map[k:v]
{{ default 42 (slice "a" "b") }} → [a b]
{{ default 42 true }} → true

<!-- As noted above, the boolean "false" is considered set -->
{{ default 42 false }} → false
```

The `default` function returns the first argument if the second argument is not set:

```go-html-template
{{ default 42 0 }} → 42
{{ default 42 "" }} → 42
{{ default 42 dict }} → 42
{{ default 42 slice }} → 42
{{ default 42 <nil> }} → 42
```

[`or`]: /functions/go-template/or/
