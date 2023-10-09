---
title: title
description: Converts the provided string to title case.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature:
  - "title STRING"
  - "strings.Title STRING"
relatedfuncs: []
---

```go-html-template
{{ title "table of contents (TOC)" }} → "Table of Contents (TOC)"
```

By default, Hugo adheres to the capitalization rules in the [Associated Press (AP) Stylebook]. Change your [site configuration] if you would prefer to follow the [Chicago Manual of Style], or to use Go's convention of capitalizing every word.

[Associated Press (AP) Stylebook]: https://www.apstylebook.com/
[Chicago Manual of Style]: https://www.chicagomanualofstyle.org/home.html
[site configuration]: /getting-started/configuration/#configure-title-case
