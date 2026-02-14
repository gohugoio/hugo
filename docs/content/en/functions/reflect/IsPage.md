---
title: reflect.IsPage
description: Reports whether the given value is a Page object.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsPage INPUT]
---

{{< new-in 0.154.0 />}}

```go-html-template {file="layouts/page.html"}
{{ with site.GetPage "/examples" }}
  {{ reflect.IsPage . }} → true
{{ end }}

{{ with .Site  }}
  {{ reflect.IsPage . }} → false
{{ end }}
```
