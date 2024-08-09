---
title: template
description: Executes the given template, optionally passing context.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/define
    - functions/partials/Include
    - functions/partials/IncludeCached
  returnType: 
  signatures: ['template NAME [CONTEXT]']
---

Use the `template` function to execute [embedded templates]. For example:

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

{{% include "functions/go-template/_common/text-template.md" %}}

[`partial`]: /functions/partials/include/
[inline partial]: /templates/partial/#inline-partials
[embedded templates]: /templates/embedded/
