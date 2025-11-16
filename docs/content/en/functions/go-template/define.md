---
title: define
description: Defines a template.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType:
    signatures: [define NAME]
---

Use with the [`block`] statement:

```go-html-template
{{ block "main" . }}
  {{ print "default value if 'main' template is empty" }}
{{ end }}

{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

Use with the [`partial`] function:

```go-html-template
{{ partial "inline/foo.html" (dict "answer" 42) }}

{{ define "_partials/inline/foo.html" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

Use with the [`template`] function:

```go-html-template
{{ template "foo" (dict "answer" 42) }}

{{ define "foo" }}
  {{ printf "The answer is %v." .answer }}
{{ end }}
```

{{% include "/_common/functions/go-template/text-template.md" %}}

[`block`]: /functions/go-template/block/
[`template`]: /functions/go-template/block/
[`partial`]: /functions/partials/include/
