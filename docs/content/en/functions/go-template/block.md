---
title: block
description: Defines a template and executes it in place.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
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

```go-html-template {file="layouts/baseof.html"}
<body>
  <main>
    {{ block "main" . }}
      {{ print "default value if 'main' template is empty" }}
    {{ end }}
  </main>
</body>
```

```go-html-template {file="layouts/page.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
{{ end }}
```

```go-html-template {file="layouts/section.html"}
{{ define "main" }}
  <h1>{{ .Title }}</h1>
  {{ .Content }}
  {{ range .Pages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
```

{{% include "/_common/functions/go-template/text-template.md" %}}
