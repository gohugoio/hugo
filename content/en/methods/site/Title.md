---
title: Title
description: Returns the title as defined in your project configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SITE.Title]
---

Project configuration:

{{< code-toggle file=hugo >}}
title = 'My Documentation Site'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Title }} → My Documentation Site
```
