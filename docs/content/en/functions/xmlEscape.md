---
title: xmlEscape
linktitle:
description: Returns the given string with the reserved XML characters escaped.
godocref:
date: 2018-10-13
publishdate: 2018-10-13
lastmod: 2018-10-13
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings, xml]
signature: ["xmlEscape INPUT"]
workson: []
hugoversion:
relatedfuncs:
deprecated: false
aliases: []
---

In the result `&` becomes `&amp;` and so on. It escapes characters with
special meaning in XML 1.0. If a character is not valid for XML 1.0, it
gets replaced with the Unicode replacement character (U+FFFD).

```
{{ xmlEscape "Hugo & Caddy > Wordpress & Apache" }} → "Hugo &amp; Caddy &gt; Wordpress &amp; Apache"
```
