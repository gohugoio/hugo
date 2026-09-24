---
title: compare.Ne
description: Reports whether the first argument is not equal to any of the subsequent arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [ne]
    returnType: bool
    signatures: ['compare.Ne ARG1 ARG2 [ARG...]']
aliases: [/functions/ne]
---

## Usage

The `compare.Ne` function reports whether the first argument is not equal to any of the subsequent arguments. You can also use this function to compare strings, boolean values, dates, and other comparable data types.

## Examples

```go-html-template
{{ compare.Ne 1 1 }} → false
{{ compare.Ne 1 2 }} → true

{{ compare.Ne 1 1 1 }} → false
{{ compare.Ne 1 1 2 }} → false
{{ compare.Ne 1 2 1 }} → false
{{ compare.Ne 1 2 2 }} → true
```

Comparing other data types:

```go-html-template
{{ compare.Ne "ab" "a" }} → true
{{ compare.Ne time.Now (time.AsTime "1964-12-30") }} → true
{{ compare.Ne true false }} → true
```
