---
title: collections.Index
description: Returns the object, element, or value associated with the given key or keys.
categories: []
keywords: []
action:
  aliases: [index]
  related: []
  returnType: any
  signatures:
    - collections.Index COLLECTION KEY...
aliases: [/functions/index,/functions/index-function]
---

Each indexed item must be a map or a slice:

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ index $s 0 }} → a
{{ index $s 1 }} → b

{{ $m := dict "a" 100 "b" 200 }}
{{ index $m "b" }} → 200
```

Use two or more keys to access a nested value:

```go-html-template
{{ $m := dict "a" 100 "b" 200 "c" (slice 10 20 30) }}
{{ index $m "c" 1 }} → 20

{{ $m := dict "a" 100 "b" 200 "c" (dict "d" 10 "e" 20) }}
{{ index $m "c" "e" }} → 20
```

You may also use a slice of keys to access a nested value:

```go-html-template
{{ $m := dict "a" 100 "b" 200 "c" (dict "d" 10 "e" 20) }}
{{ $s := slice "c" "e" }}
{{ index $m $s }} → 20
```

Use the `collections.Index` function to access a nested value when the key is variable. For example, these are equivalent:

```go-html-template
{{ .Site.Params.foo }}

{{ $k := "foo" }}
{{ index .Site.Params $k }}
```
