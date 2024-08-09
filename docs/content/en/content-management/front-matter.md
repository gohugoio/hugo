---
title: Front matter
description: Use front matter to add metadata to your content.
categories: [content management]
keywords: [front matter,yaml,toml,json,metadata,archetypes]
menu:
  docs:
    parent: content-management
    weight: 60
weight: 60
toc: true
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

[json]: https://www.json.org/
[toml]: https://toml.io/
[yaml]: https://yaml.org/

See examples of front matter delimiters by toggling between the serialization formats below.

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
date = 2024-02-02T04:14:54-08:00
draft = false
weight = 10
[params]
author = 'John Smith'
{{< /code-toggle >}}

Front matter fields may be [scalar], [arrays], or [maps] containing [boolean], [integer], [float], or [string] values. Note that the TOML format also supports date/time values using unquoted strings.

[scalar]: /getting-started/glossary/#scalar
[arrays]: /getting-started/glossary/#array
[maps]: /getting-started/glossary/#map
[boolean]: /getting-started/glossary/#boolean
[integer]: /getting-started/glossary/#integer
[float]: /getting-started/glossary/#float
[string]: /getting-started/glossary/#string

## Fields

The most common front matter fields are `date`, `draft`, `title`, and `weight`, but you can specify metadata using any of fields below.

{{% note %}}
The field names below are reserved. For example, you cannot create a custom field named `type`. Create custom fields under the `params` key. See the [parameters] section for details.

[parameters]: #parameters
{{% /note %}}

###### aliases

(`string array`) An array of one or more aliases, where each alias is a relative URL that will redirect the browser to the current location. Access these values from a template using the [`Aliases`] method on a `Page` object. See the [aliases] section for details.

[`aliases`]: /methods/page/aliases/
[aliases]: /content-management/urls/#aliases

###### build

(`map`) A map of [build options].

[build options]: /content-management/build-options/

###### cascade {#cascade-field}

(`map`) A map of front matter keys whose values are passed down to the page’s descendants unless overwritten by self or a closer ancestor’s cascade. See the [cascade] section for details.

[cascade]: #cascade

###### date

(`string`) The date associated with the page, typically the creation date. Note that the TOML format also supports date/time values using unquoted strings. Access this value from a template using the [`Date`] method on a `Page` object.

[`date`]: /methods/page/date/

###### description

(`string`) Conceptually different than the page `summary`, the description is typically rendered within a `meta` element within the `head` element of the published HTML file. Access this value from a template using the [`Description`] method on a `Page` object.

[`description`]: /methods/page/description/

###### draft

(`bool`)
If `true`, the page will not be rendered unless you pass the `--buildDrafts` flag to the `hugo` command. Access this value from a template using the [`Draft`] method on a `Page` object.

[`draft`]: /methods/page/draft/

###### expiryDate

(`string`) The page expiration date. On or after the expiration date, the page will not be rendered unless you pass the `--buildExpired` flag to the `hugo` command. Note that the TOML format also supports date/time values using unquoted strings. Access this value from a template using the [`ExpiryDate`] method on a `Page` object.

[`expirydate`]: /methods/page/expirydate/

###### headless

(`bool`) Applicable to [leaf bundles], if `true` this value sets the `render` and `list` [build options] to `never`, creating a headless bundle of [page resources].

[leaf bundles]: /content-management/page-bundles/#leaf-bundles
[page resources]: /content-management/page-resources/

###### isCJKLanguage

(`bool`) Set to `true` if the content language is in the [CJK] family. This value determines how Hugo calculates word count, and affects the values returned by the [`WordCount`], [`FuzzyWordCount`], [`ReadingTime`], and [`Summary`] methods on a `Page` object.

[`fuzzywordcount`]: /methods/page/wordcount/
[`readingtime`]: /methods/page/readingtime/
[`summary`]: /methods/page/summary/
[`wordcount`]: /methods/page/wordcount/
[cjk]: /getting-started/glossary/#cjk

###### keywords

(`string array`) An array of keywords, typically rendered within a `meta` element within the `head` element of the published HTML file, or used as a [taxonomy] to classify content. Access these values from a template using the [`Keywords`] method on a `Page` object.

[`keywords`]: /methods/page/keywords/
[taxonomy]: /getting-started/glossary/#taxonomy

