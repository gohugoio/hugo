---
title: with
description: Binds context (the dot) to the expression and executes the block if expression is truthy.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/if
    - functions/go-template/else
    - functions/go-template/end
    - functions/collections/IsSet
  returnType:
  signatures: [with EXPR]
aliases: [/functions/with]
toc: true
---

{{% include "functions/go-template/_common/truthy-falsy.md" %}}

```go-html-template
{{ $var := "foo" }}
{{ with $var }}
  {{ . }} → foo
{{ end }}
```

Use with the [`else`] statement:

```go-html-template
{{ $var := "foo" }}
{{ with $var }}
  {{ . }} → foo
{{ else }}
  {{ print "var is falsy" }}
{{ end }}
```

Use `else with` to check multiple conditions:

```go-html-template
{{ $v1 := 0 }}
{{ $v2 := 42 }}
{{ with $v1 }}
  {{ . }}
{{ else with $v2 }}
  {{ . }} → 42
{{ else }}
  {{ print "v1 and v2 are falsy" }}
{{ end }}
```

Initialize a variable, scoped to the current block:

```go-html-template
{{ with $var := 42 }}
  {{ . }} → 42
  {{ $var }} → 42
{{ end }}
{{ $var }} → undefined
```

## Understanding context

At the top of a page template, the [context] (the dot) is a `Page` object. Inside of the `with` block, the context is bound to the value passed to the `with` statement.

With this contrived example:

```go-html-template
{{ with 42 }}
  {{ .Title }}
{{ end }}
```

Hugo will throw an error:

    can't evaluate field Title in type int

The error occurs because we are trying to use the `.Title` method on an integer instead of a `Page` object. Inside of the `with` block, if we want to render the page title, we need to get the context passed into the template.

{{% note %}}
Use the `$` to get the context passed into the template.
{{% /note %}}

This template will render the page title as desired:

```go-html-template
{{ with 42 }}
  {{ $.Title }}
{{ end }}
```

{{% note %}}
Gaining a thorough understanding of context is critical for anyone writing template code.
{{% /note %}}

[context]: /getting-started/glossary/#context

{{% include "functions/go-template/_common/text-template.md" %}}

[`else`]: /functions/go-template/else/
