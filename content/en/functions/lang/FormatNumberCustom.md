---
title: lang.FormatNumberCustom
description: Returns a numeric representation of a number with the given precision using negative, decimal, and grouping options.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: ['lang.FormatNumberCustom PRECISION NUMBER [OPTIONS...]']
relatedFunctions:
  - lang.FormatAccounting
  - lang.FormatCurrency
  - lang.FormatNumber
  - lang.FormatNumberCustom
  - lang.FormatPercent
aliases: ['/functions/numfmt/']
---

This function formats a number with the given precision. The first options parameter is a space-delimited string of characters to represent negativity, the decimal point, and grouping. The default value is `- . ,`. The second options parameter defines an alternate delimiting character.

Note that numbers are rounded up at 5 or greater. So, with precision set to 0, 1.5 becomes 2, and 1.4 becomes 1.

For a simpler function that adapts to the current language, see [`lang.FormatNumber`].


```go-html-template
{{ lang.FormatNumberCustom 2 12345.6789 }} → 12,345.68
{{ lang.FormatNumberCustom 2 12345.6789 "- , ." }} → 12.345,68
{{ lang.FormatNumberCustom 6 -12345.6789 "- ." }} → -12345.678900
{{ lang.FormatNumberCustom 0 -12345.6789 "- . ," }} → -12,346
{{ lang.FormatNumberCustom 0 -12345.6789 "-|.| " "|" }} → -12 346
```

{{% note %}}
{{% readfile file="/functions/_common/locales.md" %}}
{{% /note %}}

[`lang.FormatNumber`]: /functions/lang/formatnumber
