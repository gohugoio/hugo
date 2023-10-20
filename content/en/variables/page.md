---
title: Page variables
description: Page-level variables are defined in a content file's front matter, derived from the content's file location, or extracted from the content body itself.
categories: [variables and parameters]
keywords: [pages]
menu:
  docs:
    parent: variables
    weight: 20
weight: 20
toc: true
---

The following is a list of page-level variables. Many of these will be defined in the front matter, derived from file location, or extracted from the content itself.

## Page variables

.AlternativeOutputFormats
: Contains all alternative formats for a given page; this variable is especially useful `link rel` list in your site's `<head>`. (See [Output Formats](/templates/output-formats/).)

.Aliases
: Aliases of this page

.Ancestors
: Ancestors of this page.

.BundleType
: The [bundle] type: `leaf`, `branch`, or an empty string if the page is not a bundle.

.Content
: The content itself, defined below the front matter.

.Data
: The data specific to this type of page.

.Date
: The date associated with the page. By default, this is the front matter `date` value. See [configuring dates] for a description of fallback values and precedence. See also `.ExpiryDate`, `.Lastmod`, and `.PublishDate`.

.Description
: The description for the page.

.Draft
: A boolean, `true` if the content is marked as a draft in the front matter.

.ExpiryDate
: The date on which the content is scheduled to expire. By default, this is the front matter `expiryDate` value. See [configuring dates] for a description of fallback values and precedence. See also `.Date`, `.Lastmod`, and `.PublishDate`.

.File
: Filesystem-related data for this content file. See also [File Variables].

.Fragments
: Fragments returns the fragments for this page. See [Page Fragments](#page-fragments).

.FuzzyWordCount
: The approximate number of words in the content.

.IsHome
: `true` in the context of the [homepage](/templates/homepage/).

.IsNode
: Always `false` for regular content pages.

.IsPage
: Always `true` for regular content pages.

.IsSection
: `true` if [`.Kind`](/templates/section-templates/#page-kinds) is `section`.

.IsTranslated
: `true` if there are translations to display.

.Keywords
: The meta keywords for the content.

.Kind
: The page's *kind*. Possible return values are `page`, `home`, `section`, `taxonomy`, or `term`. Note that there are also `RSS`, `sitemap`, `robotsTXT`, and `404` kinds, but these are only available during the rendering of each of these respective page's kind and therefore *not* available in any of the `Pages` collections.

.Language
: A language object that points to the language's definition in the site configuration. `.Language.Lang` gives you the language code.

.Lastmod
: The date on which the content was last modified. By default, if `enableGitInfo` is `true` in your site configuration, this is the Git author date, otherwise the front matter `lastmod` value. See [configuring dates] for a description of fallback values and precedence. See also `.Date`,`ExpiryDate`, `.PublishDate`, and [`.GitInfo`][gitinfo].

.LinkTitle
: Access when creating links to the content. If set, Hugo will use the `linktitle` from the front matter before `title`.

.Next
: Points up to the next regular page (sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath)). Example: `{{ with .Next }}{{ .Permalink }}{{ end }}`. Calling `.Next` from the first page returns `nil`.

.NextInSection
: Points up to the next regular page below the same top level section (e.g. in `/blog`)). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath). Example: `{{ with .NextInSection }}{{ .Permalink }}{{ end }}`. Calling `.NextInSection` from the first page returns `nil`.

.OutputFormats
: Contains all formats, including the current format, for a given page. Can be combined the with [`.Get` function](/functions/get/) to grab a specific format. (See [Output Formats](/templates/output-formats/).)

.Permalink
: The Permanent link for this page; see [Permalinks](/content-management/urls/)

.Plain
: The Page content stripped of HTML tags and presented as a string. You may need to pipe the result through the [`htmlUnescape`](/functions/transform/htmlunescape) function when rendering this value with the HTML [output format](/templates/output-formats#output-format-definitions).

.PlainWords
: The slice of strings that results from splitting .Plain into words, as defined in Go's [strings.Fields](https://pkg.go.dev/strings#Fields).

.Prev
: Points down to the previous regular page(sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath)). Example: `{{ if .Prev }}{{ .Prev.Permalink }}{{ end }}`.  Calling `.Prev` from the last page returns `nil`.

.PrevInSection
: Points down to the previous regular page below the same top level section (e.g. `/blog`). Pages are sorted by Hugo's [default sort](/templates/lists#default-weight--date--linktitle--filepath). Example: `{{ if .PrevInSection }}{{ .PrevInSection.Permalink }}{{ end }}`.  Calling `.PrevInSection` from the last page returns `nil`.

.PublishDate
: The date on which the content was or will be published. By default, this is the front matter `publishDate` value. See [configuring dates] for a description of fallback values and precedence. See also `.Date`, `.ExpiryDate`, and `.Lastmod`.

