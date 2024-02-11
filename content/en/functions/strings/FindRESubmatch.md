---
title: strings.FindRESubmatch
description: Returns a slice of all successive matches of the regular expression. Each element is a slice of strings holding the text of the leftmost match of the regular expression and the matches, if any, of its subexpressions.
categories: []
keywords: []
action:
  aliases: [findRESubmatch]
  related:
    - functions/strings/FindRE
    - functions/strings/Replace
    - functions/strings/ReplaceRE
  returnType: '[][]string'
  signatures: ['strings.FindRESubmatch PATTERN INPUT [LIMIT]']
aliases: [/functions/findresubmatch]
---

By default, `findRESubmatch` finds all matches. You can limit the number of matches with an optional LIMIT argument. A return value of nil indicates no match.

{{% include "functions/_common/regular-expressions.md" %}}

## Demonstrative examples

```go-html-template
{{ findRESubmatch `a(x*)b` "-ab-" }} → [["ab" ""]]
{{ findRESubmatch `a(x*)b` "-axxb-" }} → [["axxb" "xx"]]
{{ findRESubmatch `a(x*)b` "-ab-axb-" }} → [["ab" ""] ["axb" "x"]]
{{ findRESubmatch `a(x*)b` "-axxb-ab-" }} → [["axxb" "xx"] ["ab" ""]]
{{ findRESubmatch `a(x*)b` "-axxb-ab-" 1 }} → [["axxb" "xx"]]
```

## Practical example

This Markdown:

```text
- [Example](https://example.org)
- [Hugo](https://gohugo.io)
```

Produces this HTML:

```html
<ul>
  <li><a href="https://example.org">Example</a></li>
  <li><a href="https://gohugo.io">Hugo</a></li>
</ul>
```

To match the anchor elements, capturing the link destination and text:

```go-html-template
{{ $regex := `<a\s*href="(.+?)">(.+?)</a>` }}
{{ $matches := findRESubmatch $regex .Content }}
```

Viewed as JSON, the data structure of `$matches` in the code above is:

```json
[
  [
    "<a href=\"https://example.org\"></a>Example</a>",
    "https://example.org",
    "Example"
  ],
  [
    "<a href=\"https://gohugo.io\">Hugo</a>",
    "https://gohugo.io",
    "Hugo"
  ]
]
```

To render the `href` attributes:

```go-html-template
{{ range $matches }}
  {{ index . 1 }}
{{ end }}
```

Result:

```text
https://example.org
https://gohugo.io
```

{{% note %}}
You can write and test your regular expression using [regex101.com](https://regex101.com/). Be sure to select the Go flavor before you begin.
{{% /note %}}