<!-- Added in v0.123.0 but purposefully omitted from documentation. -->
<!--
kind
: The kind of page, e.g. "page", "section", "home" etc. This is usually derived from the content path.
-->

<!-- Added in v0.123.0 but purposefully omitted from documentation. -->
<!--
lang
: The language code for this page. This is usually derived from the module mount or filename.
-->

###### lastmod

(`string`) The date that the page was last modified. Note that the TOML format also supports date/time values using unquoted strings. Access this value from a template using the [`Lastmod`] method on a `Page` object.

[`lastmod`]: /methods/page/date/

###### layout

(`string`) Provide a template name to [target a specific template],  overriding the default [template lookup order]. Set the value to the base file name of the template, excluding its extension. Access this value from a template using the [`Layout`] method on a `Page` object.

[`layout`]: /methods/page/layout/
[template lookup order]: /templates/lookup-order/
[target a specific template]: templates/lookup-order/#target-a-template

###### linkTitle

(`string`) Typically a shorter version of the `title`. Access this value from a template using the [`LinkTitle`] method on a `Page` object.

[`linktitle`]: /methods/page/linktitle/

###### markup

(`string`) An identifier corresponding to one of the supported [content formats]. If not provided, Hugo determines the content renderer based on the file extension.

[content formats]: /content-management/formats/#classification

###### menus

(`string`,`string array`, or `map`) If set, Hugo adds the page to the given menu or menus. See the [menus] page for details.

[menus]: /content-management/menus/#define-in-front-matter

###### outputs

(`string array`) The [output formats] to render.

[output formats]: /templates/output-formats/

<!-- Added in v0.123.0 but purposefully omitted from documentation. -->
<!--
path
: The canonical page path.
-->

###### params

{{< new-in 0.123.0 >}}

(`map`) A map of custom [page parameters].

[page parameters]: #parameters

###### publishDate

(`string`) The page publication date. Before the publication date, the page will not be rendered unless you pass the `--buildFuture` flag to the `hugo` command. Note that the TOML format also supports date/time values using unquoted strings. Access this value from a template using the [`PublishDate`] method on a `Page` object.

[`publishdate`]: /methods/page/publishdate/

###### resources

(`map array`) An array of maps to provide metadata for [page resources].

[page-resources]: /content-management/page-resources/#page-resources-metadata

###### sitemap

(`map`) A map of sitemap options. See the [sitemap templates] page for details. Access these values from a template using the [`Sitemap`] method on a `Page` object.

[sitemap templates]: /templates/sitemap/
[`sitemap`]: /methods/page/sitemap/

###### slug

(`string`) Overrides the last segment of the URL path. Not applicable to section pages. See the [URL management] page for details. Access this value from a template using the [`Slug`] method on a `Page` object.

[`slug`]: /methods/page/slug/
[URL management]: /content-management/urls/#slug

###### summary

(`string`) Conceptually different than the page `description`, the summary either summarizes the content or serves as a teaser to encourage readers to visit the page. Access this value from a template using the [`Summary`] method on a `Page` object.

[`Summary`]: /methods/page/summary/

###### title

(`string`) The page title. Access this value from a template using the [`Title`] method on a `Page` object.

[`title`]: /methods/page/title/

###### translationKey

(`string`) An arbitrary value used to relate two or more translations of the same page, useful when the translated pages do not share a common path. Access this value from a template using the [`TranslationKey`] method on a `Page` object.

[`translationkey`]: /methods/page/translationkey/

###### type

(`string`) The [content type], overriding the value derived from the top level section in which the page resides. Access this value from a template using the [`Type`] method on a `Page` object.

[content type]: /getting-started/glossary/#content-type
[`type`]: /methods/page/type/

###### url

(`string`) Overrides the entire URL path. Applicable to regular pages and section pages. See the [URL management] page for details.

###### weight
(`int`) The page [weight], used to order the page within a [page collection]. Access this value from a template using the [`Weight`] method on a `Page` object.

[page collection]: /getting-started/glossary/#page-collection
[weight]: /getting-started/glossary/#weight
[`weight`]: /methods/page/weight/

## Parameters

{{< new-in 0.123.0 >}}

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

[`param`]: /methods/page/param/
[`params`]: /methods/page/params/

