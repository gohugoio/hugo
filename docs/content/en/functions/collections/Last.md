---
title: collections.Last
description: Returns the given collection, limited to the last N elements.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [last]
    returnType: any
    signatures: [collections.Last N COLLECTION]
aliases: [/functions/last]
---

```go-html-template
{{ slice "a" "b" "c" | last 1 }} → [c]
{{ slice "a" "b" "c" | last 2 }} → [b c]
```

Given that a string is in effect a read-only slice of bytes, this function can be used to return the specified number of bytes from the end of the string:

```go-html-template
{{ "abc" | last 1 }} → c
{{ "abc" | last 2 }} → bc
```

Note that a _character_ may consist of multiple _bytes_:

```go-html-template
{{ "Schön" | last 1 }} → n
{{ "Schön" | last 2 }} → \xb6n
{{ "Schön" | last 3 }} → ön
```

To use the `collections.Last` function with a page collection:

```go-html-template
{{ range last 5 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```

Set `N` to zero to return an empty collection:

```go-html-template
{{ $emptyPageCollection := last 0 .Pages }}
```

Use `last` and [`where`][] together:

[`where`]: /functions/collections/where/

```go-html-template
{{ range where .Pages "Section" "articles" | last 5 }}
  {{ .Render "summary" }}
{{ end }}
```
