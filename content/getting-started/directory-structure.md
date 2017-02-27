---
title: Directory Structure
linktitle: Directory Structure
description: Hugo's CLI scaffolds a project's directory structure nearly instantly and then takes that single directory and uses it as the input for creating a complete website.
date: 2017-01-02
publishdate: 2017-01-02
lastmod: 2017-01-02
categories: [project organization]
tags: [source, organization, directories,fundamentals]
weight: 50
draft: false
needsreview: true
aliases: [/overview/source-directory/]
---

<!-- copied from old overview/source-directory -->

Hugo takes a single directory and uses it as the input for creating a complete
website.

## Directory Scaffolding in New Hugo Sites

The top level of a source directory will typically have the following elements:

```bash
▸ archetypes/
▸ content/
▸ data/
▸ i18n/
▸ layouts/
▸ static/
▸ themes/
  config.toml
```

Learn more about the different directories and what their purpose is:

* [config](/getting-started/configuration/)
* [data](/templates/data-templates/)
* [i18n](/content-management/multilingual/)
* [archetypes](/content-management/archetypes/)
* [content](/content-managemt/organization/)
* [layouts](/templates/)
* [static](/themes/creating-a-theme/)
* [themes](/themes/)


## Example Hugo Project Directory

An example directory may look like the following

```bash
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
|   ├── taxonomies
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
    ├── images
    └── js
```

This directory structure tells us a lot about this website:

1. The website intends to have two different types of content: *posts* and *quotes*.
2. It will also apply two different taxonomies to that content: *categories* and *tags*.
3. It will be displaying content in 3 different views: a list, a summary, and a full-page view.

## Content for Homepage and Other List Pages

Since Hugo 0.18, "everything" is a `Page` that can have content and metadata, like `.Params`, attached to it -- and share the same set of [page variables](/variables/page-variables/).

To add content and front matter to the home page, a section, a taxonomy or a taxonomy terms listing, add a markdown file with the base name `_index` on the relevant place on the file system.

For the default Markdown content, the filename will be `_index.md`.

Se the example directory tree below.

{{% note "You Don't Have to Create `_index.md`"%}}
You don't have to create an `_index` file for every list page (i.e. section, taxonomy, taxonomy terms, etc). If Hugo does not find an `_index.md` on a list page, a default page will be created if not present but with no `{{.Content}}`  and only the default values for `.Title` etc.
{{% /note %}}

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


<!-- copied from old version of quick start -->

* **archetypes**: You can create new content files in Hugo using the `hugo new` command. When you run that command, it adds few configuration properties to the post like date and title. [Archetype][archetypes] allows you to define your own configuration properties that will be added to front matter of new content files whenever `hugo new` command is used.

* **config.toml**: Every website should have a configuration file at the root. By default, the configuration file uses `TOML` format but you can also use `YAML` or `JSON` formats as well. [TOML](https://github.com/toml-lang/toml) is minimal configuration file format that's easy to read due to obvious semantics. The configuration settings mentioned in the `config.toml` are applied to the full site. These configuration settings include `baseURL` and `title` of the website.

* **content**: This is where you will store content of the website. Inside content, you will create sub-directories for different sections. Let's suppose your website has three actions -- `blog`, `article`, and `tutorial` then you will have three different directories for each of them inside the `content` directory. The name of the section i.e. `blog`, `article`, or `tutorial` will be used by Hugo to apply a specific layout applicable to that section.

* **data**: This directory is used to store configuration files that can be
used by Hugo when generating your website. You can write these files in YAML, JSON, or TOML format.

* **layouts**: The content inside this directory is used to specify how your content will be converted into the static website.

* **static**: This directory is used to store all the static content that your website will need like images, CSS, JavaScript or other static content.
