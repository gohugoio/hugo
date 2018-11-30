---
title: Create a Theme
linktitle: Create a Theme
description: The `hugo new theme` command will scaffold the beginnings of a new theme for you to get you on your way.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [themes]
keywords: [themes, source, organization, directories]
menu:
  docs:
    parent: "themes"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: [/themes/creation/,/tutorials/creating-a-new-theme/]
toc: true
wip: true
---

{{% warning "Use Absolute Links" %}}
If you're creating a theme with plans to share it on the [Hugo Themes website](https://themes.gohugo.io/) please note that your theme's demo will be available in a sub-directory of website and for the theme's assets to load properly you will need to create absolute paths in the templates  by using either the [absURL](/functions/absurl) function or `.Permalink`. Also make sure not to use a forward slash `/` in the beginning of a `PATH`, because Hugo will turn it into a relative URL and the `absURL` function will have no effect.
{{% /warning %}}

Hugo can initialize a new blank theme directory within your existing `themes` using the `hugo new` command:

```
hugo new theme [name]
```

## Theme Folders

A theme component can provide files in one or more of the following standard Hugo folders:

layouts
: Templates used to render content in Hugo. Also see [Templates Lookup Order](/templates/lookup-order/).

static
: Static files, such as logos, CSS and JavaScript.

i18n
: Language bundles.

data
: Data files.

archetypes
: Content templates used in `hugo new`.


## Theme Configuration File

A theme component can also provide its own [Configuration File](/getting-started/configuration/), e.g. `config.toml`. There are some restrictions to what can be configured in a theme component, and it is not possible to overwrite settings in the project.

The following settings can be set:

* `params` (global and per language)
* `menu` (global and per language)
* `outputformats` and `mediatypes`


## Theme Description File

In addition to the configuration file, a theme can also provide a `theme.toml` file that describes the theme, the author and origin etc. See [Add Your Hugo Theme to the Showcase](/contribute/themes/).


{{% note "Use the Hugo Generator Tag" %}}
The [`.Hugo.Generator`](/variables/hugo/) tag is included in all themes featured in the [Hugo Themes Showcase](http://themes.gohugo.io). We ask that you include the generator tag in all sites and themes you create with Hugo to help the core team track Hugo's usage and popularity.
{{% /note %}}


