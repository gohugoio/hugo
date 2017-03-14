---
title: Front Matter
linktitle:
description: Hugo allows you to add front matter in yaml, toml, or json to your content files.
date: 2017-01-09
publishdate: 2017-01-09
lastmod: 2017-01-09
categories: [content management]
tags: ["front matter", "yaml", "toml", "json", "metadata", "archetypes"]
weight: 30
draft: false
aliases: [/content/front-matter/]
toc: true
---

**Front matter** allows you to keep metadata attached to an instance of a [content type][]---i.e., embedded inside a content file---and is one of the many features that gives Hugo its strength.

## Front Matter Formats

Hugo supports three formats for front matter, each with their own identifying tokens.

TOML
: identified by opening and closing `+++`.

YAML
: identified by opening and closing `---` *or* opening `---` and closing `...`

JSON
: a single JSON object surrounded by '`{`' and '`}`', each on their own line.

### TOML Example

```toml
+++
title = "spf13-vim 3.0 release and new website"
description = "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
tags = [ ".vimrc", "plugins", "spf13-vim", "vim" ]
date = "2012-04-06"
categories = [
  "Development",
  "VIM"
]
slug = "spf13-vim-3-0-release-and-new-website"
+++
```

### YAML Example

```yaml
---
title: "spf13-vim 3.0 release and new website"
description: "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
tags: [ ".vimrc", "plugins", "spf13-vim", "vim" ]
lastmod: 2015-12-23
date: "2012-04-06"
categories:
  - "Development"
  - "VIM"
slug: "spf13-vim-3-0-release-and-new-website"
---
```

### JSON Example

```json
{
    "title": "spf13-vim 3.0 release and new website",
    "description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
    "tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ],
    "date": "2012-04-06",
    "categories": [
        "Development",
        "VIM"
    ],
    "slug": "spf13-vim-3-0-release-and-new-website"
}
```

## Front Matter Variables

### Predefined

There are a few predefined variables that Hugo is aware of. See [Page Variables][pagevars] for how to call many of these predefined variables in your templates.

`aliases`
: an array of one or more aliases (e.g. old published path of a renamed content) that would be created to redirect to this content. See [Aliases][aliases] for details.

`date`
: the datetime at which the content was created; note this value is auto-populated according to Hugo's built-in [archetype][]

`description`
: the description for the content.

`draft`
: If true, the content will not be rendered unless `hugo` is called with `--buildDrafts`

`expirydate`
: denotes when content should no longer be published by Hugo; expired will not be rendered unless `hugo` is called with `--buildExpired`

`isCJKLanguage`
: if `true`, explicitly treat the content as CJKLanguage; both `.Summary` and `.WordCount` can work properly in CJK languages.

`keywords`
: the meta keywords for the content.

`layout`
: the layout Hugo should look for in the [lookup order][] when rendering this content.

`lastmod`
: the datetime at which the content was last modified.

`linktitle`
: used for creating links to content; if set, Hugo defaults to using the `linktitle` before the `title`. Hugo can also [order lists of content by `linktitle`][bylinktitle]

`markup`
: **experimental**; specify `"rst"` for reStructuredText (requires`rst2html`) or `"md"` (default) for Markdown

`publishdate`
: If in the future, content will not be rendered unless `hugo` is called with `--buildFuture`

<!-- `sitemap`
: see https://github.com/spf13/hugo/blob/ede452d34ef82a2d6949bf0c5a4584caf3ae03bc/hugolib/page.go#L1062 -->

`slug`
: appears as the tail of the URL. When used, it will override the segment of the URL based on the filename.

<!-- `status`
: see https://github.com/spf13/hugo/blob/ede452d34ef82a2d6949bf0c5a4584caf3ae03bc/hugolib/page.go#L1059 -->

`taxonomies`
: these will use the field name of the plural form of the index (see `tags` and
`categories` in the above front matter examples).

`title`
: the title for the content.

`type`
: the type of the content; this value will be derived from the directory (i.e., the [section][]) automatically if unset.

`url`
: the full path to the content from the web root. It makes no assumptions about the path of the content file. It also ignores any language prefixes of
the multilingual feature.

`weight`
: used for sorting; see [list templates][ordering]

{{% note "Hugo's Default URL Destinations" %}}
If neither `slug` nor `url` is present, and [permalinks are not configured otherwise](/content-management/urls/#permalinks), the filename will be used to create the URL for a page. See [Content Organization](/content-management/organization) and [URL Management](/content-management/urls/).
{{% /note %}}

### User-Defined

You can add fields to your front matter arbitrarily to meet your needs. These user-defined key-values are placed into a single `.Params` variable for use in your templates:

```yaml
include_toc: true
show_comments: false
```

These two user-defined fields can then be accessed via `.Params.include_toc` and `.Params.show_comments`, respectively. The [Variables][] section provides more information on using Hugo's page- and site-level variables in your templates.

{{% note %}}
Field names are always normalized to lowercase; e.g., `camelCase: true` is available as `.Params.camelcase`.
{{% /note %}}

## Ordering Content Through Front Matter

You can assign content-specific `weight` in the front matter of your content. These values are especially useful for [ordering][] in list views. You can use `weight` for ordering of content and the convention of [`<TAXONOMY>_weight`][taxweight] for ordering content within a taxonomy. See [Ordering and Grouping Hugo Lists][ordering] to see the role `weight` in ordering content.

## Overriding Global Blackfriday Configuration

It's possible to set some options for Markdown rendering in a content's front matter as an override to the options set in your site `config`.

See [Configuration][config] for more information on configuring Blackfriday rendering.

## Front Matter Format Specs

* [TOML Spec][]
* [YAML Spec][]
* [JSON Spec][]

[aliases]: /content-management/urls/#aliases/
[archetype]: /content-management/archetypes/
[bylinktitle]: /templates/lists/#by-link-title
[config]: /getting-started/configuration/ "Hugo documentation for site configuration"
[contentorg]: /content-management/organization/
[content type]: /content-management/types/
[JSON Spec]: /documents/ecma-404-json-spec.pdf "Specification for JSON, JavaScript Object Notation"
[lookup]: /content-management/l
[ordering]: /templates/lists/ "Hugo provides multiple ways to sort and order your content in list templates"
[pagevars]: /variables/page/
[section]: /content-management/sections/
[taxweight]: /content-management/taxonomies/
[TOML Spec]: https://github.com/toml-lang/toml "Specification for TOML, Tom's Obvious Minimal Language"
[urls]: /content-management/urls/
[YAML Spec]: http://yaml.org/spec/ "Specification for YAML, YAML Ain't Markup Language"
