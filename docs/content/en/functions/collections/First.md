---
title: collections.First
description: Returns the given collection, limited to the first N elements.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [first]
    returnType: any
    signatures: [collections.First N COLLECTION]
aliases: [/functions/first]
---

```go-html-template
{{ slice "a" "b" "c" | first 1 }} → [a]
{{ slice "a" "b" "c" | first 2 }} → [a b]
```

Given that a string is in effect a read-only slice of bytes, this function can be used to return the specified number of bytes from the beginning of the string:

```go-html-template
{{ "abc" | first 1 }} → a
{{ "abc" | first 2 }} → ab
```

Note that a _character_ may consist of multiple _bytes_:

```go-html-template
{{ "Schön" | first 3 }} → Sch
{{ "Schön" | first 4 }} → Sch\xc3
{{ "Schön" | first 5 }} → Schö
```

To use the `collections.First` function with a page collection:

```go-html-template
{{ range first 5 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```

Set `N` to zero to return an empty collection:

```go-html-template
{{ $emptyPageCollection := first 0 .Pages }}
```

Use `first` and [`where`][] together:

```go-html-template
{{ range where .Pages "Section" "articles" | first 5 }}
  {{ .Render "summary" }}
{{ end }}
```

[`where`]: /functions/collections/where/
