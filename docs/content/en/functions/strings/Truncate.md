---
title: strings.Truncate
description: Returns the given string, truncating it to a maximum length without cutting words or leaving unclosed HTML tags.
categories: []
keywords: []
action:
  aliases: [truncate]
  related: []
  returnType: template.HTML
  signatures: ['strings.Truncate SIZE [ELLIPSIS] INPUT']
aliases: [/functions/truncate]
---

Since Go templates are HTML-aware, `truncate` will intelligently handle normal strings vs HTML strings:

```go-html-template
{{ "<em>Keep my HTML</em>" | safeHTML | truncate 10 }} → <em>Keep my …</em>
```

{{% note %}}
If you have a raw string that contains HTML tags you want to remain treated as HTML, you will need to convert the string to HTML using the [`safeHTML`]function before sending the value to `truncate`. Otherwise, the HTML tags will be escaped when passed through the `truncate` function.

[`safeHTML`]: /functions/safe/html/
{{% /note %}}
