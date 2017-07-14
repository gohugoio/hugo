---
title: htmlEscape
linktitle:
description: Returns the given string with the critical reserved HTML codes escaped.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
#tags: [strings, html]
ns:
signature: ["htmlEscape INPUT"]
workson: []
hugoversion:
relatedfuncs: [htmlUnescape]
deprecated: false
aliases: []
---

`htmlEscape` returns the given string with the critical reserved HTML codes escaped, such that `&` becomes `&amp;` and so on. It escapes only: `<`, `>`, `&`, `'` and `"`.

Bear in mind that, unless content is passed to `safeHTML`, output strings are escaped usually by the processor anyway.

```
{{ htmlEscape "Hugo & Caddy > Wordpress & Apache" }} → "Hugo &amp; Caddy &gt; Wordpress &amp; Apache"
```
