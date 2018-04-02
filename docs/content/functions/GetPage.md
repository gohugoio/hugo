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
signature: [".GetPage KIND PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Every `Page` has a [`Kind` attribute][page_kinds] that shows what kind of page it is. While this attribute can be used to list pages of a certain `kind` using `where`, often it can be useful to fetch a single page by its path.

`.GetPage` returns a page of a given `Kind` and `path`.

{{% note %}}
If the `path` is `"foo/bar.md"`, it can be written as exactly that, or broken up
into multiple strings as `"foo" "bar.md"`.
{{% /note %}}

```
{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
```

This method wil return `nil` when no page could be found, so the above will not print anything if the blog section is not found.

For a regular page (whose `Kind` is `page`):

```
{{ with .Site.GetPage "page" "blog/my-post.md" }}{{ .Title }}{{ end }}
```

Note that the `path` can also be supplied like this, where the slash-separated
path elements are added as separate strings:

```
{{ with .Site.GetPage "page" "blog" "my-post.md" }}{{ .Title }}{{ end }}
```

## `.GetPage` Example

This code snippet---in the form of a [partial template][partials]---allows you to do the following:

1. Grab the index object of your `tags` [taxonomy][].
2. Assign this object to a variable, `$t`
3. Sort the terms associated with the taxonomy by popularity.
4. Grab the top two most popular terms in the taxonomy (i.e., the two most popular tags assigned to content.

{{< code file="grab-top-two-tags.html" >}}
<ul class="most-popular-tags">
{{ $t := .Site.GetPage "taxonomyTerm" "tags" }}
{{ range first 2 $t.Data.Terms.ByCount }}
    <li>{{ . }}</li>
{{ end }}
</ul>
{{< /code >}}

## `.GetPage` on Page Bundles

If the page retrieved by `.GetPage` is a [Leaf Bundle][leaf_bundle], and you
need to get the nested _**page** resources_ in that, you will need to use the
methods in `.Resources` as explained in the [Page Resources][page_resources]
section.

See the [Headless Bundle][headless_bundle] documentation for an example.


[partials]: /templates/partials/
[taxonomy]: /content-management/taxonomies/
[page_kinds]: /templates/section-templates/#page-kinds
[leaf_bundle]: /content-management/page-bundles/#leaf-bundles
[headless_bundle]: /content-management/page-bundles/#headless-bundle
[page_resources]: /content-management/page-resources/
