---
title: partials.IncludeCached
description: Executes the given template and caches the result, optionally passing one or more variant keys. If the partial template contains a return statement, returns the given value, else returns the rendered output.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [partialCached]
    returnType: any
    signatures: ['partials.IncludeCached LAYOUT CONTEXT [VARIANT...]']
aliases: [/functions/partialcached]
---

Without a [`return`][] statement, the `partialCached` function returns a string of type `template.HTML`. With a `return` statement, the `partialCached` function can return any data type.

The `partialCached` function can offer significant performance gains for complex templates that don't need to be re-rendered on every invocation.

> [!NOTE]
> Each site (or language) has its own `partialCached` cache, so each site will execute a _partial_ template once.
>
> Hugo renders pages in parallel, and will render the _partial_ template more than once with concurrent calls to the `partialCached` function. After Hugo caches the rendered _partial_ template, new pages entering the build pipeline will use the cached result.

Here is the simplest usage:

```go-html-template
{{ partialCached "footer.html" . }}
```

Pass additional arguments to `partialCached` to create variants of the cached _partial_ template. For example, if you have a complex _partial_ template that should be identical when rendered for pages within the same section, use a variant based on section so that the _partial_ template is only rendered once per section:

```go-html-template {file="layouts/baseof.html"}
{{ partialCached "footer.html" . .Section }}
```

Pass additional arguments, of any data type, as needed to create unique variants:

```go-html-template
{{ partialCached "footer.html" . .Params.country .Params.province }}
```

The variant arguments are not available to the underlying _partial_ template; they are only used to create unique cache keys.

To return a value from a _partial_ template, use the `return` statement:

```go-html-template
{{ if math.ModBool . 2 }}
  {{ return "even" }}
{{ end }}
{{ return "odd" }}
```

[`return`]: /functions/go-template/return/
