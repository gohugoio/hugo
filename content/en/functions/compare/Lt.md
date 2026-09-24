---
title: compare.Lt
description: Reports whether the first argument is less than all of the subsequent arguments.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [lt]
    returnType: bool
    signatures: ['compare.Lt ARG1 ARG2 [ARG...]']
aliases: [/functions/lt]
---

## Usage

The `compare.Lt` function reports whether the first argument is less than all of the subsequent arguments. Numbers are compared by value, regardless of type. You can also use this function to compare strings, boolean values, dates, and other comparable data types.

## Examples

```go-html-template
{{ compare.Lt 1 1 }} → false
{{ compare.Lt 1 2 }} → true
{{ compare.Lt 2 1 }} → false

{{ compare.Lt 1 1 1 }} → false
{{ compare.Lt 1 1 2 }} → false
{{ compare.Lt 1 2 1 }} → false
{{ compare.Lt 1 2 2 }} → true

{{ compare.Lt 2 1 1 }} → false
{{ compare.Lt 2 1 2 }} → false
{{ compare.Lt 2 2 1 }} → false
```

Comparing numbers of different types:

```go-html-template
{{ compare.Lt 1 1.0 }} → false
```

Comparing other data types:

```go-html-template
{{ compare.Lt "ab" "a" }} → false
{{ compare.Lt time.Now (time.AsTime "1964-12-30") }} → false
{{ compare.Lt true false }} → false
```
