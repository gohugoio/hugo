---
title: .GetPage
description: "Gets a `Page` of a given `path`."
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [sections,lists,indexes]
signature: [".GetPage PATH"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`.GetPage` returns a page of a given `path`. Both `Site` and `Page` implements this method. The `Page` variant will, if given a relative path -- i.e. a path without a leading `/` -- try look for the page relative to the current page.

{{% note %}}
**Note:** We overhauled and simplified the `.GetPage` API in Hugo 0.45. Before that you needed to provide a `Kind` attribute in addition to the path, e.g. `{{ .Site.GetPage "section" "blog" }}`. This will still work, but is now superflous.
{{% /note %}}


```go-html-template
{{ with .Site.GetPage "/blog" }}{{ .Title }}{{ end }}
```

This method wil return `nil` when no page could be found, so the above will not print anything if the blog section is not found.

To find a regular page in the blog section::

```go-html-template
{{ with .Site.GetPage "/blog/my-post.md" }}{{ .Title }}{{ end }}
```

And since `Page` also provides a `.GetPage` method, the above is the same as:

```go-html-template
{{ with .Site.GetPage "/blog" }}
{{ with .GetPage "my-post.md" }}{{ .Title }}{{ end }}
{{ end }}
```

## .GetPage and Multilingual Sites

The previous examples have used the full content filename to lookup the post. Depending on how you have organized your content (whether you have the language code in the file name or not, e.g. `my-post.en.md`), you may want to do the lookup without extension. This will get you the current language's version of the page:

```go-html-template
{{ with .Site.GetPage "/blog/my-post" }}{{ .Title }}{{ end }}
```

## .GetPage Example

This code snippet---in the form of a [partial template][partials]---allows you to do the following:

1. Grab the index object of your `tags` [taxonomy][].
2. Assign this object to a variable, `$t`
3. Sort the terms associated with the taxonomy by popularity.
4. Grab the top two most popular terms in the taxonomy (i.e., the two most popular tags assigned to content.

{{< code file="grab-top-two-tags.html" >}}
<ul class="most-popular-tags">
{{ $t := .Site.GetPage "/tags" }}
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
