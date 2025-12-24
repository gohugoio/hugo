---
title: Copyright
description: Returns the copyright notice as defined in the site configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SITE.Copyright]
---

Site configuration:

{{< code-toggle file=hugo >}}
copyright = '&copy; 2023-{{ now.Year }} ABC Widgets, Inc.'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Copyright }} → © 2023-2025 ABC Widgets, Inc.
```
