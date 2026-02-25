---
title: Rel
description: Returns the rel value of the given output format, either the default or as defined in your project configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [OUTPUTFORMAT.Rel]
---

{{% include "/_common/methods/output-formats/to-use-this-method.md" %}}

```go-html-template
{{ with .Site.Home.OutputFormats.Get "rss" }}
  {{ .Rel }} → alternate
{{ end }}
```
