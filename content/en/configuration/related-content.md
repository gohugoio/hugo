---
title: Configure related content
linkTitle: Related content
description: Configure related content.
categories: []
keywords: []
---

> [!note]
> To understand Hugo's related content identification, please refer to the [related content] page.

Hugo provides a sensible default configuration for identifying related content, but you can customize it in your site configuration, either globally or per language.

## Default configuration

This is the default configuration:

{{< code-toggle config=related />}}

> [!note]
> Adding a `related` section to your site configuration requires you to provide a full configuration. You cannot override individual default values without specifying all related settings.

## Top-level options

threshold
: (`int`) A value between 0-100, inclusive. A lower value will return more, but maybe not so relevant, matches.

includeNewer
: (`bool`) Whether to include pages newer than the current page in the related content listing. This will mean that the output for older posts may change as new related content gets added. Default is `false`.

toLower
: (`bool`) Whether to transform keywords in both the indexes and the queries to lower case. This may give more accurate results at a slight performance penalty. Default is `false`.

## Per-index options

name
: (`string`) The index name. This value maps directly to a page parameter. Hugo supports string values (`author` in the example) and lists (`tags`, `keywords` etc.) and time and date objects.

type
: (`string`) One of `basic` or `fragments`. Default is `basic`.

applyFilter
: (`string`) Apply a `type` specific filter to the result of a search. This is currently only used for the `fragments` type.

weight
: (`int`) An integer weight that indicates how important this parameter is relative to the other parameters. It can be `0`, which has the effect of turning this index off, or even negative. Test with different values to see what fits your content best. Default is `0`.

cardinalityThreshold
: (`int`) If between 1 and 100, this is a percentage. All keywords that are used in more than this percentage of documents are removed. For example, setting this to `60` will remove all keywords that are used in more than 60% of the documents in the index. If `0`, no keyword is removed from the index. Default is `0`.

pattern
: (`string`) This is currently only relevant for dates. When listing related content, we may want to list content that is also close in time. Setting "2006" (default value for date indexes) as the pattern for a date index will add weight to pages published in the same year. For busier blogs, "200601" (year and month) may be a better default.

toLower
: (`bool`) Whether to transform keywords in both the indexes and the queries to lower case. This may give more accurate results at a slight performance penalty. Default is `false`.

## Example

Imagine we're building a book review site. Our main content will be book reviews, and we'll use genres and authors as taxonomies. When someone views a book review, we want to show a short list of related reviews based on shared authors and genres.

Create the content:

```text
content/
└── book-reviews/
    ├── book-review-1.md
    ├── book-review-2.md
    ├── book-review-3.md
    ├── book-review-4.md
    └── book-review-5.md
```

Configure the taxonomies:

{{< code-toggle file=hugo >}}
[taxonomies]
author = 'authors'
genre = 'genres'
{{< /code-toggle >}}

Configure the related content identification:

{{< code-toggle file=hugo >}}
[related]
includeNewer = true
threshold = 80
toLower = true
[[related.indices]]
name = 'authors'
weight = 2
[[related.indices]]
name = 'genres'
weight = 1
{{< /code-toggle >}}

We've configured the `authors` index with a weight of `2` and the `genres` index with a weight of `1`. This means Hugo prioritizes shared `authors` as twice as significant as shared `genres`.

Then render a list of 5 related reviews with a partial template like this:

```go-html-template {file="layouts/_partials/related.html" copy=true}
{{ with site.RegularPages.Related . | first 5 }}
  <p>Related content:</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

[related content]: /content-management/related-content/
