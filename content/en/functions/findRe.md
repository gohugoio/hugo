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
relatedfuncs: [replaceRE]
---
By default, the `findRE` function finds all matches. You can limit the number of matches with an optional LIMIT parameter.

When specifying the regular expression, use a raw [string literal] (backticks) instead of an interpreted string literal (double quotes) to simplify the syntax. With an interpreted string literal you must escape backslashes.

The syntax of the regular expression is the same general syntax used by Perl, Python, and other languages. More precisely, it is the syntax accepted by [RE2] except for `\C`.

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

## findRESubmatch

In Hugo 0.110.0 we added a variant of `findRE` that returns a slice of strings holding the text of the leftmost match of the regular expression in s and the matches, if any, of its subexpressions.

This:

```go-html-template
{{ findRESubmatch §§<a\s*href="(.+?)">(.+?)</a>§§ §§<li><a href="#foo">Foo</a></li> <li><a href="#bar">Bar</a></li>§§ | print | safeHTML }}
```

Will print:

```
[[<a href=\"#foo\">Foo</a> #foo Foo] [<a href=\"#bar\">Bar</a> #bar Bar]]
```

{{< new-in "0.110.0" >}}


[RE2]: https://github.com/google/re2/wiki/Syntax
[string literal]: https://go.dev/ref/spec#String_literals
