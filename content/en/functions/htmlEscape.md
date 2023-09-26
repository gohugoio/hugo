---
title: htmlEscape
description: Returns the given string with the reserved HTML codes escaped.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: transform
relatedFuncs:
  - transform.HTMLEscape
  - transform.HTMLUnescape
signature:
  - transform.HTMLEscape INPUT
  - htmlEscape INPUT
---

In the result `&` becomes `&amp;` and so on. It escapes only: `<`, `>`, `&`, `'` and `"`.

```go-html-template
{{ htmlEscape "Hugo & Caddy > WordPress & Apache" }} → "Hugo &amp; Caddy &gt; WordPress &amp; Apache"
```
