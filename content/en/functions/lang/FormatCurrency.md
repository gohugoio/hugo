---
title: lang.FormatCurrency
description: Returns a currency representation of a number for the given currency and precision for the current language and region.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/lang/FormatAccounting
    - functions/lang/FormatNumber
    - functions/lang/FormatNumberCustom
    - functions/lang/FormatPercent
  returnType: string
  signatures: [lang.FormatCurrency PRECISION CURRENCY NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatCurrency 2 "USD" }} → $512.50
```

{{% include "functions/_common/locales.md" %}}
