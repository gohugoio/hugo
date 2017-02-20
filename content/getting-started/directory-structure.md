---
title: Directory Structure
linktitle: Directory Structure
description: Explanation of the directory structure in a typical Hugo project and how Hugo traverses the file system therein.
date: 2017-01-02
publishdate: 2017-01-02
lastmod: 2017-01-02
categories: [project organization]
tags: [source, organization, directories,fundamentals]
weight: 50
draft: false
aliases: [/overview/source-directory/]
notesforauthors:
---

<!-- copied from old version of quick start -->

* **archetypes**: You can create new content files in Hugo using the `hugo new` command. When you run that command, it adds few configuration properties to the post like date and title. [Archetype][archetypes] allows you to define your own configuration properties that will be added to front matter of new content files whenever `hugo new` command is used.

* **config.toml**: Every website should have a configuration file at the root. By default, the configuration file uses `TOML` format but you can also use `YAML` or `JSON` formats as well. [TOML](https://github.com/toml-lang/toml) is minimal configuration file format that's easy to read due to obvious semantics. The configuration settings mentioned in the `config.toml` are applied to the full site. These configuration settings include `baseURL` and `title` of the website.

* **content**: This is where you will store content of the website. Inside content, you will create sub-directories for different sections. Let's suppose your website has three actions -- `blog`, `article`, and `tutorial` then you will have three different directories for each of them inside the `content` directory. The name of the section i.e. `blog`, `article`, or `tutorial` will be used by Hugo to apply a specific layout applicable to that section.

* **data**: This directory is used to store configuration files that can be
used by Hugo when generating your website. You can write these files in YAML, JSON, or TOML format.

* **layouts**: The content inside this directory is used to specify how your content will be converted into the static website.

* **static**: This directory is used to store all the static content that your website will need like images, CSS, JavaScript or other static content.
