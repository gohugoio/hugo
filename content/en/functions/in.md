---
title: in
description: Checks if an element is in an array or slice---or a substring in a string---and returns a boolean.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs:
  - strings.Contains
  - strings.ContainsAny
  - strings.ContainsNonSpace
  - strings.HasPrefix
  - strings.HasSuffix
signature:
  - collections.In SET ITEM
  - in SET ITEM
---

The elements supported are strings, integers and floats, although only float64 will match as expected.

In addition, `in` can also check if a substring exists in a string.

```go-html-template
{{ if in .Params.tags "Git" }}Follow me on GitHub!{{ end }}
```


```go-html-template
{{ if in "this string contains a substring" "substring" }}Substring found!{{ end }}
```
