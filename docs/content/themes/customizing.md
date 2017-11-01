---
title: Customize a Theme
linktitle: Customize a Theme
description: Customize a theme by overriding theme layouts and static assets in your top-level project directories.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [themes]
keywords: [themes, source, organization, directories]
menu:
  docs:
    parent: "themes"
    weight: 20
weight: 20
sections_weight: 20
draft: false
aliases: [/themes/customize/]
toc: true
wip: true
---

The following are key concepts for Hugo site customization with themes. Hugo permits you to supplement *or* override any theme template or static file with files in your working directory.

{{% note %}}
When you use a theme cloned from its git repository, do not edit the theme's files directly. Instead, theme customization in Hugo is a matter of *overriding* the templates made available to you in a theme. This provides the added flexibility of tweaking a theme to meet your needs while staying current with a theme's upstream.
{{% /note %}}

## Override Static Files

There are times where you want to include static assets that differ from versions of the same asset that ships with a theme.

For example, a theme may use jQuery 1.8 in the following location:

```
/themes/<THEME>/static/js/jquery.min.js
```

You want to replace the version of jQuery that ships with the theme with the newer `jquery-3.1.1.js`. The easiest way to do this is to replace the file *with a file of the same name* in the same relative path in your project's root. Therefore, change `jquery-3.1.1.js` to `jquery.min.js` so that it is *identical* to the theme's version and place the file here:

```
/static/js/jquery.min.js
```

## Override Template Files

Anytime Hugo looks for a matching template, it will first check the working directory before looking in the theme directory. If you would like to modify a template, simply create that template in your local `layouts` directory.

The [template lookup order][lookup] explains the rules Hugo uses to determine which template to use for a given piece of content. Read and understand these rules carefully.

This is especially helpful when the theme creator used [partial templates][partials]. These partial templates are perfect for easy injection into the theme with minimal maintenance to ensure future compatibility.

For example:

```
/themes/<THEME>/layouts/_default/single.html
```

Would be overwritten by

```
/layouts/_default/single.html
```

{{% warning %}}
This only works for templates that Hugo "knows about" (i.e., that follow its convention for folder structure and naming). If a theme imports template files in a creatively named directory, Hugo won’t know to look for the local `/layouts` first.
{{% /warning %}}

## Override Archetypes

If the archetype that ships with the theme for a given content type (or all content types) doesn’t fit with how you are using the theme, feel free to copy it to your `/archetypes` directory and make modifications as you see fit.

{{% warning "Beware of `layouts/_default`" %}}
The `_default` directory is a very powerful force in Hugo, especially as it pertains to overwriting theme files. If a default file is located in the local [archetypes](/content-management/archetypes/) or layout directory (i.e., `archetypes/default.md` or `/layouts/_default/*.html`, respectively), it will override the file of the same name in the corresponding theme directory (i.e., `themes/<THEME>/archetypes/default.md` or `themes/<THEME>/layout/_defaults/*.html`, respectively).

It is usually better to override specific files; i.e. rather than using `layouts/_default/*.html` in your working directory.
{{% /warning %}}

[archetypes]: /content-management/archetypes/
[lookup]: /templates/lookup-order/
[partials]: /templates/partials/