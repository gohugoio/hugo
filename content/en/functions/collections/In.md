---
title: collections.In
description: Reports whether an element is in an array or slice, or if a substring is in a string.
categories: []
keywords: []
action:
  aliases: [in]
  related:
    - functions/collections/Slice
    - functions/strings/Contains
    - functions/strings/ContainsAny
    - functions/strings/ContainsNonSpace
    - functions/strings/HasPrefix
    - functions/strings/HasSuffix
  returnType: bool
  signatures: [collections.In SET ITEM]
aliases: [/functions/in]
---

```go-html-template
{{ $s := slice "a" "b" "c" }}
{{ in $s "b" }} → true
```

```go-html-template
{{ $s := slice 1 2 3 }}
{{ in $s 2 }} → true
```

```go-html-template
{{ $s := slice 1.11 2.22 3.33 }}
{{ in $s 2.22 }} → true
```

```go-html-template
{{ $s := "abc" }}
{{ in $s "b" }} → true
```
