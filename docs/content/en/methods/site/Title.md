---
title: Title
description: Returns the title as defined in the site configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SITE.Title]
---

Site configuration:

{{< code-toggle file=hugo >}}
title = 'My Documentation Site'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Title }} → My Documentation Site
```
