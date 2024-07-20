---
title: lang.FormatNumber
description: Returns a numeric representation of a number with the given precision for the current language and region.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/lang/FormatAccounting
    - functions/lang/FormatCurrency
    - functions/lang/FormatNumberCustom
    - functions/lang/FormatPercent
  returnType: string
  signatures: [lang.FormatNumber PRECISION NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatNumber 2 }} → 512.50
```

{{% include "functions/_common/locales.md" %}}
