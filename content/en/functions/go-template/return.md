---
title: return
description: Used within partial templates, terminates template execution and returns the given value, if any.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/partials/Include
    - functions/partials/IncludeCached
  returnType: any
  signatures: ['return [VALUE]']
toc: true
---

The `return` statement is a custom addition to Go's [text/template] package. Used within partial templates, the `return` statement terminates template execution and returns the given value, if any.

The returned value may be of any data type including, but not limited to, [`bool`], [`float`], [`int`], [`map`], [`resource`], [`slice`], and [`string`].

A `return` statement without a value returns an empty string of type `template.HTML`.

[`bool`]: /getting-started/glossary/#bool
[`float`]: /getting-started/glossary/#float
[`int`]: /getting-started/glossary/#int
[`map`]: /getting-started/glossary/#map
[`resource`]: /getting-started/glossary/#resource
[`slice`]: /getting-started/glossary/#slice
[`string`]: /getting-started/glossary/#string
[text/template]: https://pkg.go.dev/text/template

{{% note %}}
Unlike `return` statements in other languages, Hugo executes the first occurrence of the `return` statement regardless of its position within logical blocks. See [usage](#usage) notes below.
{{% /note %}}

## Example

By way of example, let's create a partial template that _renders_ HTML, describing whether the given number is odd or even:

{{< code file="layouts/partials/odd-or-even.html" >}}
{{ if math.ModBool . 2 }}
  <p>{{ . }} is even</p>
{{ else }}
  <p>{{ . }} is odd</p>
{{ end }}
{{< /code >}}

When called, the partial renders HTML:

```go-html-template
{{ partial "odd-or-even.html" 42 }} → <p>42 is even</p>
```

Instead of rendering HTML, let's create a partial that _returns_ a boolean value, reporting whether the given number is even:

{{< code file="layouts/partials/is-even.html" >}}
{{ return math.ModBool . 2 }}
{{< /code >}}

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

See additional examples in the [partial templates] section.

[partial templates]: /templates/partial/#returning-a-value-from-a-partial

## Usage

{{% note %}}
Unlike `return` statements in other languages, Hugo executes the first occurrence of the `return` statement regardless of its position within logical blocks.
{{% /note %}}

A partial that returns a value must contain only one `return` statement, placed at the end of the template.

For example:

{{< code file="layouts/partials/is-even.html" >}}
{{ $result := false }}
{{ if math.ModBool . 2 }}
  {{ $result = "even" }}
{{ else }}
  {{ $result = "odd" }}
{{ end }}
{{ return $result }}
{{< /code >}}

{{% note %}}
The construct below is incorrect; it contains more than one `return` statement.
{{% /note %}}

{{< code file="layouts/partials/do-not-do-this.html" >}}
{{ if math.ModBool . 2 }}
  {{ return "even" }}
{{ else }}
  {{ return "odd" }}
{{ end }}
{{< /code >}}
