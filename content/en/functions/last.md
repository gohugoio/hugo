---
title: last
description: "slices an array to only the last <em>N</em>th elements."
keywords: []
categories: [functions]
menu:
  docs:
    parent: functions
toc:
signature: ["last INDEX COLLECTION"]
relatedfuncs: []
---

```go-html-template
{{ range last 10 .Pages }}
  {{ .Render "summary" }}
{{ end }}
```
