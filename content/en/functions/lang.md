---
title: lang
package: lang
description: "TODO.."
date: 2021-07-28
categories: [functions]
keywords: [numbers]
menu:
  docs:
    parent: "functions"
toc: false
signature: ["lang.NumFmt PRECISION NUMBER [OPTIONS [DELIMITER]]"]
aliases: ['/functions/numfmt/']
type: 'template-func'
---

The default options value is `- . ,`.  The default delimiter within the options
value is a space.  If you need to use a space as one of the options, set a
custom delimiter.s

Numbers greater than or equal to 5 are rounded up. For example, if precision is set to `0`, `1.5` becomes `2`, and `1.4` becomes `1`.

```
{{ lang.NumFmt 2 12345.6789 }} → 12,345.68
{{ lang.NumFmt 2 12345.6789 "- , ." }} → 12.345,68
{{ lang.NumFmt 0 -12345.6789 "- . ," }} → -12,346
{{ lang.NumFmt 6 -12345.6789 "- ." }} → -12345.678900
{{ lang.NumFmt 6 -12345.6789 "-|.| " "|" }} → -1 2345.678900
{{ -98765.4321 | lang.NumFmt 2 }} → -98,765.43
```
