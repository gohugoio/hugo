---
title: Page Variables
linktitle:
description: Page-level variables are defined in a content file's front matter, derived from the content's file location, or extracted from the content body itself.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
keywords: [pages]
draft: false
menu:
  docs:
    title: "variables defined by a page"
    parent: "variables"
    weight: 20
weight: 20
sections_weight: 20
aliases: []
toc: true
---

The following is a list of page-level variables. Many of these will be defined in the front matter, derived from file location, or extracted from the content itself.

{{% note "`.Scratch`" %}}
See [`.Scratch`](/functions/scratch/) for page-scoped, writable variables.
{{% /note %}}

## Page Variables

.AlternativeOutputFormats
: contains all alternative formats for a given page; this variable is especially useful `link rel` list in your site's `<head>`. (See [Output Formats](/templates/output-formats/).)

.Aliases
: aliases of this page

.Content
: the content itself, defined below the front matter.

.Data
: the data specific to this type of page.

.Date
: the date associated with the page; `.Date` pulls from the `date` field in a content's front matter. See also `.ExpiryDate`, `.PublishDate`, and `.Lastmod`.

.Description
: the description for the page.

.Dir
: the path of the folder containing this content file. The path is relative to the `content` folder.

.Draft
: a boolean, `true` if the content is marked as a draft in the front matter.

.ExpiryDate
: the date on which the content is scheduled to expire; `.ExpiryDate` pulls from the `expirydate` field in a content's front matter. See also `.PublishDate`, `.Date`, and `.Lastmod`.

.File
: filesystem-related data for this content file. See also [File Variables][].

.FuzzyWordCount
: the approximate number of words in the content.

.Hugo
: see [Hugo Variables](/variables/hugo/).

.IsHome
: `true` in the context of the [homepage](/templates/homepage/).

.IsNode
: always `false` for regular content pages.

.IsPage
: always `true` for regular content pages.

.IsTranslated
: `true` if there are translations to display.

.Keywords
: the meta keywords for the content.

.Kind
: the page's *kind*. Possible return values are `page`, `home`, `section`, `taxonomy`, or `taxonomyTerm`. Note that there are also `RSS`, `sitemap`, `robotsTXT`, and `404` kinds, but these are only available during the rendering of each of these respective page's kind and therefore *not* available in any of the `Pages` collections.

.Language
: a language object that points to the language's definition in the site `config`. `.Language.Lang` gives you the language code.

.Lastmod
: the date the content was last modified. `.Lastmod` pulls from the `lastmod` field in a content's front matter.

 - If `lastmod` is not set, and `.GitInfo` feature is disabled, the front matter `date` field will be used.
 - If `lastmod` is not set, and `.GitInfo` feature is enabled, `.GitInfo.AuthorDate` will be used instead.

See also `.ExpiryDate`, `.Date`, `.PublishDate`, and [`.GitInfo`][gitinfo].

.LinkTitle
: access when creating links to the content. If set, Hugo will use the `linktitle` from the front matter before `title`.

