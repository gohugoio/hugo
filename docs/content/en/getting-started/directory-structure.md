---
title: Directory Structure
linktitle: Directory Structure
description: Hugo's CLI scaffolds a project directory structure and then takes that single directory and uses it as the input to create a complete website.
date: 2017-01-02
publishdate: 2017-02-01
lastmod: 2017-03-09
categories: [getting started,fundamentals]
keywords: [source, organization, directories]
menu:
  docs:
    parent: "getting-started"
    weight: 50
weight: 50
sections_weight: 50
draft: false
aliases: [/overview/source-directory/]
toc: true
---

## New Site Scaffolding

{{< youtube sB0HLHjgQ7E >}}

Running the `hugo new site` generator from the command line will create a directory structure with the following elements:

```
.
├── archetypes
├── config.toml
├── content
├── data
├── layouts
├── static
└── themes
```


## Directory Structure Explained

The following is a high-level overview of each of the directories with links to each of their respective sections within the Hugo docs.

[`archetypes`](/content-management/archetypes/)
: You can create new content files in Hugo using the `hugo new` command.
By default, Hugo will create new content files with at least `date`, `title` (inferred from the file name), and `draft = true`. This saves time and promotes consistency for sites using multiple content types. You can create your own [archetypes][] with custom preconfigured front matter fields as well.

[`assets`][]
: Stores all the files which need be processed by [Hugo Pipes]({{< ref "/hugo-pipes" >}}). Only the files whose `.Permalink` or `.RelPermalink` are used will be published to the `public` directory. Note: assets directory is not created by default.

[`config`](/getting-started/configuration/)
: Hugo ships with a large number of [configuration directives](https://gohugo.io/getting-started/configuration/#all-variables-yaml).
The [config directory](/getting-started/configuration/#configuration-directory) is where those directives are stored as JSON, YAML, or TOML files. Every root setting object can stand as its own file and structured by environments.
Projects with minimal settings and no need for environment awareness can use a single `config.toml` file at its root.

Many sites may need little to no configuration, but Hugo ships with a large number of [configuration directives][] for more granular directions on how you want Hugo to build your website. Note: config directory is not created by default.

[`content`][]
: All content for your website will live inside this directory. Each top-level folder in Hugo is considered a [content section][]. For example, if your site has three main sections---`blog`, `articles`, and `tutorials`---you will have three directories at `content/blog`, `content/articles`, and `content/tutorials`. Hugo uses sections to assign default [content types][].

[`data`](/templates/data-templates/)
: This directory is used to store configuration files that can be
used by Hugo when generating your website. You can write these files in YAML, JSON, or TOML format. In addition to the files you add to this folder, you can also create [data templates][] that pull from dynamic content.

[`layouts`][]
: Stores templates in the form of `.html` files that specify how views of your content will be rendered into a static website. Templates include [list pages][lists], your [homepage][], [taxonomy templates][], [partials][], [single page templates][singles], and more.

[`static`][]
: Stores all the static content: images, CSS, JavaScript, etc. When Hugo builds your site, all assets inside your static directory are copied over as-is. A good example of using the `static` folder is for [verifying site ownership on Google Search Console][searchconsole], where you want Hugo to copy over a complete HTML file without modifying its content.

{{% note %}}
From **Hugo 0.31** you can have multiple static directories.
{{% /note %}}

resources
: Caches some files to speed up generation. Can be also used by template authors to distribute built SASS files, so you don't have to have the preprocessor installed. Note: resources directory is not created by default.


[archetypes]: /content-management/archetypes/
[configuration directives]: /getting-started/configuration/#all-variables-yaml
[`content`]: /content-management/organization/
[content section]: /content-management/sections/
[content types]: /content-management/types/
[data templates]: /templates/data-templates/
[homepage]: /templates/homepage/
[`layouts`]: /templates/
[`static`]: /content-management/static-files/
[lists]: /templates/list/
[pagevars]: /variables/page/
[partials]: /templates/partials/
[searchconsole]: https://support.google.com/analytics/answer/1142414?hl=en
[singles]: /templates/single-page-templates/
[starters]: /tools/starter-kits/
[taxonomies]: /content-management/taxonomies/
[taxonomy templates]: /templates/taxonomy-templates/
[types]: /content-management/types/
[`assets`]: {{< ref "/hugo-pipes/introduction#asset-directory" >}}
