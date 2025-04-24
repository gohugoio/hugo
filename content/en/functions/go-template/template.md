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

Use the `template` function to execute any of these [embedded templates](g):

- [`disqus.html`]
- [`google_analytics.html`]
- [`opengraph.html`]
- [`pagination.html`]
- [`schema.html`]
- [`twitter_cards.html`]



For example:

```go-html-template
{{ range (.Paginate .Pages).Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{ template "_internal/pagination.html" . }}
```

You can also use the `template` function to execute a defined template:

```go-html-template
{{ template "foo" (dict "answer" 42) }}

{{ define "foo" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

The example above can be rewritten using an [inline partial] template:

```go-html-template
{{ partial "inline/foo.html" (dict "answer" 42) }}

{{ define "partials/inline/foo.html" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

The key distinctions between the preceding two examples are:

1. Inline partials are globally scoped. That means that an inline partial defined in _one_ template may be called from _any_ template.
2. Leveraging the [`partialCached`] function when calling an inline partial allows for performance optimization through result caching.
3. An inline partial can [`return`] a value of any data type instead of rendering a string.

{{% include "/_common/functions/go-template/text-template.md" %}}

[`disqus.html`]: /templates/embedded/#disqus
[`google_analytics.html`]: /templates/embedded/#google-analytics
[`opengraph.html`]: /templates/embedded/#open-graph
[`pagination.html`]: /templates/embedded/#pagination
[`partialCached`]: /functions/partials/includecached/
[`partial`]: /functions/partials/include/
[`return`]: /functions/go-template/return/
[`schema.html`]: /templates/embedded/#schema
[`twitter_cards.html`]: /templates/embedded/#x-twitter-cards
[inline partial]: /templates/partial/#inline-partials
