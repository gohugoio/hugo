---
title: truncate
description: Truncates a text to a max length without cutting words or leaving unclosed HTML tags.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: strings
relatedFuncs: []
signature:
  - strings.Truncate SIZE [ELLIPSIS] INPUT
  - truncate SIZE [ELLIPSIS] INPUT
---

Since Go templates are HTML-aware, `truncate` will intelligently handle normal strings vs HTML strings:

```go-html-template
{{ "<em>Keep my HTML</em>" | safeHTML | truncate 10 }}` → <em>Keep my …</em>`
```

{{% note %}}
If you have a raw string that contains HTML tags you want to remain treated as HTML, you will need to convert the string to HTML using the [`safeHTML` template function](/functions/safehtml) before sending the value to truncate. Otherwise, the HTML tags will be escaped when passed through the `truncate` function.
{{% /note %}}
