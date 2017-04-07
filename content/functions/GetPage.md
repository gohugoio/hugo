---
title: getpage
linktitle: GetPage
description: Looks up the index page (`_index.md`) of a given `Kind` and `path`.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [sections,lists,indexes]
ns:
signature: ["GetPage TYPE PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Every `Page` has a `Kind` attribute that shows what kind of page it is. While this attribute can be used to list pages of a certain `kind` using `where`, often it can be useful to fetch a single page by its path.

`GetPage` looks up an index page of a given `Kind` and `path`. This method may support regular pages in the future, but currently it is a convenient way of getting the index pages, such as the homepage or a section, from a template:

```
{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
```

This method wil return `nil` when no page could be found, so the above will not print anything if the blog section isn't found.

The valid page kinds are: *home, section, taxonomy and taxonomyTerm.*

## `.GetPage` Example

This code snippet---in the form of a [partial template][partials]---allows you to do the following:

1. Grab the index object of your `tags` [taxonomy][].
2. Assign this object to a variable, `$t`
3. Sort the terms associated with the taxonomy by popularity.
4. Grab the top two most popular terms in the taxonomy (i.e., the two most popular tags assigned to content.

{{% code file="grab-top-two-tags.html" %}}
```html
<ul class="most-popular-tags">
{{ $t := $.Site.GetPage("taxonomyTerm", "tags") }}
{{ range first 2 $t.Data.Terms.ByCount }}
    <li>{{.}}</li>
{{ end }}
</ul>
```
{{% /code %}}


[partials]: /templates/partials/
[taxonomy]: /content-management/taxonomies/
