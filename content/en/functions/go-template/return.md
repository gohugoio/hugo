---
title: return
description: Used within partial templates, terminates template execution and returns the given value, if any.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: any
    signatures: ['return [VALUE]']
---

The `return` statement is a non-standard extension to Go's [text/template package]. Used within partial templates, the `return` statement terminates template execution and returns the given value, if any.

The returned value may be of any data type including, but not limited to, [`bool`](g), [`float`](g), [`int`](g), [`map`](g), [`resource`](g), [`slice`](g), or [`string`](g).

A `return` statement without a value returns an empty string of type `template.HTML`.

> [!note]
> Unlike `return` statements in other languages, Hugo executes the first occurrence of the `return` statement regardless of its position within logical blocks. See [usage](#usage) notes below.

## Example

By way of example, let's create a partial template that _renders_ HTML, describing whether the given number is odd or even:

```go-html-template {file="layouts/_partials/odd-or-even.html"}
{{ if math.ModBool . 2 }}
  <p>{{ . }} is even</p>
{{ else }}
  <p>{{ . }} is odd</p>
{{ end }}
```

When called, the partial renders HTML:

```go-html-template
{{ partial "odd-or-even.html" 42 }} → <p>42 is even</p>
```

Instead of rendering HTML, let's create a partial that _returns_ a boolean value, reporting whether the given number is even:

```go-html-template {file="layouts/_partials/is-even.html"}
{{ return math.ModBool . 2 }}
```

With this template:

```go-html-template
{{ $number := 42 }}
{{ if partial "is-even.html" $number }}
  <p>{{ $number }} is even</p>
{{ else }}
  <p>{{ $number }} is odd</p>
{{ end }}
```

Hugo renders:

```html
<p>42 is even</p>
```

## Usage

> [!note]
> Unlike `return` statements in other languages, Hugo executes the first occurrence of the `return` statement regardless of its position within logical blocks.

A partial that returns a value must contain only one `return` statement, placed at the end of the template.

For example:

```go-html-template {file="layouts/_partials/is-even.html"}
{{ $result := false }}
{{ if math.ModBool . 2 }}
  {{ $result = "even" }}
{{ else }}
  {{ $result = "odd" }}
{{ end }}
{{ return $result }}
```

> [!note]
> The construct below is incorrect; it contains more than one `return` statement.

```go-html-template {file="layouts/_partials/do-not-do-this.html"}
{{ if math.ModBool . 2 }}
  {{ return "even" }}
{{ else }}
  {{ return "odd" }}
{{ end }}
```

[text/template package]: https://pkg.go.dev/text/template
