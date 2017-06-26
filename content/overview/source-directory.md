---
aliases:
- /doc/source-directory/
lastmod: 2015-02-09
date: 2013-07-01
menu:
  main:
    parent: getting started
next: /content/organization
notoc: true
prev: /overview/configuration
title: Source Organization
weight: 50
---

Hugo takes a single directory and uses it as the input for creating a complete
website.


The top level of a source directory will typically have the following elements:

    ▸ archetypes/
    ▸ content/
    ▸ data/
    ▸ i18n/
    ▸ layouts/
    ▸ static/
    ▸ themes/
      config.toml

Learn more about the different directories and what their purpose is:

* [config]({{< relref "overview/configuration.md" >}})
* [data]({{< relref "extras/datafiles.md" >}})
* [i18n]({{< relref "content/multilingual.md#translation-of-strings" >}})
* [archetypes]({{< relref "content/archetypes.md" >}})
* [content]({{< relref "content/organization.md" >}})
* [layouts]({{< relref "templates/overview.md" >}})
* [static]({{< relref "themes/creation.md#static" >}})
* [themes]({{< relref "themes/overview.md" >}})


## Example

An example directory may look like:

    .
    ├── config.toml
    ├── archetypes
    |   └── default.md
    ├── content
    |   ├── post
    |   |   ├── firstpost.md
    |   |   └── secondpost.md
    |   └── quote
    |   |   ├── first.md
    |   |   └── second.md
    ├── data
    ├── i18n
    ├── layouts
    |   ├── _default
    |   |   ├── single.html
    |   |   └── list.html
    |   ├── partials
    |   |   ├── header.html
    |   |   └── footer.html
    |   ├── taxonomy
    |   |   ├── category.html
    |   |   ├── post.html
    |   |   ├── quote.html
    |   |   └── tag.html
    |   ├── post
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── quote
    |   |   ├── li.html
    |   |   ├── single.html
    |   |   └── summary.html
    |   ├── shortcodes
    |   |   ├── img.html
    |   |   ├── vimeo.html
    |   |   └── youtube.html
    |   ├── index.html
    |   └── sitemap.xml
    ├── themes
    |   ├── hyde
    |   └── doc
    └── static
        ├── css
        └── js

This directory structure tells us a lot about this site:

1. The website intends to have two different types of content: *posts* and *quotes*.
2. It will also apply two different taxonomies to that content: *categories* and *tags*.
3. It will be displaying content in 3 different views: a list, a summary and a full page view.

## Content for home page and other list pages

Since Hugo 0.18, "everything" is a `Page` that can have content and metadata, like `.Params`, attached to it -- and share the same set of [page variables](/templates/variables/).

To add content and frontmatter to the home page, a section, a taxonomy or a taxonomy terms listing, add a markdown file with the base name `_index` on the relevant place on the file system.

For the default Markdown content, the filename will be `_index.md`. 

Se the example directory tree below. 

**Note that you don't have to create `_index` file for every section, taxonomy and similar, a default page will be created if not present, but with no content and default values for `.Title` etc.**

```bash
└── content
    ├── _index.md
    ├── categories
    │   ├── _index.md
    │   └── photo
    │       └── _index.md
    ├── post
    │   ├── _index.md
    │   └── firstpost.md
    └── tags
        ├── _index.md
        └── hugo
            └── _index.md
```
  
