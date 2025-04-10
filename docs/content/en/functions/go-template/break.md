---
title: break
description: Used with the range statement, stops the innermost iteration and bypasses all remaining iterations.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType:
    signatures: [break]
---

This template code:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ if eq . "bar" }}
    {{ break }}
  {{ end }}
  <p>{{ . }}</p>
{{ end }}
```

Is rendered to:

```html
<p>foo</p>
```

{{% include "/_common/functions/go-template/text-template.md" %}}
