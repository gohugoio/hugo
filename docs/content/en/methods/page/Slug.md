---
title: Slug
description: Returns the URL slug of the given page as defined in front matter.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [PAGE.Slug]
---

{{< code-toggle file=content/recipes/spicy-tuna-hand-rolls.md fm=true >}}
title = 'How to make spicy tuna hand rolls'
slug = 'sushi'
{{< /code-toggle >}}

This page will be served from:

    https://example.org/recipes/sushi

To get the slug value within a template:

```go-html-template
{{ .Slug }} → sushi
```
