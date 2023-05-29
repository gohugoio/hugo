---
title: findRE
description: Returns a slice of strings that match the regular expression.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [regex]
signature:
  - "findRE PATTERN INPUT [LIMIT]"
  - "strings.FindRE PATTERN INPUT [LIMIT]"
relatedfuncs: [findRESubmatch, replaceRE]
---
By default, `findRE` finds all matches. You can limit the number of matches with an optional LIMIT parameter.

When specifying the regular expression, use a raw [string literal] (backticks) instead of an interpreted string literal (double quotes) to simplify the syntax. With an interpreted string literal you must escape backslashes.

[string literal]: https://go.dev/ref/spec#String_literals

This function uses the [RE2] regular expression library. See the [RE2 syntax documentation] for details. Note that the RE2 `\C` escape sequence is not supported.

[RE2]: https://github.com/google/re2/
[RE2 syntax documentation]: https://github.com/google/re2/wiki/Syntax/

{{% note %}}
The RE2 syntax is a subset of that accepted by [PCRE], roughly speaking, and with various [caveats].

[caveats]: https://swtch.com/~rsc/regexp/regexp3.html#caveats
[PCRE]: https://www.pcre.org/
{{% /note %}}

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
