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

<dl>
<dt><code>title</code></dt><dd>The title for the content</dd>
<dt><code>description</code></dt><dd>The description for the content</dd>
<dt><code>date</code></dt><dd>The date the content will be sorted by</dd>
<dt><code>taxonomies</code></dt><dd>These will use the field name of the plural form of the index (see tags and categories above)</dd>
</dl>

### Optional variables

<dl>
<dt><code>redirect</code></dt><dd>Mark the post as a redirect post</dd>
<dt><code>draft</code></dt><dd>If true, the content will not be rendered unless <code>hugo</code> is called with <code>--buildDrafts</code></dd>
<dt><code>publishdate</code></dt><dd>If in the future, content will not be rendered unless <code>hugo</code> is called with <code>--buildFuture</code></dd>
<dt><code>type</code></dt><dd>The type of the content (will be derived from the directory automatically if unset)</dd>
<dt><code>weight</code></dt><dd>Used for sorting</dd>
<dt><code>markup</code></dt><dd><em>(Experimental)</em> Specify <code>"rst"</code> for reStructuredText (requires <code>rst2html</code>) or <code>"md"</code> (default) for Markdown</dd>
<dt><code>slug</code></dt><dd>The token to appear in the tail of the URL, <em>or</em></dd>
<dt><code>url</code></dt><dd>The full path to the content from the web root.<br></dd>
</dl>

*If neither `slug` or `url` is present, the filename will be used.*

## Configure Blackfriday rendering

It's possible to set some options for Markdown rendering in the page's front matter, as an override to the site wide configuration.

See [Configuration]({{< ref "overview/configuration.md#configure-blackfriday-rendering" >}}) for more.

