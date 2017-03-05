---
title: Creating a Theme
linktitle: Creating a Theme
description: The `hugo new theme` command will scaffold the beginnings of a new theme for you to get you on your way.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [themes]
tags: [themes, source, organization, directories]
weight: 30
draft: false
aliases: [/themes/creation/]
toc: true
wip: true
---

{{% warning "Use Relative Links" %}}
If you're creating a them with plans to share it with the community, it is not always safe to assume that the end user of the theme is working from the root of their published website. This is especially important for assets such as stylesheets. See [relURL](/functions/relurl) and [absURL](/functions/absurl).
{{% /warning %}}

Hugo has the ability to create a new theme in your themes directory for you
using the `hugo new` command.

`hugo new theme [name]`

This command will initialize all of the files and directories a basic theme
would need.

## Theme Components

A theme consists of templates and static assets such as javascript and css
files. Themes can also optionally provide [archetypes](/content/archetypes/),
which are archetypal content types used by the `hugo new` command to scaffold new conte files with preconfigured front matter.


{{% note "Use the Hugo Generator Tag" %}}
The [`.Hugo.Generator`](/variables/other/) tag is included in all themes featured in the [Hugo Themes Showcase](http://themes.gohugo.io). We ask that you include the generator tag in all sites and themes you create with Hugo to help the core team track Hugo's usage and popularity.
{{% /note %}}

## Layouts

Hugo is built around the concept that things should be as simple as possible.
Fundamentally, website content is displayed in two different ways, a single
piece of content and a list of content items. With Hugo, a theme layout starts
with the defaults. As additional layouts are defined, they are used for the
content type or section they apply to. This keeps layouts simple, but permits
a large amount of flexibility.

## Single Content

The default single file layout is located at `layouts/_default/single.html`.

## List of Contents

The default list file layout is located at `layouts/_default/list.html`.

## Partial Templates

Theme creators should liberally use [partial templates](/templates/partials/)
throughout their theme files. Not only is a good DRY practice to include shared
code, but partials are a special template type that enables the themes end user
to be able to overwrite just a small piece of a file or inject code into the
theme from their local /layouts. These partial templates are perfect for easy
injection into the theme with minimal maintenance to ensure future
compatibility.

## Static

Everything in the static directory will be copied directly into the final site
when rendered. No structure is provided here to enable complete freedom. It is
common to organize the static content into:

```
/css
/js
/img
```

The actual structure is entirely up to you, the theme creator, on how you would like to organize your files.

## Archetypes

If your theme makes use of specific keys in the front matter, it is a good idea
to provide an archetype for each content type you have. [Read more about archetypes][archetypes].

[archetypes]: /content-management/archetypes/