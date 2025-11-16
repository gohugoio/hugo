---
title: Page
description: Returns the Page object from which the shortcode was called.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: hugolib.pageForShortcode
    signatures: [SHORTCODE.Page]
---

With this content:

{{< code-toggle file=content/books/les-miserables.md fm=true >}}
title = 'Les Misérables'
author = 'Victor Hugo'
publication_year = 1862
isbn = '978-0451419439'
{{< /code-toggle >}}

Calling this shortcode:

```text
{{</* book-details */>}}
```

We can access the front matter values using the `Page` method:

```go-html-template {file="layouts/_shortcodes/book-details.html"}
<ul>
  <li>Title: {{ .Page.Title }}</li>
  <li>Author: {{ .Page.Params.author }}</li>
  <li>Published: {{ .Page.Params.publication_year }}</li>
  <li>ISBN: {{ .Page.Params.isbn }}</li>
</ul>
```
