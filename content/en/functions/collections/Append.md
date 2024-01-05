---
title: collections.Append
description: Appends one or more elements to a slice and returns the resulting slice.
categories: []
keywords: []
action:
  aliases: [append]
  related:
    - functions/collections/Merge
  returnType: any
  signatures:
    - collections.Append ELEMENT [ELEMENT...] COLLECTION
    - collections.Append COLLECTION1 COLLECTION2
aliases: [/functions/append]
---

This function appends all elements, excluding the last, to the last element. This allows [pipe](/getting-started/glossary/#pipeline) constructs as shown below.

Append a single element to a slice:

```go-html-template
{{ $s := slice "a" "b" }}
{{ $s }} → [a b]

{{ $s = $s | append "c" }}
{{ $s }} → [a b c]
```

Append two elements to a slice:

```go-html-template
{{ $s := slice "a" "b" }}
{{ $s }} → [a b]

{{ $s = $s | append "c" "d" }}
{{ $s }} → [a b c d]
```

Append two elements, as a slice, to a slice. This produces the same result as the previous example:

```go-html-template
{{ $s := slice "a" "b" }}
{{ $s }} → [a b]

{{ $s = $s | append (slice "c" "d") }}
{{ $s }} → [a b c d]
```

Start with an empty slice:

```go-html-template
{{ $s := slice }}
{{ $s }} → []

{{ $s = $s | append "a" }}
{{ $s }} → [a]

{{ $s = $s | append "b" "c" }}
{{ $s }} → [a b c]

{{ $s = $s | append (slice "d" "e") }}
{{ $s }} → [a b c d e]
```

If you start with a slice of a slice:

```go-html-template
{{ $s := slice (slice "a" "b") }}
{{ $s }} → [[a b]]

{{ $s = $s | append (slice "c" "d") }}
{{ $s }} → [[a b] [c d]]
```

To create a slice of slices, starting with an empty slice:

```go-html-template
{{ $s := slice }}
{{ $s }} → []

{{ $s = $s | append (slice (slice "a" "b")) }}
{{ $s }} → [[a b]]

{{ $s = $s | append (slice "c" "d") }}
{{ $s }} → [[a b] [c d]]
```

Although the elements in the examples above are strings, you can use the `append` function with any data type, including Pages. For example, on the home page of a corporate site, to display links to the two most recent press releases followed by links to the four most recent articles:

```go-html-template
{{ $p := where site.RegularPages "Type" "press-releases" | first 2 }}
{{ $p = $p | append (where site.RegularPages "Type" "articles" | first 4) }}

{{ with $p }}
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```
