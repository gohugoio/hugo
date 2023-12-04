---
title: collections.Merge
description: Returns the result of merging two or more maps.
categories: []
keywords: []
action:
  aliases: [merge]
  related:
    - functions/collections/Append
  returnType: any
  signatures: [collections.Merge MAP MAP...]
aliases: [/functions/merge]
---

Returns the result of merging two or more maps from left to right. If a key already exists, `merge` updates its value. If a key is absent, `merge` inserts the value under the new key.

Key handling is case-insensitive.

The following examples use these map definitions:

```go-html-template
{{ $m1 := dict "x" "foo" }}
{{ $m2 := dict "x" "bar" "y" "wibble" }}
{{ $m3 := dict "x" "baz" "y" "wobble" "z" (dict "a" "huey") }}
```

Example 1

```go-html-template
{{ $merged := merge $m1 $m2 $m3 }}

{{ $merged.x }}   → baz
{{ $merged.y }}   → wobble
{{ $merged.z.a }} → huey
```

Example 2

```go-html-template
{{ $merged := merge $m3 $m2 $m1 }}

{{ $merged.x }}   → foo
{{ $merged.y }}   → wibble
{{ $merged.z.a }} → huey
```

Example 3

```go-html-template
{{ $merged := merge $m2 $m3 $m1 }}

{{ $merged.x }}   → foo
{{ $merged.y }}   → wobble
{{ $merged.z.a }} → huey
```

Example 4

```go-html-template
{{ $merged := merge $m1 $m3 $m2 }}

{{ $merged.x }}   → bar
{{ $merged.y }}   → wibble
{{ $merged.z.a }} → huey
```

{{% note %}}
Regardless of depth, merging only applies to maps. For slices, use [append](/functions/collections/append).
{{% /note %}}
