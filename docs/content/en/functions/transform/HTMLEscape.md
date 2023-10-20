---
title: transform.HTMLEscape
linkTitle: htmlEscape
description: Returns the given string with the reserved HTML codes escaped.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [htmlEscape]
  returnType: string
  signatures: [transform.HTMLEscape INPUT]
relatedFunctions:
  - transform.HTMLEscape
  - transform.HTMLUnescape
aliases: [/functions/htmlescape]
---

In the result `&` becomes `&amp;` and so on. It escapes only: `<`, `>`, `&`, `'` and `"`.

```go-html-template
{{ htmlEscape "Hugo & Caddy > WordPress & Apache" }} → "Hugo &amp; Caddy &gt; WordPress &amp; Apache"
```
