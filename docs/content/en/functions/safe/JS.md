---
title: safe.JS
description: Declares the given string as a safe JavaScript expression.
categories: []
keywords: []
action:
  aliases: [safeJS]
  related:
    - functions/safe/CSS
    - functions/safe/HTML
    - functions/safe/HTMLAttr
    - functions/safe/JSStr
    - functions/safe/URL
  returnType: template.JS
  signatures: [safe.JS INPUT]
toc: true
aliases: [/functions/safejs]
---

## Introduction

{{% include "functions/_common/go-html-template-package.md" %}}

## Usage

Use the `safe.JS` function to encapsulate a known safe EcmaScript5 Expression.

Template authors are responsible for ensuring that typed expressions do not break the intended precedence and that there is no statement/expression ambiguity as when passing an expression like `{ foo: bar() }\n['foo']()`, which is both a valid Expression and a valid Program with a very different meaning.

Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.

Using the `safe.JS` function to include valid but untrusted JSON is not safe. A safe alternative is to parse the JSON with the [`transform.Unmarshal`] function and then pass the resultant object into the template, where it will be converted to sanitized JSON when presented in a JavaScript context.

[`transform.Unmarshal`]: /functions/transform/unmarshal/

See the [Go documentation] for details.

[Go documentation]: https://pkg.go.dev/html/template#JS

## Example

Without a safe declaration:

```go-html-template
{{ $js := "x + y" }}
<script>const a = {{ $js }}</script>
```

Hugo renders the above to:

```html
<script>const a = "x + y"</script>
```

To declare the string as safe:

```go-html-template
{{ $js := "x + y" }}
<script>const a = {{ $js | safeJS }}</script>
```

Hugo renders the above to:

```html
<script>const a = x + y</script>
```
