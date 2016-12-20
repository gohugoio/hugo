---
aliases:
- /doc/front-matter/
lastmod: 2015-12-23
date: 2013-07-01
menu:
  main:
    parent: content
next: /content/sections
prev: /content/organization
title: Front Matter
weight: 20
toc: true
---

The **front matter** is one of the features that gives Hugo its strength. It enables
you to include the meta data of the content right with it. Hugo supports a few
different formats, each with their own identifying tokens.

Supported formats:

  * **[TOML][]**, identified by '`+++`'.
  * **[YAML][]**, identified by '`---`'.
  * **[JSON][]**, a single JSON object which is surrounded by '`{`' and '`}`', followed by a newline.

[TOML]: https://github.com/toml-lang/toml "Tom's Obvious, Minimal Language"
[YAML]: http://www.yaml.org/ "YAML Ain't Markup Language"
[JSON]: http://www.json.org/ "JavaScript Object Notation"

## TOML Example

<pre><code class="language-toml">+++
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
</code><code class="language-markdown">Content of the file goes Here
</code></pre>

## YAML Example

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

Content of the file goes Here
```

## JSON Example

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

Content of the file goes Here
```

## Variables

There are a few predefined variables that Hugo is aware of and utilizes. The user can also create
any variable they want. These will be placed into the `.Params` variable available to the templates.
Field names are always normalized to lowercase (e.g. `camelCase: true` is available as `.Params.camelcase`).

### Required variables

* **title** The title for the content
* **description** The description for the content
* **date** The date the content will be sorted by
* **taxonomies** These will use the field name of the plural form of the index (see tags and categories above)

### Optional variables

* **aliases** An array of one or more aliases
              (e.g. old published path of a renamed content)
              that would be created to redirect to this content.
              See [Aliases]({{< relref "extras/aliases.md" >}}) for details.
* **draft** If true, the content will not be rendered unless `hugo` is called with `--buildDrafts`
* **publishdate** If in the future, content will not be rendered unless `hugo` is called with `--buildFuture`
* **expirydate** Content already expired will not be rendered unless `hugo` is called with `--buildExpired`
* **type** The type of the content (will be derived from the directory automatically if unset)
* **isCJKLanguage** If true, explicitly treat the content as CJKLanguage (`.Summary` and `.WordCount` can work properly in CJKLanguage)
* **weight** Used for sorting
* **markup** *(Experimental)* Specify `"rst"` for reStructuredText (requires
            `rst2html`) or `"md"` (default) for Markdown
* **slug** appears as tail of the url. It can be used to change the part of the url that is based on the filename.
* **url** The full path to the content from the web root. It makes no assumptions about the path of the content file. It also ignores any language prefixes of the multilingual feature.

*If neither `slug` or `url` is present, the filename will be used.*

## Configure Blackfriday rendering

It's possible to set some options for Markdown rendering in the page's front matter as an override to the site wide configuration.

See [Configuration]({{< ref "overview/configuration.md#configure-blackfriday-rendering" >}}) for more.

