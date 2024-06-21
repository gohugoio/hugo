---
title: safe.CSS
description: Declares the given string as a safe CSS string.
categories: []
keywords: []
action:
  aliases: [safeCSS]
  related:
    - functions/safe/HTML
    - functions/safe/HTMLAttr
    - functions/safe/JS
    - functions/safe/JSStr
    - functions/safe/URL
  returnType: template.CSS
  signatures: [safe.CSS INPUT]
toc: true
aliases: [/functions/safecss]
---

## Introduction

{{% include "functions/_common/go-html-template-package.md" %}}

## Usage

Use the `safe.CSS` function to encapsulate known safe content that matches any of:

1. The CSS3 stylesheet production, such as `p { color: purple }`.
2. The CSS3 rule production, such as `a[href=~"https:"].foo#bar`.
3. CSS3 declaration productions, such as `color: red; margin: 2px`.
4. The CSS3 value production, such as `rgba(0, 0, 255, 127)`.

Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.

See the [Go documentation] for details.

[Go documentation]: https://pkg.go.dev/html/template#CSS

## Example

Without a safe declaration:

```go-html-template
{{ $style := "color: red;" }}
<p style="{{ $style }}">foo</p>
```

Hugo renders the above to:

```html
<p style="ZgotmplZ">foo</p>
```

{{% note %}}
`ZgotmplZ` is a special value that indicates that unsafe content reached a CSS or URL context at runtime.
{{% /note %}}

To declare the string as safe:

```go-html-template
{{ $style := "color: red;" }}
<p style="{{ $style | safeCSS }}">foo</p>
```

Hugo renders the above to:

```html
<p style="color: red;">foo</p>
```
