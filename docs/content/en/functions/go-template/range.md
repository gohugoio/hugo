---
title: range
description: Iterates over a non-empty collection, binds context (the dot) to successive elements, and executes the block.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: 
    signatures: [range COLLECTION]
aliases: [/functions/range]
---

The collection may be a slice, a map, or an integer.

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ . }} → foo bar baz
{{ end }}
```

Use with the [`else`] statement:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  <p>{{ . }}</p>
{{ else }}
  <p>The collection is empty</p>
{{ end }}
```

Within a range block:

- Use the [`continue`] statement to stop the innermost iteration and continue to the next iteration
- Use the [`break`] statement to stop the innermost iteration and bypass all remaining iterations

## Understanding context

At the top of a page template, the [context](g) (the dot) is a `Page` object. Within the `range` block, the context is bound to each successive element.

With this contrived example:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ .Title }} 
{{ end }}
```

Hugo will throw an error:

```text
can't evaluate field Title in type int
```

The error occurs because we are trying to use the `.Title` method on a string instead of a `Page` object. Within the `range` block, if we want to render the page title, we need to get the context passed into the template.

> [!note]
> Use the `$` to get the context passed into the template.

This template will render the page title three times:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ $.Title }} 
{{ end }}
```

> [!note]
> Gaining a thorough understanding of context is critical for anyone writing template code.

## Examples

### Slice of scalars

This template code:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  <p>{{ . }}</p>
{{ end }}
```

Is rendered to:

```html
<p>foo</p>
<p>bar</p>
<p>baz</p>
```

This template code:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $v := $s }}
  <p>{{ $v }}</p>
{{ end }}
```

Is rendered to:

```html
<p>foo</p>
<p>bar</p>
<p>baz</p>
```

This template code:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $k, $v := $s }}
  <p>{{ $k }}: {{ $v }}</p>
{{ end }}
```

Is rendered to:

```html
<p>0: foo</p>
<p>1: bar</p>
<p>2: baz</p>
```

### Slice of maps

This template code:

```go-html-template
{{ $m := slice
  (dict "name" "John" "age" 30)
  (dict "name" "Will" "age" 28)
  (dict "name" "Joey" "age" 24)
}}
{{ range $m }}
  <p>{{ .name }} is {{ .age }}</p>
{{ end }}
```

Is rendered to:

```html
<p>John is 30</p>
<p>Will is 28</p>
<p>Joey is 24</p>
```

### Slice of pages

This template code:

```go-html-template
{{ range where site.RegularPages "Type" "articles" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

Is rendered to:

```html
<h2><a href="/articles/article-3/">Article 3</a></h2>
<h2><a href="/articles/article-2/">Article 2</a></h2>
<h2><a href="/articles/article-1/">Article 1</a></h2>
```

### Maps

This template code:

```go-html-template
{{ $m :=  dict "name" "John" "age" 30 }}
{{ range $k, $v := $m }}
  <p>key = {{ $k }} value = {{ $v }}</p>
{{ end }}
```

Is rendered to:

```go-html-template
<p>key = age value = 30</p>
<p>key = name value = John</p>
```

Unlike ranging over an array or slice, Hugo sorts by key when ranging over a map.

### Integers

{{< new-in 0.123.0 />}}

Ranging over a positive integer `n` executes the block `n` times, with the context starting at zero and incrementing by one in each iteration.

```go-html-template
{{ $s := slice }}
{{ range 1 }}
  {{ $s = $s | append . }}
{{ end }}
{{ $s }} → [0]
```

```go-html-template
{{ $s := slice }}
{{ range 3 }}
  {{ $s = $s | append . }}
{{ end }}
{{ $s }} → [0 1 2]
```

Ranging over a non-positive integer executes the block zero times.

{{% include "/_common/functions/go-template/text-template.md" %}}

[`break`]: /functions/go-template/break/
[`continue`]: /functions/go-template/continue/
[`else`]: /functions/go-template/else/
