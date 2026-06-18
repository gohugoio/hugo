---
title: Keywords
description: Returns a slice of keywords as defined in front matter.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: '[]string'
    signatures: [PAGE.Keywords]
---

By default, Hugo evaluates the keywords when creating collections of [related content][].

{{< code-toggle file=content/recipes/sushi.md fm=true >}}
title = 'How to make spicy tuna hand rolls'
keywords = ['tuna','sriracha','nori','rice']
{{< /code-toggle >}}

To list the keywords within a template:

```go-html-template
{{ range .Keywords }}
  {{ . }}
{{ end }}
```

Or use the [`delimit`][] function:

```go-html-template
{{ delimit .Keywords ", " ", and " }} → tuna, sriracha, nori, and rice
```

Keywords are also a useful [taxonomy][]:

{{< code-toggle file=hugo >}}
[taxonomies]
tag = 'tags'
keyword = 'keywords'
category = 'categories'
{{< /code-toggle >}}

[`delimit`]: /functions/collections/delimit/
[related content]: /content-management/related-content/
[taxonomy]: /content-management/taxonomies/
