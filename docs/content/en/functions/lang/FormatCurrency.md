---
title: lang.FormatCurrency
description: Returns a currency representation of a number for the given currency and precision for the current language.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [lang.FormatAccounting PRECISION CURRENCY NUMBER]
relatedFunctions:
  - lang.FormatAccounting
  - lang.FormatCurrency
  - lang.FormatNumber
  - lang.FormatNumberCustom
  - lang.FormatPercent
---

```go-html-template
{{ 512.5032 | lang.FormatCurrency 2 "USD" }} → $512.50
```

{{% note %}}
{{% readfile file="/functions/_common/locales.md" %}}
{{% /note %}}
