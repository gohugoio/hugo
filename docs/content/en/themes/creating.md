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
If you're creating a theme with plans to share it on the [Hugo Themes website](https://themes.gohugo.io/) please note the following: 
- If using inline styles you will need to use absolute URLs, for the linked assets to be served properly, e.g. `<div style="background: url('{{ "images/background.jpg" | absURL }}')">`
- Make sure not to use a forward slash `/` in the beginning of a `URL`, because it will point to the host root. Your theme's demo will be available in a subdirectory of the Hugo website and in this scenario Hugo will not generate the correct `URL` for theme assets.
- If using external CSS and JS from a CDN, make sure to load these assets over `https`. Please do not use relative protocol URLs in your theme's templates.
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


