---
title: Page Variables
linktitle:
description: Page-level variables are defined in a content file's front matter, derived from the content's file location, or extracted from the content body itself.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
tags: [pages]
draft: false
weight: 20
aliases: []
toc: true
---

The following is a list of page-level variables that can be defined for a piece of content. Many of these will be defined in the front matter, derived from file location, or extracted from the content itself.

{{% note "`.Scratch`" %}}
See [`.Scratch`](/functions/scratch/) for page-scoped writable variables.
{{% /note %}}

## Page Variables List

`.Content`
: the content itself, defined below the front matter.

`.Data`
: the data specific to this type of page.

`.Date`
: the date associated with the page.

`.Description`
: the description for the page.

`.Draft`
: a boolean, `true` if the content is marked as a draft in the front matter.

`.ExpiryDate`
: the date on which the content is scheduled to expire.

`.FuzzyWordCount`
: the approximate number of words in the content.

`.Hugo`
: see [Hugo Variables](/variables-and-params/shortcode-git-and-hugo-variables/).

`.IsHome`
: `true` in the context of the [home page](/templates/homepage-template/).

`.IsNode`
: always `false` for regular content pages.

`.IsPage`
: always `true` for regular content pages.

`.IsTranslated`
: `true` if there are translations to display.

`.Keywords`
: the meta keywords for the content.

`.Kind`
: the page's *kind*. Possible return values are `page`, `home`, `section`, `taxonomy`, or `taxonomyTerm`. Note that there are also `RSS`, `sitemap`, `robotsTXT`, and `404` kinds, but these are only available during the rendering of each of these respective page's kind and therefore *not* available in any of the `Pages` collections.

`.Lang`
: language taken from the language extension notation.

`.Language`
: a language object that points to the language's definition in the site
`config`.

`.Lastmod`
: the date the content was last modified (i.e. from `lastmod` in the content's front matter)

`.LinkTitle`
: access when creating links to the content. If set, Hugo will use the `linktitle` from the front matter before `title`.

`.Next`
: pointer to the following content (based on `publishdate` in front matter).

`.NextInSection`
: pointer to the following content within the same section (based on `publishdate` in front matter).

`.Pages`
: a collection of associated pages. This value will be `nil` for regular content pages. `.Pages` is an alias for `.Data.Pages`.

`.Permalink`
: the Permanent link for this page; see [Permalinks](/content-management/url-management/)

`.Prev`
: Pointer to the previous content (based on `publishdate` in front matter).

`.PrevInSection`
: Pointer to the previous content within the same section (based on `publishdate` in front matter). For example, `{{if .PrevInSection}}{{.PrevInSection.Permalink}}{{end}}`.

`.PublishDate`
: the date on which the content was or will be published.

`.RSSLink`
: link to the taxonomies' RSS link.

`.RawContent`
: raw markdown content without the front matter. Useful with [remarkjs.com](
http://remarkjs.com)

`.ReadingTime`
: the estimated time, in minutes, it takes to read the content.

`.Ref`
: returns the permalink for a given reference (e.g., `.Ref "sample.md"`).  `.Ref` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

`.RelPermalink`
: the relative permanent link for this page.

`.RelRef`
: returns the relative permalink for a given reference (e.g., `RelRef
"sample.md"`). `.RelRef` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

`.Section`
: the [section](/content-management/sections/) this content belongs to.

`.Site`
: see [Site Variables](/variables-and-params/site-variables/).

`.Summary`
: a generated summary of the content for easily showing a snippet in a summary view. The breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page. See [Content Summaries](/content-management/content-summaries/) for more details.

`.TableOfContents`
: the rendered [table of contents](/content-management/table-of-contents/) for the page.

`.Title`
: the title for this page.

`.Translations`
: a list of translated versions of the current page. See [Multilingual Mode](/content-management/multilingual-mode/) for more information.

`.Truncated`
: a boolean, `true` if the `.Summary` is truncated. Useful for showing a "Read more..." link only when necessary.  See [Summaries](/content-management/content-summaries/) for more information.

`.Type`
: the [content type](/content-management/content-types/) of the content (e.g., `post`).

`.URL`
: the relative URL for the page. Note that the `URL` set directly in front
matter overrides the default relative URL for the page.

`.UniqueID`
: the MD5-checksum of the content file's path.

`.Weight`
: assigned weight (in the front matter) to this content, used in sorting.

`.WordCount`
: the number of words in the content.

## Page-level Params

Any other value defined in the front matter in a content file, including taxonomies, will be made available as part of the `.Params` variable.

```yaml
---
title: My First Post
date: date: 2017-02-20T15:26:23-06:00
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

```yaml
---
...
affiliatelink: "http://www.my-book-link.here"
recommendedby: "My Mother"
...
---
```

These fields would then be accessible to the `/themes/yourtheme/layouts/review/single.html` template through `.Params.affiliatelink` and `.Params.recommendedby`, respectively.

Two common situations where this type of front matter field could be introduced is as a value of a certain attribute like `href=""` or by itself to be displayed as text to the website's visitors.

{{% input "/themes/yourtheme/layouts/review/single.html" %}}
```html
<h3><a href={{ printf "%s" $.Params.affiliatelink }}>Buy this book</a></h3>
<p>It was recommended by {{ .Params.recommendedby }}.</p>
```
{{% /input %}}

This template would render as follows, assuming you've set [`uglyURLs`](/content-management/url-management/) to `false` in your [site `config`](/getting-started/configuration/):

{{% output "yourbaseurl/review/book01/index.html" %}}
```html
<h3><a href="http://www.my-book-link.here">Buy this book</a></h3>
<p>It was recommended by my Mother.</p>
```
{{% /output %}}

{{% note %}}
See [Archetypes](/content-management/archetyps) for consistency of `Params` across pieces of content.
{{% /note %}}

### The `.Param` Method

In Hugo, you can declare params in individual pages and globally for your entire website. A common use case is to have a general value for the site param and a more specific value for some of the pages (i.e., a header image):

```golang
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value according to it's definition in a page parameter (i.e. in the content's front matter) or a site parameter (i.e., in your `config`).

### Accessing Nested Fields in Front Matter

When front matter contains nested fields like the following:

```yaml
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

```golang
{{ $.Param "favorites.flavor" }}
=> vanilla
```
