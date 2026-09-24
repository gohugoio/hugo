---
title: template
description: Executes the given template, optionally passing context.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType:
    signatures: ['template NAME [CONTEXT]']
---

Use the `template` statement to execute a defined template:

```go-html-template
{{ template "foo" (dict "answer" 42) }}

{{ define "foo" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

The example above can be rewritten using an inline _partial_ template:

```go-html-template
{{ partial "inline/foo.html" (dict "answer" 42) }}

{{ define "_partials/inline/foo.html" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

The key distinctions between the preceding two examples are:

1. Inline _partial_ templates are globally scoped. That means that an inline _partial_ template defined in one template may be called from any template.
1. Leveraging the [`partialCached`][] function when calling an inline _partial_ template allows for performance optimization through result caching.
1. An inline _partial_ template can [`return`][] a value of any data type instead of rendering a string.

{{% include "/_common/functions/go-template/text-template.md" %}}

[`partialCached`]: /functions/partials/includecached/
[`return`]: /functions/go-template/return/
