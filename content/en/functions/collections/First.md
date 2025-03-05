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
{{ range first 5 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```

Set `N` to zero to return an empty collection.

```go-html-template
{{ $emptyPageCollection := first 0 .Pages }}
```

Use `first` and [`where`] together.

```go-html-template
{{ range where .Pages "Section" "articles" | first 5 }}
  {{ .Render "summary" }}
{{ end }}
```

[`where`]: /functions/collections/where/
