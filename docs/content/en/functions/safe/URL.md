---
title: safe.URL
description: Declares the given string as a safe URL or URL substring.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [safeURL]
    returnType: template.URL
    signatures: [safe.URL INPUT]
aliases: [/functions/safeurl]
---

## Introduction

{{% include "/_common/functions/go-html-template-package.md" %}}

## Usage

Use the `safe.URL` function to encapsulate a known safe URL or URL substring. Schemes other than the following are considered unsafe:

- `http:`
- `https:`
- `mailto:`

Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.

See the [Go documentation] for details.

## Example

Without a safe declaration:

```go-html-template
{{ $href := "irc://irc.freenode.net/#golang" }}
<a href="{{ $href }}">IRC</a>
```

Hugo renders the above to:

```html
<a href="#ZgotmplZ">IRC</a>
```

> [!note]
> `ZgotmplZ` is a special value that indicates that unsafe content reached a CSS or URL context at runtime.

To declare the string as safe:

```go-html-template
{{ $href := "irc://irc.freenode.net/#golang" }}
<a href="{{ $href | safeURL }}">IRC</a>
```

Hugo renders the above to:

```html
<a href="irc://irc.freenode.net/#golang">IRC</a>
```

[Go documentation]: https://pkg.go.dev/html/template#URL
