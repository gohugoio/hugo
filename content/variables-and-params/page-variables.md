---
title: Page Variables
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
tags: [pages]
draft: false
weight: 20
aliases: []
toc: false
needsreview: true
notesforauthors:
---

The following is a list of page-level variables that can be defined for a piece of content. Many of these will be defined in the front matter, derived from file location, or extracted from the content itself.

{{% note "`.Scratch`" %}}
See [`.Scratch`](/functions/scratch/) for page-scoped writable variables.
{{% /note %}}

* `.Content` The content itself, defined below the front matter.
* `.Data` The data specific to this type of page.
* `.Date` The date the page is associated with.
* `.Description` The description for the page.
* `.Draft` A boolean, `true` if the content is marked as a draft in the front matter.
* `.ExpiryDate` The date where the content is scheduled to expire on.
* `.FuzzyWordCount` The approximate number of words in the content.
* `.Hugo` See [Hugo Variables][hugovariables]
* `.IsHome` True if this is the home page.
* `.IsNode` Always false for regular content pages.
* `.IsPage` Always true for regular content pages.
* `.IsTranslated` Whether there are any translations to display.
* `.Keywords` The meta keywords for this content.
* `.Kind` What *kind* of page is this: is one of *page, home, section, taxonomy or taxonomyTerm.* There are also *RSS, sitemap, robotsTXT and 404*, but these will only available during rendering of that kind of page, and not available in any of the `Pages` collections.
* `.Lang` Language taken from the language extension notation.
* `.Language` A language object that points to this the language's definition in the site config.
* `.Lastmod` The date the content was last modified.
* `.LinkTitle` Access when creating links to this content. Will use `linktitle` if set in front matter, else `title`.
* `.Next` Pointer to the following content (based on pub date).
* `.NextInSection` Pointer to the following content within the same section (based on pub date)
* `.Pages` a collection of associated pages. This will be nil for regular content pages. This is an alias for `.Data.Pages`.
* `.Permalink` The Permanent link for this page.
* `.Prev` Pointer to the previous content (based on pub date).
* `.PrevInSection` Pointer to the previous content within the same section (based on `.PublishDate`). For example, `{{if .PrevInSection}}{{.PrevInSection.Permalink}}{{end}}`.
* `.PublishDate` The date the content is published on.
* `.RSSLink` Link to the taxonomies' RSS link.
* `.RawContent` Raw markdown content without the front matter. Useful with [remarkjs.com](http://remarkjs.com)
* `.ReadingTime` The estimated time it takes to read the content in minutes.
* `.Ref` Returns the permalink for a given reference;e.g., `.Ref "sample.md"`. See [Cross References][crossreferences]. Does not handle in-page fragments correctly.
* `.RelPermalink` The Relative permanent link for this page.
* `.RelRef` Returns the relative permalink for a given reference.  Example: `RelRef "sample.md"`. See [Cross References][crossreferences]. This does not handle in-page fragments.
* `.Section` The [section](/content-management/content-sections/) this content belongs to.
* `.Site` See [Site Variables][sitevariables] below.
* `.Summary` A generated summary of the content for easily showing a snippet in a summary view. Note that the breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page. See [Summaries](/content/summaries/) for more details.
* `.TableOfContents` The rendered [table of contents](/content-management/table-of-contents/) for the page.
* `.Title` The title for this page.
* `.Translations` A list of translated versions of the current page. See [Multilingual](/content-management/multilingual-mode/) for more info.
* `.Truncated` A boolean, `true` if the `.Summary` is truncated.  Useful for showing a "Read more..." link only if necessary.  See [Summaries](/content/summaries/) for more details.
* `.Type` The [content type][] (e.g., `post`).
* `.URL` The relative URL for this page. Note that if `URL` is set directly in front matter, that URL is returned as-is.
* `.UniqueID` The MD5-checksum of the content file's path
* `.Weight` Assigned weight (in the front matter) to this content, used in sorting.
* `.WordCount` The number of words in the content.

## Page-level Params

Any other value defined in the front matter, including taxonomies, will be made available as part of the `.Params variable.

For example, the *tags* and *categories* taxonomies are accessed with:

* `.Params.tags`
* `.Params.categories`

{{% note "Casing of Params" %}}
Page-level `.Params` are *only* accessible in lowercase.
{{% /note %}}

This is particularly useful for the introduction of user-defined fields in content files. For example, a Hugo website on book reviews could have the following front matter in `/content/review/book01.md`:

```yaml
---
...
affiliatelink: "http://www.my-book-link.here"
recommendedby: "My Mother"
---
```

Which would then be accessible to a template at `/themes/yourtheme/layouts/review/single.html` through `.Params.affiliatelink` and `.Params.recommendedby`, respectively. Two common situations where these could be introduced are as a value of a certain attribute (like `href=""` below) or by itself to be displayed. Sample syntaxes include:

```html
<h3><a href={{ printf "%s" $.Params.affiliatelink }}>Buy this book</a></h3>
<p>It was recommended by {{ .Params.recommendedby }}.</p>
```

which would render

```html
<h3><a href="http://www.my-book-link.here">Buy this book</a></h3>
<p>It was recommended by my Mother.</p>
```

{{% note %}}
See [Archetypes](/content-management/archetyps) for consistency of `Params` across pieces of content.
{{% /note %}}

### Param method

In Hugo, you can declare params both for the site and the individual page. A
common use case is to have a general value for the site and a more specific
value for some of the pages (i.e., a header image):

```golang
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value whether it's
in a page parameter or a site parameter.

When front matter contains nested fields like the following:

```yaml
---
author:
  given_name: John
  family_name: Feminella
  display_name: John Feminella
---
```
`.Param` can access them by concatenating the field names together with a
dot:

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
```

[content type]: /content-management/content-types/
[crossreferences]: /content-management/cross-references/
[hugovariables]: /variables-and-params/hugo-variables/
[sitevariables]: /variables-and-params/site-variables/