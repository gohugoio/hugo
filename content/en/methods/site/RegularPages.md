---
title: RegularPages
description: Returns a collection of all regular pages.
categories: []
keywords: []
action:
  related:
    - methods/site/AllPages
    - methods/site/RegularPages
    - methods/site/Sections
  returnType: page.Pages
  signatures: [SITE.RegularPages]
---

The `RegularPages` method on a `Site` object returns a collection of all [regular pages].

[regular pages]: /getting-started/glossary/#regular-page

```go-html-template
{{ range .Site.RegularPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
```

By default, Hugo sorts page collections by:

1. The page `weight` as defined in front matter
1. The page `date` as defined in front matter
1. The page `linkTitle` as defined in front matter
1. The file path

If the `linkTitle` is not defined, Hugo evaluates the `title` instead.

To change the sort order, use any of the `Pages` [sorting methods]. For example:

```go-html-template
{{ range .Site.RegularPages.ByTitle }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
```

[sorting methods]: /methods/pages/
