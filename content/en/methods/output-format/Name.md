---
title: Name
description: Returns the identifier of the given output format.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [OUTPUTFORMAT.Name]
---

{{% include "/_common/methods/output-formats/to-use-this-method.md" %}}

```go-html-template
{{ with .Site.Home.OutputFormats.Get "rss" }}
  {{ .Name }} → rss
{{ end }}
```
