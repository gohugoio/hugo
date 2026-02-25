---
title: Copyright
description: Returns the copyright notice as defined in your project configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SITE.Copyright]
---

Project configuration:

{{< code-toggle file=hugo >}}
copyright = '© 2023 ABC Widgets, Inc.'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Copyright }} → © 2023 ABC Widgets, Inc.
```
