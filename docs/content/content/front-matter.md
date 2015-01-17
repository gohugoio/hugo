---
aliases:
- /doc/front-matter/
date: 2013-07-01
menu:
  main:
    parent: content
next: /content/sections
prev: /content/organization
title: Front Matter
weight: 20
---

The **front matter** is one of the features that gives Hugo its strength. It enables
you to include the meta data of the content right with it. Hugo supports a few
different formats, each with their own identifying tokens.

Supported formats:

  * **[TOML][]**, identified by '`+++`'.
  * **[YAML][]**, identified by '`---`'.
  * **[JSON][]**, a single JSON object which is surrounded by '`{`' and '`}`', each on their own line.

[TOML]: https://github.com/toml-lang/toml "Tom's Obvious, Minimal Language"
[YAML]: http://www.yaml.org/ "YAML Ain't Markup Language"
[JSON]: http://www.json.org/ "JavaScript Object Notation"

### TOML Example

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
    
    Content of the file goes Here

### YAML Example

    ---
    title: "spf13-vim 3.0 release and new website"
    description: "spf13-vim is a cross platform distribution of vim plugins and resources for Vim."
    tags: [ ".vimrc", "plugins", "spf13-vim", "vim" ]
    date: "2012-04-06"
    categories:
      - "Development"
      - "VIM"
    slug: "spf13-vim-3-0-release-and-new-website"
    ---
    
    Content of the file goes Here

### JSON Example

    {
        "title": "spf13-vim 3.0 release and new website",
        "description": "spf13-vim is a cross platform distribution of vim plugins and resources for Vim.",
        "tags": [ ".vimrc", "plugins", "spf13-vim", "vim" ],
        "date": "2012-04-06",
        "categories": [
            "Development",
            "VIM"
        ],
        "slug": "spf13-vim-3-0-release-and-new-website",
    }
    
    Content of the file goes Here

## Variables

There are a few predefined variables that Hugo is aware of and utilizes. The user can also create
any variable they want to. These will be placed into the `.Params` variable available to the templates.
Field names are always normalized to lowercase (e.g. `camelCase: true` is available as `.Params.camelcase`).

### Required variables

* **title** The title for the content
* **description** The description for the content
* **date** The date the content will be sorted by
* **taxonomies** These will use the field name of the plural form of the index (see tags and categories above)

### Optional variables

* **redirect** Mark the post as a redirect post
* **draft** If true, the content will not be rendered unless `hugo` is called with `--buildDrafts`
* **publishdate** If in the future, content will not be rendered unless `hugo` is called with `--buildFuture`
* **type** The type of the content (will be derived from the directory automatically if unset)
* **weight** Used for sorting
* **markup** *(Experimental)* Specify `"rst"` for reStructuredText (requires
            `rst2html`) or `"md"` (default) for Markdown
* **slug** The token to appear in the tail of the URL,
   *or*<br>
* **url** The full path to the content from the web root.<br>

*If neither `slug` or `url` is present, the filename will be used.*

## Configure Blackfriday rendering

It's possible to set some options for Markdown rendering in the page's front matter, as an override to the site wide configuration.

See [Configuration]({{< ref "overview/configuration.md#configure-blackfriday-rendering" >}}) for more.

