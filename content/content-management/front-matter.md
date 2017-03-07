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
wip: true
toc: true
notesforauthors:
---

**Front matter** allows you to keep metadata attached to an instance of a [content type][]---i.e., embedded inside a content file---and is one of the many features that gives Hugo its strength. Hugo supports a few different formats for front matter, each with their own identifying tokens.

## Supported Front Matter Formats

* TOML. identified by opening and closing `+++`.
* YAML. identified by opening and closing `---` *or* opening `---` and closing `...`
* JSON. a single JSON object which is surrounded by '`{`' and '`}`', each on their own line.

### TOML Front Matter Example

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

### YAML Front Matter Example

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

### JSON Front Matter Example

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

## Variables

There are a few predefined variables that Hugo is aware of and utilizes. The user can also create any variable they want. These will be placed into the `.Params` variable available to the templates.

{{% note %}}
Field names are always normalized to lowercase (e.g., `camelCase: true` is available as `.Params.camelcase`) for both Hugo's built-in *and* user-defined variables.
{{% /note %}}

### Required variables

`title`
: The title for the content

`description`
: The description for the content

`date`
: The date the content will be sorted by

`taxonomies`
: These will use the field name of the plural form of the index (see tags and
categories above)

### Optional variables

`aliases`
: An array of one or more aliases (e.g. old published path of a renamed content) that would be created to redirect to this content. See [Aliases][] for details.

`draft`
: If true, the content will not be rendered unless `hugo` is called with `--buildDrafts`

`expirydate`
: Content already expired will not be rendered unless `hugo` is called with `--buildExpired`

`isCJKLanguage`
: If true, explicitly treat the content as CJKLanguage (`.Summary` and `.WordCount` can work properly in CJKLanguage)

`markup`
: **Experimental**; specify `"rst"` for reStructuredText (requires`rst2html`) or `"md"` (default) for Markdown

`publishdate`
: If in the future, content will not be rendered unless `hugo` is called with `--buildFuture`

`slug`
: appears as tail of the url. It can be used to change the part of the url that is based on the filename

`type`
: The type of the content (will be derived from the directory automatically if unset)

`url`
: The full path to the content from the web root. It makes no assumptions about the path of the content file. It also ignores any language prefixes of
the multilingual feature.

weight
: Used for sorting

{{% note "Hugo's Default URLs" %}}
If neither `slug` nor `url` is present, and [permalinks are not configured otherwise](/content-management/urls/#permalinks), the filename will be used to create the URL for a page. See [Content Organization][contentorg] and [URL Management][urls].
{{% /note %}}

## User-Defined Variables

In addition to....

## Override Global Blackfriday Configuration

It's possible to set some options for Markdown rendering in the page's front matter as an override to the site-wide configuration.

See [site configuration][config] for more information on setting up global Blackfriday options.

User-Define Parameter (`.Params`)

## Assigning `weight` for Ordering


### Assigning `weight` for Ordering Content

### Assigning `weight` for Ordering Taxonomies

## Front Matter Format Specs

* [TOML Spec][]
* [YAML Spec][]
* [JSON Spec][]

[aliases]: /content-management/urls/#aliases/
[config]: /getting-started/configuration/ "Hugo documentation for site configuration"
[contentorg]: /content-management/organization/
[content type]: /content-management/types/
[JSON Spec]: /documents/ecma-404-json-spec.pdf "Specification for JSON, JavaScript Object Notation"
[TOML Spec]: https://github.com/toml-lang/toml "Specification for TOML, Tom's Obvious Minimal Language"
[urls]: /content-management/urls/
[YAML Spec]: http://yaml.org/spec/ "Specification for YAML, YAML Ain't Markup Language"
