---
title: Front Matter
description: Hugo allows you to add front matter in yaml, toml, or json to your content files.
categories: [content management]
keywords: ["front matter", "yaml", "toml", "json", "metadata", "archetypes"]
menu:
  docs:
    parent: content-management
    weight: 60
toc: true
weight: 60
aliases: [/content/front-matter/]
---

**Front matter** allows you to keep metadata attached to an instance of a [content type]---i.e., embedded inside a content file---and is one of the many features that gives Hugo its strength.

{{< youtube Yh2xKRJGff4 >}}

## Front Matter Formats

Hugo supports four formats for front matter, each with their own identifying tokens.

TOML
: identified by opening and closing `+++`.

YAML
: identified by opening and closing `---`.

JSON
: a single JSON object surrounded by '`{`' and '`}`', followed by a new line.

ORG
: a group of Org mode keywords in the format '`#+KEY: VALUE`'. Any line that does not start with `#+` ends the front matter section.
  Keyword values can be either strings (`#+KEY: VALUE`) or a whitespace separated list of strings (`#+KEY[]: VALUE_1 VALUE_2`).

### Example

{{< code-toggle >}}
title = "spf13-vim 3.0 release and new website"
description = "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
tags = [ ".vimrc", "plugins", "spf13-vim", "vim" ]
date = "2012-04-06"
categories = [
  "Development",
  "VIM"
]
slug = "spf13-vim-3-0-release-and-new-website"
{{< /code-toggle >}}

## Front Matter Variables

### Predefined

There are a few predefined variables that Hugo is aware of. See [Page Variables][pagevars] for how to call many of these predefined variables in your templates.

aliases
: An array of one or more aliases (e.g., old published paths of renamed content) that will be created in the output directory structure . See [Aliases][aliases] for details.

audio
: An array of paths to audio files related to the page; used by the `opengraph` [internal template](/templates/internal) to populate `og:audio`.

