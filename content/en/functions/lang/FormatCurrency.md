---
title: lang.FormatCurrency
description: Returns a currency representation of a number for the given currency and precision for the current language and region.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [lang.FormatCurrency PRECISION CURRENCY NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatCurrency 2 "USD" }} → $512.50
```

{{% include "/_common/functions/locales.md" %}}
