---
title: lang.FormatAccounting
description: Returns a currency representation of a number for the given currency and precision for the current language and region in accounting notation.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [lang.FormatAccounting PRECISION CURRENCY NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatAccounting 2 "NOK" }} → NOK512.50
```

{{% include "/_common/functions/locales.md" %}}
