---
title: strings.Truncate
linkTitle: truncate
description: Truncates a text to a max length without cutting words or leaving unclosed HTML tags.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [truncate]
  returnType: template.HTML
  signatures: ['strings.Truncate SIZE [ELLIPSIS] INPUT']
relatedFunctions: []
aliases: [/functions/truncate]
---

Since Go templates are HTML-aware, `truncate` will intelligently handle normal strings vs HTML strings:

```go-html-template
{{ "<em>Keep my HTML</em>" | safeHTML | truncate 10 }} → <em>Keep my …</em>
```

{{% note %}}
If you have a raw string that contains HTML tags you want to remain treated as HTML, you will need to convert the string to HTML using the [`safeHTML` template function](/functions/safe/html) before sending the value to truncate. Otherwise, the HTML tags will be escaped when passed through the `truncate` function.
{{% /note %}}
