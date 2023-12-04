---
title: block
description: Defines a template and executes it in place.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/define
    - functions/go-template/end
  returnType:
  signatures: [block NAME CONTEXT]
---

A block is shorthand for defining a template:

```go-html-template
{{ define "name" }} T1 {{ end }}
```

and then executing it in place:

```go-html-template
{{ template "name" pipeline }}
```
The typical use is to define a set of root templates that are then customized by redefining the block templates within.

{{< code file=layouts/_default/baseof.html >}}
<body>
  <main>
    {{ block "main" . }}
      {{ print "default value if 'main' template is empty" }}
    {{ end }}
  </main>
</body>
{{< /code >}}

{{< code file=layouts/_default/single.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
{{< /code >}}

{{< code file=layouts/_default/list.html >}}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
{{< /code >}}

{{% include "functions/go-template/_common/text-template.md" %}}
