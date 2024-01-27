---
title: strings.Title
description: Returns the given string, converting it to title case.
categories: []
keywords: []
action:
  aliases: [title]
  related:
    - functions/strings/FirstUpper
    - functions/strings/ToLower
    - functions/strings/ToUpper
  returnType: string
  signatures: [strings.Title STRING]
aliases: [/functions/title]
---

```go-html-template
{{ title "table of contents (TOC)" }} → Table of Contents (TOC)
```

By default, Hugo follows the capitalization rules published in the [Associated Press Stylebook]. Change your [site configuration] if you would prefer to:

- Follow the capitalization rules published in the [Chicago Manual of Style]
- Capitalize the first letter of every word
- Capitalize the first letter of the first word
- Disable the effects of the `title` function

The last option is useful if your theme uses the `title` function, and you would prefer to manually capitalize strings as needed.

[Associated Press Stylebook]: https://www.apstylebook.com/
[Chicago Manual of Style]: https://www.chicagomanualofstyle.org/home.html
[site configuration]: /getting-started/configuration/#configure-title-case
