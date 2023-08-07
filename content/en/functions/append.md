---
title: append
description: "`append` appends one or more values to a slice and returns the resulting slice."
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [collections]
signature: ["COLLECTION | append VALUE [VALUE]...", "COLLECTION | append COLLECTION"]
relatedfuncs: [last,first,where,slice]
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

If a slice contains other slices, further slices will be appended as values:

```go-html-template
{{ $s := slice (slice "a" "b") (slice "c" "d") }}
{{ $s = $s | append (slice "e" "f") (slice "g" "h") }}
{{/* $s now contains a [][]string containing four slices:  ["a" "b"],  ["c" "d"], ["e" "f"], and ["g" "h"] */}}
```

The `append` function works for all types, including `Pages`.
