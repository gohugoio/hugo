---
title: continue
description: Used with the range statement, stops the innermost iteration and continues to the next iteration.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/go-template/break
    - functions/go-template/range
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

{{% include "functions/go-template/_common/text-template.md" %}}
