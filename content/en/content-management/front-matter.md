---
title: Front matter
description: Use front matter to add metadata to your content.
categories: []
keywords: []
aliases: [/content/front-matter/]
---

## Overview

The front matter at the top of each content file is metadata that:

- Describes the content
- Augments the content
- Establishes relationships with other content
- Controls the published structure of your site
- Determines template selection

Provide front matter using a serialization format, one of [JSON], [TOML], or [YAML]. Hugo determines the front matter format by examining the delimiters that separate the front matter from the page content.

See examples of front matter delimiters by toggling between the serialization formats below.

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
date = 2024-02-02T04:14:54-08:00
draft = false
weight = 10
[params]
author = 'John Smith'
{{< /code-toggle >}}

Front matter fields may be [boolean](g), [integer](g), [float](g), [string](g), [arrays](g), or [maps](g). Note that the TOML format also supports unquoted date/time values.

## Fields

The most common front matter fields are `date`, `draft`, `title`, and `weight`, but you can specify metadata using any of fields below.

> [!note]
> The field names below are reserved. For example, you cannot create a custom field named `type`. Create custom fields under the `params` key. See the [parameters] section for details.

[parameters]: #parameters

aliases
: (`string array`) An array of one or more aliases, where each alias is a relative URL that will redirect the browser to the current location. Access these values from a template using the [`Aliases`] method on a `Page` object. See the [aliases] section for details.

build
: (`map`) A map of [build options].

cascade
: (`map`) A map of front matter keys whose values are passed down to the page's descendants unless overwritten by self or a closer ancestor's cascade. See the [cascade] section for details.

