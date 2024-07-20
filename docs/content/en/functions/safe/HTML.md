---
title: safe.HTML
description: Declares the given string as a safeHTML string.
categories: []
keywords: []
action:
  aliases: [safeHTML]
  related:
    - functions/safe/CSS
    - functions/safe/HTMLAttr
    - functions/safe/JS
    - functions/safe/JSStr
    - functions/safe/URL
  returnType: template.HTML
  signatures: [safe.HTML INPUT]
toc: true
aliases: [/functions/safehtml]
---

## Introduction

{{% include "functions/_common/go-html-template-package.md" %}}

## Usage

Use the `safe.HTML` function to encapsulate a known safe HTML document fragment. It should not be used for HTML from a third-party, or HTML with unclosed tags or comments.

Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.

See the [Go documentation] for details.

[Go documentation]: https://pkg.go.dev/html/template#HTML

## Example

Without a safe declaration:

```go-html-template
{{ $html := "<em>emphasized</em>" }}
{{ $html }}
```

Hugo renders the above to:

```html
&lt;em&gt;emphasized&lt;/em&gt;
```

To declare the string as safe:

```go-html-template
{{ $html := "<em>emphasized</em>" }}
{{ $html | safeHTML }}
```

Hugo renders the above to:

```html
<em>emphasized</em>
```
