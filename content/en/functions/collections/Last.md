---
title: collections.Last
description: Slices an array to the last N elements.
categories: []
keywords: []
action:
  aliases: [last]
  related:
    - functions/collections/After
    - functions/collections/First
  returnType: any
  signatures: [collections.Last INDEX COLLECTION]
aliases: [/functions/last]
---

```go-html-template
{{ range last 10 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```
