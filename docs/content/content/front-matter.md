+++
title = "Front Matter"
date = "2013-07-01"
aliases = ["/doc/front-matter/"]
weight = 40

[menu.main]
    parent = "content"
+++

The front matter is one of the features that gives Hugo its strength. It enables
you to include the meta data of the content right with it. Hugo supports a few
different formats each with their own identifying tokens.

Supported formats: <br>
  **YAML**, identified by '\-\-\-'. <br>
  **TOML**, indentified with '+++'.<br>
  **JSON**, a single JSON object which is surrounded by '{' and '}' each on their own line.

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
**Field names are case insensitive.**

### Required

* **title** The title for the content
* **description** The description for the content
* **date** The date the content will be sorted by
* **indexes** These will use the field name of the plural form of the index (see tags and categories above)

### Optional

* **redirect** Mark the post as a redirect post
* **draft** If true the content will not be rendered unless hugo is called with --build-drafts
* **type** The type of the content (will be derived from the directory automatically if unset)
* **markup** (Experimental) Specify "rst" for reStructuredText (requires
            `rst2html`,) or "md" (default) for the Markdown
* **slug** The token to appear in the tail of the url
   *or*<br>
* **url** The full path to the content from the web root.<br>

*If neither slug or url is present the filename will be used.*