.RawContent
: Raw markdown content without the front matter. Useful with [remarkjs.com](
https://remarkjs.com)

.RenderShortcodes
: See [Render Shortcodes](#rendershortcodes).

.ReadingTime
: The estimated time, in minutes, it takes to read the content.

.Resources
: Resources such as images and CSS that are associated with this page

.Ref
: Returns the permalink for a given reference (e.g., `.Ref "sample.md"`).  `.Ref` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

.RelPermalink
: The relative permanent link for this page.

.RelRef
: Returns the relative permalink for a given reference (e.g., `RelRef
"sample.md"`). `.RelRef` does *not* handle in-page fragments correctly. See [Cross References](/content-management/cross-references/).

.Site
: See [Site Variables](/variables/site/).

.Sites
: Returns all sites (languages). A typical use case would be to link back to the main language: `<a href="{{ .Sites.First.Home.RelPermalink }}">...</a>`.

.Sites.First
: Returns the site for the first language. If this is not a multilingual setup, it will return itself.

.Summary
: A generated summary of the content for easily showing a snippet in a summary view. The breakpoint can be set manually by inserting <code>&lt;!&#x2d;&#x2d;more&#x2d;&#x2d;&gt;</code> at the appropriate place in the content page, or the summary can be written independent of the page text.  See [Content Summaries](/content-management/summaries/) for more details.

.TableOfContents
: The rendered [table of contents](/content-management/toc/) for the page.

.Title
: The title for this page.

.Translations
: A list of translated versions of the current page. See [Multilingual Mode](/content-management/multilingual/) for more information.

.TranslationKey
: The key used to map language translations of the current page. See [Multilingual Mode](/content-management/multilingual/) for more information.

.Truncated
: A boolean, `true` if the `.Summary` is truncated. Useful for showing a "Read more..." link only when necessary.  See [Summaries](/content-management/summaries/) for more information.

.Type
: The [content type](/content-management/types/) of the content (e.g., `posts`).

.Weight
: Assigned weight (in the front matter) to this content, used in sorting.

.WordCount
: The number of words in the content.

## Page collections

List pages receive the following page collections in [context](/getting-started/glossary/#context):

.Pages
: Regular pages within the current section (not recursive), and section pages for immediate descendant sections (not recursive).

.RegularPages
: Regular pages within the current section (not recursive).

.RegularPagesRecursive
: Regular pages within the current section, and regular pages within all descendant sections.

## Writable page-scoped variables

[.Scratch][scratch]
: Returns a Scratch to store and manipulate data. In contrast to the [`.Store`][store] method, this scratch is reset on server rebuilds.

[.Store][store]
: Returns a Scratch to store and manipulate data. In contrast to the [`.Scratch`][scratch] method, this scratch is not reset on server rebuilds.

## Section variables and methods

Also see [Sections](/content-management/sections/).

.CurrentSection
: The page's current section. The value can be the page itself if it is a section or the homepage.

.FirstSection
: The page's first section below root, e.g. `/docs`, `/blog` etc.

.InSection $anotherPage
: Whether the given page is in the current section.

.IsAncestor $anotherPage
: Whether the current page is an ancestor of the given page.

.IsDescendant $anotherPage
: Whether the current page is a descendant of the given page.

.Parent
: A section's parent section or a page's section.

.Section
: The [section](/content-management/sections/) this content belongs to. **Note:** For nested sections, this is the first path element in the directory, for example, `/blog/funny/mypost/ => blog`.

.Sections
: The [sections](/content-management/sections/) below this content.

## Page fragments

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
    <p>Page contains duplicate "my-identifier" fragments</p>
{{ end }}
```

.HeadingsMap
: Holds a map of headings for the current page. Can be used to start the table of contents from a specific heading.

Also see the [Go Doc](https://pkg.go.dev/github.com/gohugoio/hugo@v0.111.0/markup/tableofcontents#Fragments) for the return type.

### Fragments in hooks and shortcodes

`.Fragments` are safe to call from render hooks, even on the page you're on (`.Page.Fragments`). For shortcodes we recommend that all `.Fragments` usage is nested inside the `{{</**/>}}` shortcode delimiter (`{{%/**/%}}` takes part in the ToC creation so it's easy to end up in a situation where you bite yourself in the tail).

## The `.RenderShortcodes` method {#rendershortcodes}

{{< new-in "0.117.0" >}} This renders all the shortcodes in the content, preserving the surrounding markup (e.g. Markdown) as is.

The common use case this is to composing a page from multiple content files while preserving a global context for table of contents and foot notes.

This method is most often used in shortcode templates. A simple example of shortcode template including content from another page would look like:

```go-html-template
{{ $p := site.GetPage (.Get 0) }}
{{ $p.RenderShortcodes }}
```

In the above it's important to understand  and the difference between the two delimiters used when including a shortcode:

* `{{</* myshortcode */>}}` tells Hugo that the rendered shortcode does not need further processing (e.g. it's HTML).
* `{{%/* myshortcode */%}}` tells Hugo that the rendered shortcode needs further processing (e.g. it's Markdown).

The latter is what you want to use for the include shortcode outlined above:

```md
## Mypage
{{%/* include "mypage" */%}}
``````


Also see [Use Shortcodes](/content-management/shortcodes/#use-shortcodes).

## Page-level params

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

### The `.Param` method

In Hugo, you can declare parameters in individual pages and globally for your entire website. A common use case is to have a general value for the site parameter and a more specific value for some of the pages (i.e., a header image):

```go-html-template
{{ $.Param "header_image" }}
```

The `.Param` method provides a way to resolve a single value according to it's definition in a page parameter (i.e. in the content's front matter) or a site parameter (i.e., in your site configuration).

### Access nested fields in front matter

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

[configuring dates]: /getting-started/configuration/#configure-dates
[gitinfo]: /variables/git/
[File Variables]: /variables/files/
[bundle]: /content-management/page-bundles
[scratch]: /functions/scratch
[store]: /functions/store
