---
title: compare.Le
description: Reports whether the first argument is less than or equal to all of the subsequent arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [le]
    returnType: bool
    signatures: ['compare.Le ARG1 ARG2 [ARG...]']
aliases: [/functions/le]
---

## Usage

The `compare.Le` function reports whether the first argument is less than or equal to all of the subsequent arguments. Numbers are compared by value, regardless of type. You can also use this function to compare strings, boolean values, dates, and other comparable data types.

## Examples

```go-html-template
{{ compare.Le 1 1 }} → true
{{ compare.Le 1 2 }} → true
{{ compare.Le 2 1 }} → false

{{ compare.Le 1 1 1 }} → true
{{ compare.Le 1 1 2 }} → true
{{ compare.Le 1 2 1 }} → true
{{ compare.Le 1 2 2 }} → true

{{ compare.Le 2 1 1 }} → false
{{ compare.Le 2 1 2 }} → false
{{ compare.Le 2 2 1 }} → false
```

Comparing numbers of different types:

```go-html-template
{{ compare.Le 1 1.0 }} → true
```

Comparing other data types:

```go-html-template
{{ compare.Le "ab" "a" }} → false
{{ compare.Le time.Now (time.AsTime "1964-12-30") }} → false
{{ compare.Le true false }} → false
```
