---
title: compare.Default
description: Returns the second argument if set, else the first argument.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [default]
    returnType: any
    signatures: [compare.Default DEFAULT INPUT]
aliases: [/functions/default]
---

## Usage

The `compare.Default` function returns the second argument if set, else the first argument.

> [!NOTE]
> When the second argument is the boolean `false` value, the `compare.Default` function returns `false`. All _other_ falsy values are considered unset.
>
> The falsy values are `false`, `0`, any `nil` pointer or interface value, any array, slice, map, or string of length zero, and zero `time.Time` values.
>
> Everything else is truthy.
>
> To set a default value based on truthiness, use the [`or`][] operator instead.

## Examples

When the second argument is set:

```go-html-template
{{ 1             | compare.Default 42 }} → 1
{{ "foo"         | compare.Default 42 }} → foo
{{ dict "k" "v"  | compare.Default 42 }} → map[k:v]
{{ slice "a" "b" | compare.Default 42 }} → [a b]
{{ true          | compare.Default 42 }} → true

<!-- As noted above, the boolean "false" is considered set -->
{{ false         | compare.Default 42 }} → false
```

When the second argument is not set:

```go-html-template
{{ 0     | compare.Default 42 }} → 42
{{ ""    | compare.Default 42 }} → 42
{{ dict  | compare.Default 42 }} → 42
{{ slice | compare.Default 42 }} → 42

```

[`or`]: /functions/go-template/or/
