---
title: partials.IncludeCached
description: Executes the given template and caches the result, optionally passing context. If the partial template contains a return statement, returns the given value, else returns the rendered output.
categories: []
keywords: []
action:
  aliases: [partialCached]
  related:
    - functions/go-template/return
    - functions/partials/Include
    - functions/go-template/template
    - methods/page/Render
  returnType: any
  signatures: ['partials.IncludeCached LAYOUT CONTEXT [VARIANT...]']
signatures: 
  - partials.IncludeCached NAME CONTEXT [VARIANT...]
  - partialCached NAME CONTEXT [VARIANT...]
aliases: [/functions/partialcached]
---

Without a [`return`] statement, the `partialCached` function returns a string of type `template.HTML`. With a `return` statement, the `partialCached` function can return any data type.

The `partialCached` function can offer significant performance gains for complex templates that don't need to be re-rendered on every invocation.

{{% note %}}
Each Site (or language) has its own `partialCached` cache, so each site will execute a partial once.

Hugo renders pages in parallel, and will render the partial more than once with concurrent calls to the `partialCached` function. After Hugo caches the rendered partial, new pages entering the build pipeline will use the cached result.
{{% /note %}}

Here is the simplest usage:

```go-html-template
{{ partialCached "footer.html" . }}
```

Pass additional arguments to `partialCached` to create variants of the cached partial. For example, if you have a complex partial that should be identical when rendered for pages within the same section, use a variant based on section so that the partial is only rendered once per section:

{{< code file=partial-cached-example.html >}}
{{ partialCached "footer.html" . .Section }}
{{< /code >}}

Pass additional arguments, of any data type, as needed to create unique variants:

```go-html-template
{{ partialCached "footer.html" . .Params.country .Params.province }}
```

The variant arguments are not available to the underlying partial template; they are only used to create unique cache keys. 

To return a value from a partial template, it must contain only one `return` statement, placed at the end of the template:

```go-html-template
{{ $result := "" }}
{{ if math.ModBool . 2 }}
  {{ $result = "even" }}
{{ else }}
  {{ $result = "odd" }}
{{ end }}
{{ return $result }}
```

See&nbsp;[details][`return`].

[`return`]: /functions/go-template/return/
