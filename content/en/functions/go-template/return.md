---
title: return
description: Terminates execution of the current template and, when used within a partial template, returns the given value, if any.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: any
    signatures: ['return [VALUE]']
---

The `return` statement is a non-standard extension to Go's [`text/template`][] package. It terminates execution of the current template only; execution continues in any calling template. When used within a _partial_ template, the `return` statement may also return a value to the caller.

{{< new-in 0.166.0 >}}
In earlier versions, the `return` statement was only supported within partial templates, limited to one `return` statement per template, executed regardless of its position within logical blocks. The `return` statement now follows normal flow control: you can use it in any template, in any position, any number of times, and Hugo executes it only when reached.
{{< /new-in >}}

## Return a value

Within a _partial_ template, the `return` statement may return a value of any data type: [`bool`](g), [`float`](g), [`int`](g), [`map`](g), [`resource`](g), [`slice`](g), [`string`](g), and others. When returning a value, any output rendered before the `return` statement is discarded.

Using the `return` statement with a value in any other template type produces an error.

For example, a _partial_ template that returns a string value:

```go-html-template {file="layouts/_partials/parity.html"}
{{ if math.ModBool . 2 }}
  {{ return "even" }}
{{ end }}
{{ return "odd" }}
```

```go-html-template {file="layouts/single.html"}
{{ partial "parity.html" 42 }} → even
```

A more practical example is a _partial_ template that returns the cover image of the first page in a section that has one:

```go-html-template {file="layouts/_partials/section-cover.html"}
{{ range .Pages }}
  {{ with .Resources.GetMatch "cover.*" }}
    {{ return . }}
  {{ end }}
{{ end }}
```

```go-html-template {file="layouts/section.html"}
{{ with partial "section-cover.html" . }}
  <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
{{ end }}
```

## Return early

Use the `return` statement without a value to stop execution of the current template:

```go-html-template {file="layouts/single.html"}
<h2>{{ .Title }}</h2>
{{ if .Draft }}
  <p>This article is a draft.</p>
  {{ return }}
{{ end }}
{{ .Content }}
```

Within a _shortcode_ template, use a `return` statement after each validation check to avoid deeply nested conditional blocks:

```go-html-template {file="layouts/_shortcodes/img.html"}
{{ if not (.Get "src") }}
  {{ errorf "The %q shortcode requires a src argument. See %s" .Name .Position }}
  {{ return }}
{{ end }}

{{ if not (.Get "alt") }}
  {{ errorf "The %q shortcode requires an alt argument. See %s" .Name .Position }}
  {{ return }}
{{ end }}

<img src="{{ .Get "src" }}" alt="{{ .Get "alt" }}">
```

## Limitations

The `return` statement must be the last command in a pipeline. This produces an error:

```go-html-template
{{ return "even" | strings.ToUpper }}
```

[`text/template`]: https://pkg.go.dev/text/template