date
: (`string`) The date associated with the page, typically the creation date. Note that the TOML format also supports unquoted date/time values. See the [dates](#dates) section for examples. Access this value from a template using the [`Date`] method on a `Page` object.

description
: (`string`) Conceptually different than the page `summary`, the description is typically rendered within a `meta` element within the `head` element of the published HTML file. Access this value from a template using the [`Description`] method on a `Page` object.

draft
: (`bool`) Whether to disable rendering unless you pass the `--buildDrafts` flag to the `hugo` command. Access this value from a template using the [`Draft`] method on a `Page` object.

expiryDate
: (`string`) The page expiration date. On or after the expiration date, the page will not be rendered unless you pass the `--buildExpired` flag to the `hugo` command. Note that the TOML format also supports unquoted date/time values. See the [dates](#dates) section for examples. Access this value from a template using the [`ExpiryDate`] method on a `Page` object.

headless
: (`bool`) Applicable to [leaf bundles], whether to set the `render` and `list` [build options] to `never`, creating a headless bundle of [page resources].

isCJKLanguage
: (`bool`) Whether the content language is in the [CJK](g) family. This value determines how Hugo calculates word count, and affects the values returned by the [`WordCount`], [`FuzzyWordCount`], [`ReadingTime`], and [`Summary`] methods on a `Page` object.

keywords
: (`string array`) An array of keywords, typically rendered within a `meta` element within the `head` element of the published HTML file, or used as a [taxonomy](g) to classify content. Access these values from a template using the [`Keywords`] method on a `Page` object.

lastmod
: (`string`) The date that the page was last modified. Note that the TOML format also supports unquoted date/time values. See the [dates](#dates) section for examples. Access this value from a template using the [`Lastmod`] method on a `Page` object.

layout
: (`string`) Provide a template name to [target a specific template],  overriding the default [template lookup order]. Set the value to the base file name of the template, excluding its extension. Access this value from a template using the [`Layout`] method on a `Page` object.

linkTitle
: (`string`) Typically a shorter version of the `title`. Access this value from a template using the [`LinkTitle`] method on a `Page` object.

markup
: (`string`) An identifier corresponding to one of the supported [content formats]. If not provided, Hugo determines the content renderer based on the file extension.

menus
: (`string`, `string array`, or `map`) If set, Hugo adds the page to the given menu or menus. See the [menus] page for details.

modified
: Alias to [lastmod](#lastmod).

outputs
: (`string array`) The [output formats] to render. See [configure outputs] for more information.

params
: {{< new-in 0.123.0 />}}
: (`map`) A map of custom [page parameters].

pubdate
: Alias to [publishDate](#publishdate).

publishDate
: (`string`) The page publication date. Before the publication date, the page will not be rendered unless you pass the `--buildFuture` flag to the `hugo` command. Note that the TOML format also supports unquoted date/time values. See the [dates](#dates) section for examples. Access this value from a template using the [`PublishDate`] method on a `Page` object.

published
: Alias to [publishDate](#publishdate).

resources
: (`map array`) An array of maps to provide metadata for [page resources].

sitemap
: (`map`) A map of sitemap options. See the [sitemap templates] page for details. Access these values from a template using the [`Sitemap`] method on a `Page` object.

slug
: (`string`) Overrides the last segment of the URL path. Not applicable to `home`, `section`, `taxonomy`, or `term` pages. See the [URL management] page for details. Access this value from a template using the [`Slug`] method on a `Page` object.

summary
: (`string`) Conceptually different than the page `description`, the summary either summarizes the content or serves as a teaser to encourage readers to visit the page. Access this value from a template using the [`Summary`] method on a `Page` object.

title
: (`string`) The page title. Access this value from a template using the [`Title`] method on a `Page` object.

translationKey
: (`string`) An arbitrary value used to relate two or more translations of the same page, useful when the translated pages do not share a common path. Access this value from a template using the [`TranslationKey`] method on a `Page` object.

type
: (`string`) The [content type](g), overriding the value derived from the top-level section in which the page resides. Access this value from a template using the [`Type`] method on a `Page` object.

unpublishdate
: Alias to [expirydate](#expirydate).

url
: (`string`) Overrides the entire URL path. Applicable to regular pages and section pages. See the [URL management] page for details.

weight
: (`int`) The page [weight](g), used to order the page within a [page collection](g). Access this value from a template using the [`Weight`] method on a `Page` object.

## Parameters

{{< new-in 0.123.0 />}}

Specify custom page parameters under the `params` key in front matter:

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
date = 2024-02-02T04:14:54-08:00
draft = false
weight = 10
[params]
author = 'John Smith'
{{< /code-toggle >}}

Access these values from a template using the [`Params`] or [`Param`] method on a `Page` object.

Hugo provides [embedded templates] to optionally insert meta data within the `head` element of your rendered pages. These embedded templates expect the following front matter parameters:

Parameter|Data type|Used by these embedded templates
:--|:--|:--
`audio`|`[]string`|[`opengraph.html`]
`images`|`[]string`|[`opengraph.html`], [`schema.html`], [`twitter_cards.html`]
`videos`|`[]string`|[`opengraph.html`]

The embedded templates will skip a parameter if not provided in front matter, but will throw an error if the data type is unexpected.

## Taxonomies

Classify content by adding taxonomy terms to front matter. For example, with this site configuration:

{{< code-toggle file=hugo >}}
[taxonomies]
tag = 'tags'
genre = 'genres'
{{< /code-toggle >}}

Add taxonomy terms as shown below:

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
date = 2024-02-02T04:14:54-08:00
draft = false
weight = 10
tags = ['red','blue']
genres = ['mystery','romance']
[params]
author = 'John Smith'
{{< /code-toggle >}}

You can add taxonomy terms to the front matter of any these [page kinds](g):

- `home`
- `page`
- `section`
- `taxonomy`
- `term`

Access taxonomy terms from a template using the [`Params`] or [`GetTerms`] method on a `Page` object. For example:

```go-html-template {file="layouts/page.html"}
{{ with .GetTerms "tags" }}
  <p>Tags</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
```

[`GetTerms`]: /methods/page/getterms/

## Cascade

A [node](g) can cascade front matter values to its descendants. However, this cascading will be prevented if the descendant already defines the field, or if a closer ancestor node has already cascaded a value for that same field.

For example, to cascade a "color" parameter from the home page to all its descendants:

{{< code-toggle file=content/_index.md fm=true >}}
title = 'Home'
[cascade.params]
color = 'red'
{{< /code-toggle >}}

### Target

<!-- TODO
Update the <version> and <date> below when we actually get around to deprecating _target.

We deprecated the `_target` front matter key in favor of `target` in <version> on <date>. Remove footnote #1 on or after 2026-03-10 (15 months after deprecation).
-->

The `target`[^1] keyword allows you to target specific pages or [environments](g). For example, to cascade a "color" parameter from the home page only to pages within the "articles" section, including the "articles" section page itself:

[^1]: The `_target` alias for `target` is deprecated and will be removed in a future release.

{{< code-toggle file=content/_index.md fm=true >}}
title = 'Home'
[cascade.params]
color = 'red'
[cascade.target]
path = '{/articles,/articles/**}'
{{< /code-toggle >}}

Use any combination of these keywords to target pages and/or environments:

environment
: (`string`) A [glob](g) pattern matching the build [environment](g). For example: `{staging,production}`.

kind
: (`string`) A [glob](g) pattern matching the [page kind](g). For example: `{taxonomy,term}`.

path
: (`string`) A [glob](g) pattern matching the page's [logical path](g). For example: `{/books,/books/**}`.

### Array

Define an array of cascade parameters to apply different values to different targets. For example:

{{< code-toggle file=content/_index.md fm=true >}}
title = 'Home'
[[cascade]]
[cascade.params]
color = 'red'
[cascade.target]
path = '{/books/**}'
kind = 'page'
[[cascade]]
[cascade.params]
color = 'blue'
[cascade.target]
path = '{/films/**}'
kind = 'page'
{{< /code-toggle >}}

> [!note]
> For multilingual sites, defining cascade values in your site configuration is often more efficient. This avoids repeating the same cascade values on the home, section, taxonomy, or term page for each language. See&nbsp;[details](/configuration/cascade/).
>
> If you choose to define cascade values in front matter for a multilingual site, you must create a corresponding home, section, taxonomy, or term page for every language.

## Emacs Org Mode

If your [content format] is [Emacs Org Mode], you may provide front matter using Org Mode keywords. For example:

```text {file="content/example.org"}
#+TITLE: Example
#+DATE: 2024-02-02T04:14:54-08:00
#+DRAFT: false
#+AUTHOR: John Smith
#+GENRES: mystery
#+GENRES: romance
#+TAGS: red
#+TAGS: blue
#+WEIGHT: 10
```

Note that you can also specify array elements on a single line:

```text {file="content/example.org"}
#+TAGS[]: red blue
```

[content format]: /content-management/formats/
[emacs org mode]: https://orgmode.org/

## Dates

When populating a date field, whether a [custom page parameter](#parameters) or one of the four predefined fields ([`date`](#date), [`expiryDate`](#expirydate), [`lastmod`](#lastmod), [`publishDate`](#publishdate)), use one of these parsable formats:

{{% include "/_common/parsable-date-time-strings.md" %}}

To override the default time zone, set the [`timeZone`](/configuration/all/#timezone) in your site configuration. The order of precedence for determining the time zone is:

1. The time zone offset in the date/time string
1. The time zone specified in your site configuration
1. The `Etc/UTC` time zone

[`aliases`]: /methods/page/aliases/
[`date`]: /methods/page/date/
[`description`]: /methods/page/description/
[`draft`]: /methods/page/draft/
[`expirydate`]: /methods/page/expirydate/
[`fuzzywordcount`]: /methods/page/wordcount/
[`keywords`]: /methods/page/keywords/
[`lastmod`]: /methods/page/date/
[`layout`]: /methods/page/layout/
[`linktitle`]: /methods/page/linktitle/
[`opengraph.html`]: <{{% eturl opengraph %}}>
[`Param`]: /methods/page/param/
[`Params`]: /methods/page/params/
[`publishdate`]: /methods/page/publishdate/
[`readingtime`]: /methods/page/readingtime/
[`schema.html`]: <{{% eturl schema %}}>
[`sitemap`]: /methods/page/sitemap/
[`slug`]: /methods/page/slug/
[`Summary`]: /methods/page/summary/
[`title`]: /methods/page/title/
[`translationkey`]: /methods/page/translationkey/
[`twitter_cards.html`]: <{{% eturl twitter_cards %}}>
[`type`]: /methods/page/type/
[`weight`]: /methods/page/weight/
[`wordcount`]: /methods/page/wordcount/
[aliases]: /content-management/urls/#aliases
[build options]: /content-management/build-options/
[cascade]: #cascade-1
[configure outputs]: /configuration/outputs/#outputs-per-page
[content formats]: /content-management/formats/#classification
[embedded templates]: /templates/embedded/
[json]: https://www.json.org/
[leaf bundles]: /content-management/page-bundles/#leaf-bundles
[menus]: /content-management/menus/#define-in-front-matter
[output formats]: /configuration/output-formats/
[page parameters]: #parameters
[page resources]: /content-management/page-resources/#metadata
[sitemap templates]: /templates/sitemap/
[target a specific template]: /templates/lookup-order/#target-a-template
[template lookup order]: /templates/lookup-order/
[toml]: https://toml.io/
[URL management]: /content-management/urls/#slug
[yaml]: https://yaml.org/
