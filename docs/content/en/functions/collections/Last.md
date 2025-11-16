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
{{ range last 10 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```

Set `N` to zero to return an empty collection.

```go-html-template
{{ $emptyPageCollection := last 0 .Pages }}
```

Use `last` and [`where`] together.

[`where`]: /functions/collections/where/

```go-html-template
{{ range where .Pages "Section" "articles" | last 5 }}
  {{ .Render "summary" }}
{{ end }}
```
