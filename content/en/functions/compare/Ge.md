---
title: compare.Ge
description: Reports whether the first argument is greater than or equal to all of the subsequent arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [ge]
    returnType: bool
    signatures: ['compare.Ge ARG1 ARG2 [ARG...]']
aliases: [/functions/ge]
---

## Usage

The `compare.Ge` function reports whether the first argument is greater than or equal to all of the subsequent arguments. Numbers are compared by value, regardless of type. You can also use this function to compare strings, boolean values, dates, and other comparable data types.

## Examples

```go-html-template
{{ compare.Ge 1 1 }} → true
{{ compare.Ge 1 2 }} → false
{{ compare.Ge 2 1 }} → true

{{ compare.Ge 1 1 1 }} → true
{{ compare.Ge 1 1 2 }} → false
{{ compare.Ge 1 2 1 }} → false
{{ compare.Ge 1 2 2 }} → false

{{ compare.Ge 2 1 1 }} → true
{{ compare.Ge 2 1 2 }} → true
{{ compare.Ge 2 2 1 }} → true
```

Comparing numbers of different types:

```go-html-template
{{ compare.Ge 1 1.0 }} → true
```

Comparing other data types:

```go-html-template
{{ compare.Ge "ab" "a" }} → true
{{ compare.Ge time.Now (time.AsTime "1964-12-30") }} → true
{{ compare.Ge true false }} → true
```
