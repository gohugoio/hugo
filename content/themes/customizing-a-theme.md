---
title: Customizing a Theme
linktitle: Customizing a Theme
description: Customize a theme by overriding theme layouts and static assets in your top-level project directories.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [themes]
tags: [themes, source, organization, directories]
weight: 20
draft: false
aliases: [/themes/customizing/,/themes/customize/]
toc: true
needsreview: true
---

## Introduction

The following are key concepts for Hugo site customization. Hugo permits you to supplement *or* override any theme template or static file, with files in your working directory.

{{% note %}}
When you use a theme cloned from its git repository, do not edit the theme's files directly. Instead, theme customization in Hugo is a matter of *overriding* the templates made available to you in a theme. This provides the added flexibility of tweaking a theme to meet your needs while staying current with a theme's upstream.
{{% /note %}}

## Understanding the Theme Lookup Order

{{< readfile file="readfile-content/lookupexplanation.md" markdown="true" >}}

## Replacing Static Files

There are times where you want to include static assets that differ from versions of the same asset that ships with a theme. For example, if you would like to use a more recent version of jQuery than what the theme happens to include, simply place an identically-named file in the same relative location but in your working directory.

Let's assume the theme you are using has jQuery 1.8 in the following location:

```bash
/themes/<THEME>/static/js/jquery.min.js
```

You want to replace jQuery with jQuery 1.7. The easiest way to do this is to replace the file *with a file of the same name* in the same relative path in your project's root.

So, to replace jQuery 1.7 from the theme, take your version of jQuery (e.g., `jquery-3.1.1.js`), change the file name so that it is *identical* to the theme file you are trying to use (`jquery.min.js`) and place it here:

```bash
/static/js/jquery.min.js
```

## Replacing Template Files

Anytime Hugo looks for a matching template, it will first check the working directory before looking in the theme directory. If you would like to modify a template, simply create that template in your local `layouts` directory.

In the [template documentation](/templates/overview/) _each different template type explains the rules it uses to determine which template to use_. Read and understand these rules carefully.

This is especially helpful when the theme creator used [partial templates](/templates/partials/). These partial templates are perfect for easy injection into the theme with minimal maintenance to ensure future compatibility.

For example:

```bash
/themes/themename/layouts/_default/single.html
```

Would be overwritten by

```bash
/layouts/_default/single.html
```

{{% warning %}}
This only works for templates that Hugo "knows about" (i.e., that follow its convention for folder structure and naming). If a theme imports template files in a creatively named directory, Hugo won’t know to look for the local `/layouts` first.
{{% /warning %}}

## Replace an Archetype

If the archetype that ships with the theme for a given content type (or all content types) doesn’t fit with how you are using the theme, feel free to copy it to your `/archetypes` directory and make modifications as you see fit.

{{% warning "Beware of `layouts/_default`" %}}
The `_default` directory is a very powerful force in Hugo, especially as it pertains to overwriting theme files. If a default file is located in the local [archetype](/content-management/archetypes/) or layout directory (i.e., `archetypes/default.md` or `/layouts/_default/*.html`, respectively), it will override the file of the same name in the corresponding theme directory (i.e., `themes/<THEME>/archetypes/default.md` or `themes/<THEME>/layout/_defaults/*.html`, respectively).

It is usually better to override specific files; i.e. rather than using `layouts/_default/*.html` in your working directory.
{{% /warning %}}