cascade
: A map of front matter keys whose values are passed down to the page's descendants unless overwritten by self or a closer ancestor's cascade. See [Front Matter Cascade](#front-matter-cascade) for details.

date
: The datetime assigned to this page. This is usually fetched from the `date` field in front matter, but this behavior is configurable.

description
: The description for the content.

draft
: If `true`, the content will not be rendered unless the `--buildDrafts` flag is passed to the `hugo` command.

expiryDate
: The datetime at which the content should no longer be published by Hugo; expired content will not be rendered unless the `--buildExpired` flag is passed to the `hugo` command.

headless
: If `true`, sets a leaf bundle to be [headless][headless-bundle].

images
: An array of paths to images related to the page; used by [internal templates](/templates/internal) such as `_internal/twitter_cards.html`.

isCJKLanguage
: If `true`, Hugo will explicitly treat the content as a CJK language; both `.Summary` and `.WordCount` work properly in CJK languages.

keywords
: The meta keywords for the content.

layout
: The layout Hugo should select from the [lookup order][lookup] when rendering the content. If a `type` is not specified in the front matter, Hugo will look for the layout of the same name in the layout directory that corresponds with a content's section. See [Content Types][content type].

lastmod
: The datetime at which the content was last modified.

linkTitle
: Used for creating links to content; if set, Hugo defaults to using the `linktitle` before the `title`. Hugo can also [order lists of content by `linktitle`][bylinktitle].

markup
: **experimental**; specify `"rst"` for reStructuredText (requires`rst2html`) or `"md"` (default) for Markdown.

outputs
: Allows you to specify output formats specific to the content. See [output formats][outputs].

publishDate
: If in the future, content will not be rendered unless the `--buildFuture` flag is passed to `hugo`.

resources
: Used for configuring page bundle resources. See [Page Resources][page-resources].

series
: An array of series this page belongs to, as a subset of the `series` [taxonomy](/content-management/taxonomies/); used by the `opengraph` [internal template](/templates/internal) to populate `og:see_also`.

slug
: Overrides the last segment of the URL path. Not applicable to section pages. See [URL Management](/content-management/urls/#slug) for details.

summary
: Text used when providing a summary of the article in the `.Summary` page variable; details available in the [content-summaries](/content-management/summaries/) section.

title
: The title for the content.

type
: The type of the content; this value will be automatically derived from the directory (i.e., the [section]) if not specified in front matter.

url
: Overrides the entire URL path. Applicable to regular pages and section pages. See [URL Management](/content-management/urls/#url) for details.

videos
: An array of paths to videos related to the page; used by the `opengraph` [internal template](/templates/internal) to populate `og:video`.

weight
: used for [ordering your content in lists][ordering]. Lower weight gets higher precedence. So content with lower weight will come first. If set, weights should be non-zero, as 0 is interpreted as an *unset* weight.

\<taxonomies\>
: Field name of the *plural* form of the index. See `tags` and `categories` in the above front matter examples. *Note that the plural form of user-defined taxonomies cannot be the same as any of the predefined front matter variables.*

{{% note %}}
If neither `slug` nor `url` is present and [permalinks are not configured otherwise in your site configuration file](/content-management/urls/#permalinks), Hugo will use the filename of your content to create the output URL. See [Content Organization](/content-management/organization) for an explanation of paths in Hugo and [URL Management](/content-management/urls/) for ways to customize Hugo's default behaviors.
{{% /note %}}

### User-Defined

You can add fields to your front matter arbitrarily to meet your needs. These user-defined key-values are placed into a single `.Params` variable for use in your templates.

The following fields can be accessed via `.Params.include_toc` and `.Params.show_comments`, respectively. The [Variables] section provides more information on using Hugo's page- and site-level variables in your templates.

{{< code-toggle copy=false >}}
include_toc: true
show_comments: false
{{</ code-toggle >}}

## Front Matter Cascade

Any node or section can pass down to descendants a set of Front Matter values as long as defined underneath the reserved `cascade` Front Matter key.

### Target Specific Pages

The `cascade` block can be a slice with a optional `_target` keyword, allowing for multiple `cascade` values targeting different page sets.

{{< code-toggle copy=false >}}
title ="Blog"
[[cascade]]
background = "yosemite.jpg"
[cascade._target]
path="/blog/**"
lang="en"
kind="page"
[[cascade]]
background = "goldenbridge.jpg"
[cascade._target]
kind="section"
{{</ code-toggle >}}

Keywords available for `_target`:

path
: A [Glob](https://github.com/gobwas/glob) pattern matching the content path below /content. Expects Unix-styled slashes. Note that this is the virtual path, so it starts at the mount root. The matching supports double-asterisks so you can match for patterns like `/blog/*/**` to match anything from the third level and down.

kind
: A Glob pattern matching the Page's Kind(s), e.g. "{home,section}".

lang
: A Glob pattern matching the Page's language, e.g. "{en,sv}".

environment
: A Glob pattern matching the build environment, e.g. "{production,development}"

Any of the above can be omitted.

### Example

In `content/blog/_index.md`

{{< code-toggle copy=false >}}
title: Blog
cascade:
  banner: images/typewriter.jpg
{{</ code-toggle >}}

With the above example the Blog section page and its descendants will return `images/typewriter.jpg` when `.Params.banner` is invoked unless:

- Said descendant has its own `banner` value set
- Or a closer ancestor node has its own `cascade.banner` value set.

## Order Content Through Front Matter

You can assign content-specific `weight` in the front matter of your content. These values are especially useful for [ordering][ordering] in list views. You can use `weight` for ordering of content and the convention of [`<TAXONOMY>_weight`][taxweight] for ordering content within a taxonomy. See [Ordering and Grouping Hugo Lists][lists] to see how `weight` can be used to organize your content in list views.

## Override Global Markdown Configuration

It's possible to set some options for Markdown rendering in a content's front matter as an override to the [Rendering options set in your project configuration][config].

## Front Matter Format Specs

- [TOML Spec][toml]
- [YAML Spec][yaml]
- [JSON Spec][json]

[variables]: /variables/
[aliases]: /content-management/urls/#aliases
[archetype]: /content-management/archetypes/
[bylinktitle]: /templates/lists/#by-link-title
[config]: /getting-started/configuration/ "Hugo documentation for site configuration"
[content type]: /content-management/types/
[contentorg]: /content-management/organization/
[headless-bundle]: /content-management/page-bundles/#headless-bundle
[json]: https://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf "Specification for JSON, JavaScript Object Notation"
[lists]: /templates/lists/#order-content "See how to order content in list pages; for example, templates that look to specific _index.md for content and front matter."
[lookup]: /templates/lookup-order/ "Hugo traverses your templates in a specific order when rendering content to allow for DRYer templating."
[ordering]: /templates/lists/ "Hugo provides multiple ways to sort and order your content in list templates"
[outputs]: /templates/output-formats/ "With the release of v22, you can output your content to any text format using Hugo's familiar templating"
[page-resources]: /content-management/page-resources/
[pagevars]: /variables/page/
[section]: /content-management/sections/
[taxweight]: /content-management/taxonomies/
[toml]: https://github.com/toml-lang/toml "Specification for TOML, Tom's Obvious Minimal Language"
[urls]: /content-management/urls/
[variables]: /variables/
[yaml]: https://yaml.org/spec/ "Specification for YAML, YAML Ain't Markup Language"
