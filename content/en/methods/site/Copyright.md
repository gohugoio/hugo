---
title: Copyright
description:  Returns the copyright notice as defined in the site configuration.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [SITE.Copyright]
---

Site configuration:

{{< code-toggle file=hugo >}}
copyright = '© 2023 ABC Widgets, Inc.'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.Copyright }} → © 2023 ABC Widgets, Inc.
```
