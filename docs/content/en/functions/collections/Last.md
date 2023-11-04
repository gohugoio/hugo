---
title: collections.Last
linkTitle: last
description: Slices an array to the last N elements.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [last]
  returnType: any
  signatures: [collections.Last INDEX COLLECTION]
relatedFunctions:
  - collections.After
  - collections.First
  - collections.Last
aliases: [/functions/last]
---

```go-html-template
{{ range last 10 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```
