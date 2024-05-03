---
title: Description
description: Returns the description of the given page as defined in front matter.
categories: []
keywords: []
action:
  related:
    - methods/page/Summary
  returnType: string
  signatures: [PAGE.Description]
---

Conceptually different from a [content summary], a page description is typically used in metadata about the page.

{{< code-toggle file=content/recipes/sushi.md fm=true >}}
title = 'How to make spicy tuna hand rolls'
description = 'Instructions for making spicy tuna hand rolls.'
{{< /code-toggle >}}

{{< code file=layouts/baseof.html  >}}
<head>
  ...
  <meta name="description" content="{{ .Description }}">
  ...
</head>
{{< /code >}}

[content summary]: /content-management/summaries/
