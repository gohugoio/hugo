---
title: reflect.IsSite
description: Reports whether the given value is a Site object.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsSite INPUT]
---

{{< new-in 0.154.0 />}}

```go-html-template {file="layouts/page.html"}
{{ with .Site  }}
  {{ reflect.IsSite . }} → true
{{ end }}

{{ with site.GetPage "/examples" }}
  {{ reflect.IsSite . }} → false
{{ end }}
```
