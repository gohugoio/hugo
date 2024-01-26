---
title: strings.FindRE
description: Returns a slice of strings that match the regular expression.
categories: []
keywords: []
action:
  aliases: [findRE]
  related:
    - functions/strings/FindRESubmatch
    - functions/strings/Replace
    - functions/strings/ReplaceRE
  returnType: '[]string'
  signatures: ['strings.FindRE PATTERN INPUT [LIMIT]']
aliases: [/functions/findre]
---
By default, `findRE` finds all matches. You can limit the number of matches with an optional LIMIT argument.

{{% include "functions/_common/regular-expressions.md" %}}

This example returns a slice of all second level headings (`h2` elements) within the rendered `.Content`:

```go-html-template
{{ findRE `(?s)<h2.*?>.*?</h2>` .Content }}
```

The `s` flag causes `.` to match `\n` as well, allowing us to find an `h2` element that contains newlines.

To limit the number of matches to one:

```go-html-template
{{ findRE `(?s)<h2.*?>.*?</h2>` .Content 1 }}
```

{{% note %}}
You can write and test your regular expression using [regex101.com](https://regex101.com/). Be sure to select the Go flavor before you begin.
{{% /note %}}
