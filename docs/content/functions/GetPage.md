---
title: .GetPage
description: "Gets a `Page` of a given `Kind` and `path`."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [sections,lists,indexes]
signature: [".GetPage TYPE PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Every `Page` has a `Kind` attribute that shows what kind of page it is. While this attribute can be used to list pages of a certain `kind` using `where`, often it can be useful to fetch a single page by its path.

`.GetPage` looks up a page of a given `Kind` and `path`.

```
{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
```

This method wil return `nil` when no page could be found, so the above will not print anything if the blog section isn't found.

For a regular page:

```
{{ with .Site.GetPage "page" "blog" "my-post.md" }}{{ .Title }}{{ end }}
```

Note that the path can also be supplied like this:

```
{{ with .Site.GetPage "page" "blog/my-post.md" }}{{ .Title }}{{ end }}
```

The valid page kinds are: *page, home, section, taxonomy and taxonomyTerm.*

## `.GetPage` Example

This code snippet---in the form of a [partial template][partials]---allows you to do the following:

1. Grab the index object of your `tags` [taxonomy][].
2. Assign this object to a variable, `$t`
3. Sort the terms associated with the taxonomy by popularity.
4. Grab the top two most popular terms in the taxonomy (i.e., the two most popular tags assigned to content.

{{< code file="grab-top-two-tags.html" >}}
<ul class="most-popular-tags">
{{ $t := $.Site.GetPage "taxonomyTerm" "tags" }}
{{ range first 2 $t.Data.Terms.ByCount }}
    <li>{{.}}</li>
{{ end }}
</ul>
{{< /code >}}


[partials]: /templates/partials/
[taxonomy]: /content-management/taxonomies/
