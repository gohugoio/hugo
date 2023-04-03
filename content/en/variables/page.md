---
title: Page Variables
description: Page-level variables are defined in a content file's front matter, derived from the content's file location, or extracted from the content body itself.
categories: [variables and params]
keywords: [pages]
menu:
  docs:
    parent: variables
    weight: 20
weight: 20
toc: true
---

The following is a list of page-level variables. Many of these will be defined in the front matter, derived from file location, or extracted from the content itself.

## Page Variables

.AlternativeOutputFormats
: contains all alternative formats for a given page; this variable is especially useful `link rel` list in your site's `<head>`. (See [Output Formats](/templates/output-formats/).)

.Aliases
: aliases of this page

.Ancestors
: get the ancestors of each page, simplify [breadcrumb navigation](/content-management/sections#example-breadcrumb-navigation) implementation complexity  

.BundleType
: the [bundle] type: `leaf`, `branch`, or an empty string if the page is not a bundle.

.Content
: the content itself, defined below the front matter.

.Data
: the data specific to this type of page.

.Date
: the date associated with the page; `.Date` pulls from the `date` field in a content's front matter. See also `.ExpiryDate`, `.PublishDate`, and `.Lastmod`.

.Description
: the description for the page.

.Draft
: a boolean, `true` if the content is marked as a draft in the front matter.

.ExpiryDate
: the date on which the content is scheduled to expire; `.ExpiryDate` pulls from the `expirydate` field in a content's front matter. See also `.PublishDate`, `.Date`, and `.Lastmod`.

.File
: filesystem-related data for this content file. See also [File Variables].

.Fragments
: Fragments returns the fragments for this page. See [Page Fragments](#page-fragments).

.FuzzyWordCount
: the approximate number of words in the content.

.IsHome
: `true` in the context of the [homepage](/templates/homepage/).

.IsNode
: always `false` for regular content pages.

.IsPage
: always `true` for regular content pages.

.IsSection
: `true` if [`.Kind`](/templates/section-templates/#page-kinds) is `section`.

.IsTranslated
: `true` if there are translations to display.

.Keywords
: the meta keywords for the content.

.Kind
: the page's *kind*. Possible return values are `page`, `home`, `section`, `taxonomy`, or `term`. Note that there are also `RSS`, `sitemap`, `robotsTXT`, and `404` kinds, but these are only available during the rendering of each of these respective page's kind and therefore *not* available in any of the `Pages` collections.

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
: Points up to the next [regular page](/variables/site/#site-pages) (sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath)). Example: `{{ with .Next }}{{ .Permalink }}{{ end }}`. Calling `.Next` from the first page returns `nil`.

.NextInSection
: Points up to the next [regular page](/variables/site/#site-pages) below the same top level section (e.g. in `/blog`)). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath). Example: `{{ with .NextInSection }}{{ .Permalink }}{{ end }}`. Calling `.NextInSection` from the first page returns `nil`.

.OutputFormats
: contains all formats, including the current format, for a given page. Can be combined the with [`.Get` function](/functions/get/) to grab a specific format. (See [Output Formats](/templates/output-formats/).)

.Pages
: a collection of associated pages. This value will be `nil` within
  the context of regular content pages. See [`.Pages`](#pages).

.Permalink
: the Permanent link for this page; see [Permalinks](/content-management/urls/)

.Plain
: the Page content stripped of HTML tags and presented as a string. You may need to pipe the result through the [`htmlUnescape`](/functions/htmlunescape/) function when rendering this value with the HTML [output format](/templates/output-formats#output-format-definitions).

.PlainWords
: the slice of strings that results from splitting .Plain into words, as defined in Go's [strings.Fields](https://pkg.go.dev/strings#Fields).

.Prev
: Points down to the previous [regular page](/variables/site/#site-pages) (sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath)). Example: `{{ if .Prev }}{{ .Prev.Permalink }}{{ end }}`.  Calling `.Prev` from the last page returns `nil`.

.PrevInSection
: Points down to the previous [regular page](/variables/site/#site-pages) below the same top level section (e.g. `/blog`). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath). Example: `{{ if .PrevInSection }}{{ .PrevInSection.Permalink }}{{ end }}`.  Calling `.PrevInSection` from the last page returns `nil`.

.PublishDate
: the date on which the content was or will be published; `.Publishdate` pulls from the `publishdate` field in a content's front matter. See also `.ExpiryDate`, `.Date`, and `.Lastmod`.

.RawContent
: raw markdown content without the front matter. Useful with [remarkjs.com](
https://remarkjs.com)

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

.Weight
: assigned weight (in the front matter) to this content, used in sorting.

.WordCount
: the number of words in the content.

## Writable Page-scoped Variables

[.Scratch][scratch]
: returns a Scratch to store and manipulate data. In contrast to the [`.Store`][store] method, this scratch is reset on server rebuilds.

[.Store][store]
: returns a Scratch to store and manipulate data. In contrast to the [`.Scratch`][scratch] method, this scratch is not reset on server rebuilds.

## Section Variables and Methods

Also see [Sections](/content-management/sections/).

{{< readfile file="/content/en/readfiles/sectionvars.md" markdown="true" >}}

## The `.Pages` Variable {#pages}

`.Pages` is an alias to `.Data.Pages`. It is conventional to use the
aliased form `.Pages`.

### `.Pages` compared to `.Site.Pages`

{{< getcontent path="readfiles/pages-vs-site-pages.md" >}}

## Page Fragments

{{< new-in "0.111.0" >}}

The `.Fragments` method returns a list of fragments for the current page.

.Headings
: A recursive list of headings for the current page. Can be used to generate a table of contents.

{{< todo >}}add .Headings toc example{{< /todo >}}

.Identifiers
: A sorted list of identifiers for the current page. Can be used to check if a page contains a specific identifier or if a page contains duplicate identifiers:

```go-html-template
{{ if .Fragments.Identifiers.Contains "my-identifier" }}
    <p>Page contains identifier "my-identifier"</p>
{{ end }}

{{ if gt (.Fragments.Identifiers.Count "my-identifier")  1 }}
    <p>Page contains duplicate "my-idenfifier" fragments</p>
{{ end }}
``` 

.HeadingsMap
: Holds a map of headings for the current page. Can be used to start the table of contents from a specific heading.

Also see the [Go Doc](https://pkg.go.dev/github.com/gohugoio/hugo@v0.111.0/markup/tableofcontents#Fragments) for the return type.

### Fragments in hooks and shortcodes

`.Fragments` are safe to call from render hooks, even on the page you're on (`.Page.Fragments`). For shortcodes we recommend that all `.Fragments` usage is nested inside the `{{</**/>}}` shortcode delimiter (`{{%/**/%}}` takes part in the ToC creation so it's easy to end up in a situation where you bite yourself in the tail).


## The global page function

{{< new-in "0.111.1" >}}

Hugo almost always passes a `Page` as the data context into the top level template (e.g. `single.html`) (the one exception is the multihost sitemap template). This means that you can access the current page with the `.` variable in the template.

But when you're deeply nested inside `.Render`, partial etc., accessing that `Page` object isn't always practical or possible.

For this reason, Hugo provides a global `page` function that you can use to access the current page from anywhere in any template.

```go-html-template
{{ page.Title }}
```

There are one caveat with this, and this isn't new, but it's worth mentioning here: There are situations in Hugo where you may see a cached value, e.g. when using `partialCached` or in a shortcode. 

## Page-level Params

Any other value defined in the front matter in a content file, including taxonomies, will be made available as part of the `.Params` variable.

{{< code-toggle file="content/example.md" fm=true copy=false >}}
title: Example
categories: [one]
tags: [two,three,four]
{{< /code-toggle >}}

With the above front matter, the `tags` and `categories` taxonomies are accessible via the following:

* `.Params.tags`
* `.Params.categories`

The `.Params` variable is particularly useful for the introduction of user-defined front matter fields in content files. For example, a Hugo website on book reviews could have the following front matter:

{{< code-toggle file="content/example.md" fm=true copy=false >}}
title: Example
affiliatelink: "http://www.my-book-link.here"
recommendedby: "My Mother"
{{< /code-toggle >}}

These fields would then be accessible to via `.Params.affiliatelink` and `.Params.recommendedby`.

```go-html-template
<h3><a href="{{ .Params.affiliatelink }}">Buy this book</a></h3>
<p>It was recommended by {{ .Params.recommendedby }}.</p>
```

This template would render as follows:

```html
<h3><a href="http://www.my-book-link.here">Buy this book</a></h3>
<p>It was recommended by my Mother.</p>
```

{{% note %}}
See [Archetypes](/content-management/archetypes/) for consistency of `Params` across pieces of content.
{{% /note %}}

### The `.Param` Method

In Hugo, you can declare params in individual pages and globally for your entire website. A common use case is to have a general value for the site param and a more specific value for some of the pages (i.e., a header image):

```go-html-template
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value according to it's definition in a page parameter (i.e. in the content's front matter) or a site parameter (i.e., in your `config`).

### Access Nested Fields in Front Matter

When front matter contains nested fields like the following:

{{< code-toggle file="content/example.md" fm=true copy=false >}}
title: Example
author:
  given_name: John
  family_name: Feminella
  display_name: John Feminella
{{< /code-toggle >}}

`.Param` can access these fields by concatenating the field names together with a dot:

```go-html-template
{{ $.Param "author.display_name" }}
```

[gitinfo]: /variables/git/
[File Variables]: /variables/files/
[bundle]: /content-management/page-bundles
[scratch]: /functions/scratch
[store]: /functions/store
