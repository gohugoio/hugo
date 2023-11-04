---
title: lang.FormatNumber
description: Returns a numeric representation of a number with the given precision for the current language.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: [lang.FormatNumber PRECISION NUMBER]
relatedFunctions:
  - lang.FormatAccounting
  - lang.FormatCurrency
  - lang.FormatNumber
  - lang.FormatNumberCustom
  - lang.FormatPercent
---

```go-html-template
{{ 512.5032 | lang.FormatNumber 2 }} → 512.50
```

{{% note %}}
{{% readfile file="/functions/_common/locales.md" %}}
{{% /note %}}