.Next
: Points up to the next [regular page](/variables/site/#site-pages) (sorted by Hugo's [default sort](/templates/lists#default-weight-date-linktitle-filepath)). Example: `{{with .Next}}{{.Permalink}}{{end}}`. Calling `.Next` from the first page returns `nil`.

.NextInSection
: Points up to the next [regular page](/variables/site/#site-pages) below the same top level section (e.g. in `/blog`)). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight-date-linktitle-filepath). Example: `{{with .NextInSection}}{{.Permalink}}{{end}}`. Calling `.NextInSection` from the first page returns `nil`.

.OutputFormats
: contains all formats, including the current format, for a given page. Can be combined the with [`.Get` function](/functions/get/) to grab a specific format. (See [Output Formats](/templates/output-formats/).)

.Pages
: a collection of associated pages. This value will be `nil` within
  the context of regular content pages. See [`.Pages`](#pages).

.Permalink
: the Permanent link for this page; see [Permalinks](/content-management/urls/)

.Plain
: the Page content stripped of HTML tags and presented as a string.

.PlainWords
: the Page content stripped of HTML as a `[]string` using Go's [`strings.Fields`](https://golang.org/pkg/strings/#Fields) to split `.Plain` into a slice.

.Prev
: Points down to the previous [regular page](/variables/site/#site-pages) (sorted by Hugo's [default sort](/templates/lists#default-weight-date-linktitle-filepath)). Example: `{{if .PrevPage}}{{.PrevPage.Permalink}}{{end}}`.  Calling `.Prev` from the last page returns `nil`.

.PrevInSection
: Points down to the previous [regular page](/variables/site/#site-pages) below the same top level section (e.g. `/blog`). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight-date-linktitle-filepath). Example: `{{if .PrevInSection}}{{.PrevInSection.Permalink}}{{end}}`.  Calling `.PrevInSection` from the last page returns `nil`.

.PublishDate
: the date on which the content was or will be published; `.Publishdate` pulls from the `publishdate` field in a content's front matter. See also `.ExpiryDate`, `.Date`, and `.Lastmod`.

.RSSLink (deprecated)
: link to the page's RSS feed. This is deprecated. You should instead do something like this: `{{ with .OutputFormats.Get "RSS" }}{{ .RelPermalink }}{{ end }}`.

.RawContent
: raw markdown content without the front matter. Useful with [remarkjs.com](
http://remarkjs.com)

.ReadingTime
: the estimated time, in minutes, it takes to read the content.

.Resources
: resources such as images and CSS that are associated with this page

.Ref
: returns the permalink for a given reference (e.g., `.Ref "sample.md"`).  `.Ref` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

.RelPermalink
: the relative permanent link for this page.

.RelRef
: returns the relative permalink for a given reference (e.g., `RelRef
"sample.md"`). `.RelRef` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

.Site
: see [Site Variables](/variables/site/).

.Sites
: returns all sites (languages). A typical use case would be to link back to the main language: `<a href="{{ .Sites.First.Home.RelPermalink }}">...</a>`.

.Sites.First
: returns the site for the first language. If this is not a multilingual setup, it will return itself.

.Summary
: a generated summary of the content for easily showing a snippet in a summary view. The breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page, or the summary can be written independent of the page text.  See [Content Summaries](/content-management/summaries/) for more details.

.TableOfContents
: the rendered [table of contents](/content-management/toc/) for the page.

.Title
: the title for this page.

.Translations
: a list of translated versions of the current page. See [Multilingual Mode](/content-management/multilingual/) for more information.

.TranslationKey
: the key used to map language translations of the current page. See [Multilingual Mode](/content-management/multilingual/) for more information.

.Truncated
: a boolean, `true` if the `.Summary` is truncated. Useful for showing a "Read more..." link only when necessary.  See [Summaries](/content-management/summaries/) for more information.

.Type
: the [content type](/content-management/types/) of the content (e.g., `posts`).

.UniqueID
: the MD5-checksum of the content file's path.

.Weight
: assigned weight (in the front matter) to this content, used in sorting.

.WordCount
: the number of words in the content.

## Section Variables and Methods

Also see [Sections](/content-management/sections/).

{{< readfile file="/content/en/readfiles/sectionvars.md" markdown="true" >}}

## The `.Pages` Variable {#pages}

`.Pages` is an alias to `.Data.Pages`. It is conventional to use the
aliased form `.Pages`.

### `.Pages` compared to `.Site.Pages`

{{< readfile file="/content/en/readfiles/pages-vs-site-pages.md" markdown="true" >}}

## Page-level Params

Any other value defined in the front matter in a content file, including taxonomies, will be made available as part of the `.Params` variable.

```
---
title: My First Post
date: 2017-02-20T15:26:23-06:00
categories: [one]
tags: [two,three,four]
```

With the above front matter, the `tags` and `categories` taxonomies are accessible via the following:

* `.Params.tags`
* `.Params.categories`

{{% note "Casing of Params" %}}
Page-level `.Params` are *only* accessible in lowercase.
{{% /note %}}

The `.Params` variable is particularly useful for the introduction of user-defined front matter fields in content files. For example, a Hugo website on book reviews could have the following front matter in `/content/review/book01.md`:

```
---
...
affiliatelink: "http://www.my-book-link.here"
recommendedby: "My Mother"
...
---
```

These fields would then be accessible to the `/themes/yourtheme/layouts/review/single.html` template through `.Params.affiliatelink` and `.Params.recommendedby`, respectively.

Two common situations where this type of front matter field could be introduced is as a value of a certain attribute like `href=""` or by itself to be displayed as text to the website's visitors.

{{< code file="/themes/yourtheme/layouts/review/single.html" >}}
<h3><a href={{ printf "%s" $.Params.affiliatelink }}>Buy this book</a></h3>
<p>It was recommended by {{ .Params.recommendedby }}.</p>
{{< /code >}}

This template would render as follows, assuming you've set [`uglyURLs`](/content-management/urls/) to `false` in your [site `config`](/getting-started/configuration/):

{{< output file="yourbaseurl/review/book01/index.html" >}}
<h3><a href="http://www.my-book-link.here">Buy this book</a></h3>
<p>It was recommended by my Mother.</p>
{{< /output >}}

{{% note %}}
See [Archetypes](/content-management/archetypes/) for consistency of `Params` across pieces of content.
{{% /note %}}

### The `.Param` Method

In Hugo, you can declare params in individual pages and globally for your entire website. A common use case is to have a general value for the site param and a more specific value for some of the pages (i.e., a header image):

```
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value according to it's definition in a page parameter (i.e. in the content's front matter) or a site parameter (i.e., in your `config`).

### Access Nested Fields in Front Matter

When front matter contains nested fields like the following:

```
---
author:
  given_name: John
  family_name: Feminella
  display_name: John Feminella
---
```
`.Param` can access these fields by concatenating the field names together with a dot:

```
{{ $.Param "author.display_name" }}
```

If your front matter contains a top-level key that is ambiguous with a nested key, as in the following case:

```
---
favorites.flavor: vanilla
favorites:
  flavor: chocolate
---
```

The top-level key will be preferred. Therefore, the following method, when applied to the previous example, will print `vanilla` and not `chocolate`:

```
{{ $.Param "favorites.flavor" }}
=> vanilla
```

[gitinfo]: /variables/git/
[File Variables]: /variables/files/
