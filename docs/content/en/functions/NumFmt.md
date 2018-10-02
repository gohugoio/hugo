---
title: lang.NumFmt
description: "Formats a number with a given precision using the requested `negative`, `decimal`, and `grouping` options. The `options` parameter is a string consisting of `<negative> <decimal> <grouping>`."
godocref: ""
workson: []
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-08-21
categories: [functions]
keywords: [numbers]
menu:
  docs:
    parent: "functions"
toc: false
signature: ["lang.NumFmt PRECISION NUMBER [OPTIONS [DELIMITER]]"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
comments:
---

The default options value is `- . ,`.  The default delimiter within the options
value is a space.  If you need to use a space as one of the options, set a
custom delimiter.

Numbers greater than or equal to 5 are rounded up. For example, if precision is set to `0`, `1.5` becomes `2`, and `1.4` becomes `1`.

```
{{ lang.NumFmt 2 12345.6789 }} → 12,345.68
{{ lang.NumFmt 2 12345.6789 "- , ." }} → 12.345,68
{{ lang.NumFmt 0 -12345.6789 "- . ," }} → -12,346
{{ lang.NumFmt 6 -12345.6789 "- ." }} → -12345.678900
{{ lang.NumFmt 6 -12345.6789 "-|.| " "|" }} → -1 2345.678900
{{ -98765.4321 | lang.NumFmt 2 }} → -98,765.43
```
