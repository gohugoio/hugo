---
title: strings.Title
linkTitle: title
description: Converts the provided string to title case.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [title]
  returnType: string
  signatures: [strings.Title STRING]
relatedFunctions:
  - strings.FirstUpper
  - strings.Title
  - strings.ToLower
  - strings.ToUpper
aliases: [/functions/title]
---

```go-html-template
{{ title "table of contents (TOC)" }} → "Table of Contents (TOC)"
```

By default, Hugo adheres to the capitalization rules in the [Associated Press (AP) Stylebook]. Change your [site configuration] if you would prefer to follow the [Chicago Manual of Style], or to use Go's convention of capitalizing every word.

[Associated Press (AP) Stylebook]: https://www.apstylebook.com/
[Chicago Manual of Style]: https://www.chicagomanualofstyle.org/home.html
[site configuration]: /getting-started/configuration/#configure-title-case
