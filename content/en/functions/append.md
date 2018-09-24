---
title: append
description: "`append` appends one or more values to a slice and returns the resulting slice."
date: 2018-09-14
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [collections]
signature: ["COLLECTION | append VALUE [VALUE]...", "COLLECTION | append COLLECTION"]
workson: []
hugoversion: "0.49"
relatedfuncs: [last,first,where,slice]
aliases: []
---

An example appending single values:

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ $s = $s | append "d" "e" }}
{{/* $s now contains a []string with elements "a", "b", "c", "d", and "e" */}}

```

The same example appending a slice to a slice:


```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ $s = $s | append (slice "d" "e") }}
```

The `append` function works for all types, including `Pages`.