Hugo provides [embedded templates] to optionally insert meta data within the `head` element of your rendered pages. These embedded templates expect the following front matter parameters:

Parameter|Data type|Used by these embedded templates
:--|:--|:--
`audio`|`[]string`|[`opengraph.html`]
`images`|`[]string`|[`opengraph.html`], [`schema.html`], [`twitter_cards.html`]
`videos`|`[]string`|[`opengraph.html`]

The embedded templates will skip a parameter if not provided in front matter, but will throw an error if the data type is unexpected. 

[`opengraph.html`]: {{% eturl opengraph %}}
[`schema.html`]: {{% eturl schema %}}
[`twitter_cards.html`]: {{% eturl twitter_cards %}}
[embedded templates]: /templates/embedded/

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

You can add taxonomy terms to the front matter of any these [page kinds]:

- `home`
- `page`
- `section`
- `taxonomy`
- `term`

[page kinds]: /getting-started/glossary/#page-kind

Access taxonomy terms from a template using the [`Params`] or [`GetTerms`] method on a `Page` object. For example:

{{< code file=layouts/_default/single.html >}}
{{ with .GetTerms "tags" }}
  <p>Tags</p>
  <ul>
    {{ range . }}
      <li><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></li>
    {{ end }}
  </ul>
{{ end }}
{{< /code >}}

[`Params`]: /methods/page/params/
[`GetTerms`]: /methods/page/getterms/

## Cascade

Any [node] can pass down to its descendants a set of front matter values.

[node]: /getting-started/glossary/#node

### Target specific pages

The `cascade` block can be an array with an optional `_target` keyword, allowing you to target different page sets while cascading values.

{{< code-toggle file=content/_index.md fm=true >}}
title ="Home"
[[cascade]]
[cascade.params]
background = "yosemite.jpg"
[cascade._target]
path="/articles/**"
lang="en"
kind="page"
[[cascade]]
[cascade.params]
background = "goldenbridge.jpg"
[cascade._target]
kind="section"
{{</ code-toggle >}}

Use any combination of these keywords to target a set of pages:

###### path {#cascade-path}

(`string`) A [Glob](https://github.com/gobwas/glob) pattern matching the content path below /content. Expects Unix-styled slashes. Note that this is the virtual path, so it starts at the mount root. The matching supports double-asterisks so you can match for patterns like `/blog/*/**` to match anything from the third level and down.

###### kind {#cascade-kind}

(`string`) A Glob pattern matching the Page's Kind(s), e.g. "{home,section}".

###### lang {#cascade-lang}

(`string`) A Glob pattern matching the Page's language, e.g. "{en,sv}".

###### environment {#cascade-environment}

(`string`) A Glob pattern matching the build environment, e.g. "{production,development}"

Any of the above can be omitted.

{{% note %}}
With a multilingual site it may be more efficient to define the `cascade` values in your site configuration to avoid duplicating the `cascade` values on the section, taxonomy, or term page for each language.

With a multilingual site, if you choose to define the `cascade` values in front matter, you must create a section, taxonomy, or term page for each language; the `lang` keyword is ignored.
{{% /note %}}

### Example

{{< code-toggle file=content/posts/_index.md fm=true >}}
date = 2024-02-01T21:25:36-08:00
title = 'Posts'
[cascade]
  [cascade.params]
    banner = 'images/typewriter.jpg'
{{</ code-toggle >}}

With the above example the posts section page and its descendants will return `images/typewriter.jpg` when `.Params.banner` is invoked unless:

- Said descendant has its own `banner` value set
- Or a closer ancestor node has its own `cascade.banner` value set.

## Emacs Org Mode

If your [content format] is [Emacs Org Mode], you may provide front matter using Org Mode keywords. For example:

{{< code file=content/example.org lang=text >}}
#+TITLE: Example
#+DATE: 2024-02-02T04:14:54-08:00
#+DRAFT: false
#+AUTHOR: John Smith
#+GENRES: mystery
#+GENRES: romance
#+TAGS: red
#+TAGS: blue
#+WEIGHT: 10
{{< /code >}}

Note that you can also specify array elements on a single line:

{{< code file=content/example.org lang=text >}}
#+TAGS[]: red blue
{{< /code >}}

[content format]: /content-management/formats/
[emacs org mode]: https://orgmode.org/
