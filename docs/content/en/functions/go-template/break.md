---
title: break
description: Used with the range statement, stops the innermost iteration and bypasses all remaining iterations.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/continue
    - functions/go-template/range
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

{{% include "functions/go-template/_common/text-template.md" %}}
