---
title: transform.HTMLEscape
description: Returns the given string, escaping special characters by replacing them with HTML entities.
categories: []
keywords: []
action:
  aliases: [htmlEscape]
  related:
    - functions/transform/HTMLUnescape
  returnType: string
  signatures: [transform.HTMLEscape INPUT]
aliases: [/functions/htmlescape]
---

The `transform.HTMLEscape` function escapes five special characters by replacing them with [HTML entities]:

- `&` → `&amp;`
- `<` → `&lt;`
- `>` → `&gt;`
- `'` → `&#39;`
- `"` → `&#34;`

For example:

```go-html-template
{{ htmlEscape "Lilo & Stitch" }} → Lilo &amp; Stitch
{{ htmlEscape "7 > 6" }} → 7 &gt; 6
```

[html entities]: https://developer.mozilla.org/en-US/docs/Glossary/Entity
