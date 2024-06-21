---
title: define
description: Defines a template.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/block
    - functions/go-template/end
    - functions/go-template/template
    - functions/partials/Include
    - functions/partials/IncludeCached
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

{{ define "partials/inline/foo.html" }}
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

[`block`]: /functions/go-template/block/
[`template`]: /functions/go-template/block/
[`partial`]: /functions/partials/include/

{{% include "functions/go-template/_common/text-template.md" %}}
