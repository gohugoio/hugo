---
title: lang.FormatPercent
description: Returns a percentage representation of a number with the given precision for the current language and region.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/lang/FormatAccounting
    - functions/lang/FormatCurrency
    - functions/lang/FormatNumber
    - functions/lang/FormatNumberCustom
  returnType: string
  signatures: [lang.FormatPercent PRECISION NUMBER]
---

```go-html-template
{{ 512.5032 | lang.FormatPercent 2 }} → 512.50%
```

{{% include "functions/_common/locales.md" %}}
