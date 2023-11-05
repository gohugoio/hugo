---
title: Page
description: Returns the Page object from which the shortcode was called.
categories: []
keywords: []
action:
  related: []
  returnType: hugolib.pageForShortcode
  signatures: [SHORTCODE.Page]
---

With this content:

{{< code-toggle file=content/books/les-miserables.md copy=false fm=true >}}
title = 'Les Misérables'
author = 'Victor Hugo'
published = 1862
isbn = '978-0451419439'
{{< /code-toggle >}}

Calling this shortcode:

```text
{{</* book-details */>}}
```

We can access the front matter values using the `Page` method:

{{< code file="layouts/shortcodes/book-details.html" lang=go-html-template copy=false >}}
<ul>
  <li>Title: {{ .Page.Title }}</li>
  <li>Author: {{ .Page.Params.author }}</li>
  <li>Published: {{ .Page.Params.publication_year }}</li>
  <li>ISBN: {{ .Page.Params.isbn }}</li>
</ul>
{{< /code >}}
