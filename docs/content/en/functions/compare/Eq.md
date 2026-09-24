---
title: compare.Eq
description: Reports whether the first argument is equal to any of the subsequent arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [eq]
    returnType: bool
    signatures: ['compare.Eq ARG1 ARG2 [ARG...]']
aliases: [/functions/eq]
---

## Usage

The `compare.Eq` function reports whether the first argument is equal to any of the subsequent arguments. You can also use this function to compare strings, boolean values, dates, and other comparable data types.

## Examples

```go-html-template
{{ compare.Eq 1 1 }} → true
{{ compare.Eq 1 2 }} → false

{{ compare.Eq 1 1 1 }} → true
{{ compare.Eq 1 1 2 }} → true
{{ compare.Eq 1 2 1 }} → true
{{ compare.Eq 1 2 2 }} → false
```

Comparing other data types:

```go-html-template
{{ compare.Eq "ab" "a" }} → false
{{ compare.Eq time.Now (time.AsTime "1964-12-30") }} → false
{{ compare.Eq true false }} → false
```
