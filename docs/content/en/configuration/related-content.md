---
title: Configure related content
linkTitle: Related content
description: Configure related content.
categories: []
keywords: []
---

> [!NOTE]
> Use the [`Pages.Related`][] method to render related content in your templates.

Hugo provides a sensible default configuration for identifying related content, but you can customize it in your project configuration, either globally or per language.

## Default configuration

This is the default configuration:

{{< code-toggle config=related />}}

> [!NOTE]
> Adding a `related` section to your project configuration requires you to provide a full configuration. You cannot override individual default values without specifying all related settings.

## Top-level settings

`threshold`
: (`int`) A value in the range `[0, 100]`. A lower value will return more, but maybe not so relevant, matches.

`includeNewer`
: (`bool`) Whether to include pages newer than the current page in the related content listing. The output for older posts may change as new related content is added. Default is `false`.

`toLower`
: (`bool`) Whether to transform keywords in both the indexes and the queries to lower case. This may give more accurate results at a slight performance penalty. Default is `false`.

## Per-index settings

`applyFilter`
: (`bool`) Apply a `type` specific filter to the result of a search. This is only used for the `fragments` type. Default is `false`.

`cardinalityThreshold`
: (`int`) The percentage threshold, in the range `[1, 100]`, above which indexed values are removed from the index. For example, a value of `60` removes all indexed values that appear in more than 60% of the documents. A value of `0` disables filtering. Default is `0`.

`minTokenLength`
: {{< new-in 0.166.0 />}}
: (`int`) When [`tokenize`](#tokenize) is `true`, the minimum token length in Unicode [code points](g). Use this to exclude short words such as `a`, `in`, and `the`, regardless of their frequency. This setting is applied before [`cardinalityThreshold`](#cardinalitythreshold), so excluded tokens never enter the index. Default is `0`, meaning no minimum.

`name`
: (`string`) The index name. This value maps directly to a page parameter. Hugo supports string values such as `author`, lists such as `tags` and `keywords`, and time and date objects.

`pattern`
: (`string`) This is only relevant for dates. When listing related content, you may want to list content that is also close in time. Setting `2006`, the default value for date indexes, as the pattern for a date index will add weight to pages published in the same year. For busier blogs, `200601`, representing year and month, may be a better default.

`tokenize`
: {{< new-in 0.166.0 />}}
: (`bool`) Whether to tokenize string values by splitting on whitespace before indexing.

  For indices of type `fragments`, non-heading fragments such as [description list terms][] are always indexed by their exact identifier.

  Tokenizing increases index size in proportion to the number of words per value. Consider setting [`minTokenLength`](#mintokenlength) to exclude short words and [`cardinalityThreshold`](#cardinalitythreshold) to remove high-frequency words from the index.

  Singular and plural forms such as `apple` and `apples` are treated as distinct words, and [CJK](g) languages and other scripts without whitespace word boundaries are not supported.

  Default is `false`.

`toLower`
: (`bool`) Whether to transform keywords in both the indexes and the queries to lower case. This may give more accurate results at a slight performance penalty. Default is `false`.

`type`
: (`string`) One of `basic` or `fragments`. Default is `basic`.

`weight`
: (`int`) An integer weight that indicates how important this parameter is relative to the other parameters. It can be `0`, which has the effect of turning this index off, or even negative. Test with different values to see what fits your content best. Default is `0`.

## Examples

The following examples demonstrate common configuration patterns.

### Related by taxonomy terms

Imagine we're building a book review site. Our main content will be book reviews, and we'll use genres and authors as taxonomies. When someone views a book review, we want to show a short list of related reviews based on shared authors and genres.

Create the content:

```tree
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

Then render a list of 5 related reviews with a _partial_ template like this:

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

### Related by title words

Imagine we're building a blog. We want to show related posts based on shared words in the title. For example, a post titled "Hugo image processing" should match any post whose title contains `Hugo`, `image`, or `processing`.

Configure the related content identification:

{{< code-toggle file=hugo >}}
[related]
threshold = 50
toLower = true
[[related.indices]]
name = 'title'
weight = 1
tokenize = true
minTokenLength = 4
{{< /code-toggle >}}

With `tokenize` set to `true`, Hugo splits each page title on whitespace and indexes the individual words. Setting `minTokenLength` to `4` excludes short words such as `and`, `the`, and `for`. Then use the same partial template shown in the previous example to render the results.

[`Pages.Related`]: /methods/pages/related/
[description list terms]: /configuration/markup/#parserautodefinitiontermid
