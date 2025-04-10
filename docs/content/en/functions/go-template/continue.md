---
title: continue
description: Used with the range statement, stops the innermost iteration and continues to the next iteration.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType:
    signatures: [continue]
---

This template code:

```go-html-template
{{ $s := slice "foo" "bar" "baz" }}
{{ range $s }}
  {{ if eq . "bar" }}
    {{ continue }}
  {{ end }}
  <p>{{ . }}</p>
{{ end }}
```

Is rendered to:

```html
<p>foo</p>
<p>baz</p>
```

{{% include "/_common/functions/go-template/text-template.md" %}}
