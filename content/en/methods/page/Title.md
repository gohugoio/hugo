---
title: Title
description: Returns the title of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/LinkTitle
  returnType: string
  signatures: [PAGE.Title]
---

For pages backed by a file, the `Title` method returns the `title` field in front matter. For section pages that are automatically generated, the `Title` method returns the section name.

{{< code-toggle file=content/about.md fm=true >}}
title = 'About us'
{{< /code-toggle >}}

```go-html-template
{{ .Title }} → About us
```

With automatic section pages, the title is capitalized by default, using the `titleCaseStyle` defined in your site configuration. To disable this behavior, set `pluralizeListTitles` to `false` in your site configuration.
