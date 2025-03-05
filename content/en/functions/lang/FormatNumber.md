---
title: lang.FormatNumber
description: Returns a numeric representation of a number with the given precision for the current language and region.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [lang.FormatNumber PRECISION NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatNumber 2 }} → 512.50
```

{{% include "/_common/functions/locales.md" %}}
