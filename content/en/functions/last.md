---
title: last
description: Slices an array to the last N elements.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: collections
relatedFuncs:
  - collections.After
  - collections.First
  - collections.Last
signature: [last INDEX COLLECTION]
---

```go-html-template
{{ range last 10 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```
