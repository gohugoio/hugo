---
title: range
description: Iterates over slices, maps, and page collections.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: 
  signatures: [range COLLECTION]
relatedFunctions:
  - with
  - range
aliases: [/functions/range]
toc: true
---

{{% readfile file="/functions/_common/go-template-functions.md" %}}

## Slices

Template:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  <p>{{ . }}</p>
{{ end }}
```

Result:

```html
<p>foo</p>
<p>bar</p>
<p>baz</p>
```

Template:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $v := $s }}
  <p>{{ $v }}</p>
{{ end }}
```

Result:

```html
<p>foo</p>
<p>bar</p>
<p>baz</p>
```

Template:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $k, $v := $s }}
  <p>{{ $k }}: {{ $v }}</p>
{{ end }}
```

Result:

```html
<p>0: foo</p>
<p>1: bar</p>
<p>2: baz</p>
```

## Maps

Template:

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

Result:

```html
<p>John is 30</p>
<p>Will is 28</p>
<p>Joey is 24</p>
```

## Page collections

Template:

```go-html-template
{{ range where site.RegularPages "Type" "articles" }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

Result:

```html
<h2><a href="/articles/article-3/">Article 3</a></h2>
<h2><a href="/articles/article-2/">Article 2</a></h2>
<h2><a href="/articles/article-1/">Article 1</a></h2>
```

## Break

Use the `break` statement to stop the innermost iteration and bypass all remaining iterations.

Template:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ if eq . "bar" }}
    {{ break }}
  {{ end }}
  <p>{{ . }}</p>
{{ end }}
```

Result:

```html
<p>foo</p>
```

## Continue

Use the `continue` statement to stop the innermost iteration and continue to the next iteration.

Template:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ if eq . "bar" }}
    {{ continue }}
  {{ end }}
  <p>{{ . }}</p>
{{ end }}
```

Result:

```html
<p>foo</p>
<p>baz</p>
```
