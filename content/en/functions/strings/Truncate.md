---
title: strings.Truncate
description: Returns the given string, truncating it to a maximum length without cutting words or leaving unclosed HTML tags.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [truncate]
    returnType: template.HTML
    signatures: ['strings.Truncate SIZE [ELLIPSIS] STRING']
aliases: [/functions/truncate]
---

When truncating a value marked as safe HTML, such as one returned by the [`safe.HTML`][] function, `strings.Truncate` closes any tag left open by the truncation instead of cutting in the middle of it:

```go-html-template
{{ "<em>Keep my HTML</em>" | safeHTML | strings.Truncate 10 }} → <em>Keep my …</em>
```

> [!NOTE]
> If you have a raw string containing HTML tags that you want treated as HTML, convert it with the [`safe.HTML`][] function first. Otherwise, `strings.Truncate` escapes the tags.

[`safe.HTML`]: /functions/safe/html/
