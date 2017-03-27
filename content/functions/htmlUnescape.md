---
title: htmlUnescape
linktitle: htmlUnescape
description: Returns the given string with HTML escape codes un-escaped.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
ns:
signature: ["htmlUnescape INPUT"]
workson: []
hugoversion:
relatedfuncs: [htmlEscape]
deprecated: false
aliases: []
---

`htmlUnescape` returns the given string with HTML escape codes un-escaped. This un-escapes more codes than `htmlEscape` escapes, including `#` codes and pre-UTF8 escapes for accented characters. It defers completely to the Go `html.UnescapeString` function, so functionality is consistent with that codebase.

Remember to pass the output of this to `safeHTML` if fully un-escaped characters are desired. Otherwise, the output will be escaped again as normal.

```
{{ htmlUnescape "Hugo &amp; Caddy &gt; Wordpress &amp; Apache" }} → "Hugo & Caddy > Wordpress & Apache"
```
