---
title: htmlEscape
description: Returns the given string with the reserved HTML codes escaped.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings, html]
signature: ["htmlEscape INPUT"]
relatedfuncs: [htmlUnescape]
---

In the result `&` becomes `&amp;` and so on. It escapes only: `<`, `>`, `&`, `'` and `"`.

```go-html-template
{{ htmlEscape "Hugo & Caddy > WordPress & Apache" }} → "Hugo &amp; Caddy &gt; WordPress &amp; Apache"
```
