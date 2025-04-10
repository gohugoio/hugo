---
title: len
description: Returns the length of a string, slice, map, or collection.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: int
    signatures: [len VALUE]
aliases: [/functions/len]
---

With a string:

```go-html-template
{{ "ab" | len }} → 2
{{ "" | len }} → 0
```

With a slice:

```go-html-template
{{ slice "a" "b" | len }} → 2
{{ slice | len }} → 0
```

With a map:

```go-html-template
{{ dict "a" 1 "b" 2  | len }} → 2
{{ dict | len }} → 0
```

With a collection:

```go-html-template
{{ site.RegularPages | len }} → 42
```

You may also determine the number of pages in a collection with:

```go-html-template
{{ site.RegularPages.Len }} → 42
```

{{% include "/_common/functions/go-template/text-template.md" %}}
