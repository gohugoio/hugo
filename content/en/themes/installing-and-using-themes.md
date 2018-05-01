---
title: Install and Use Themes
linktitle: Install and Use Themes
description: Install and use a theme from the Hugo theme showcase easily through the CLI.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [themes]
keywords: [install, themes, source, organization, directories,usage]
menu:
  docs:
    parent: "themes"
    weight: 10
weight: 10
sections_weight: 10
draft: false
aliases: [/themes/usage/,/themes/installing/]
toc: true
wip: true
---

{{% note "No Default Theme" %}}
Hugo currently doesn’t ship with a “default” theme. This decision is intentional. We leave it up to you to decide which theme best suits your Hugo project.
{{% /note %}}

## Assumptions

1. You have already [installed Hugo on your development machine][install].
2. You have git installed on your machine and you are familiar with basic git usage.

## Install Themes

{{< youtube L34JL_3Jkyc >}}

The community-contributed themes featured on [themes.gohugo.io](//themes.gohugo.io/) are hosted in a [centralized GitHub repository][themesrepo]. The Hugo Themes Repo at <https://github.com/gohugoio/hugoThemes> is really a meta repository that contains pointers to a set of contributed themes.

{{% warning "Get `git` First" %}}
Without [Git](https://git-scm.com/) installed on your computer, none of the following theme instructions will work. Git tutorials are beyond the scope of the Hugo docs, but [GitHub](https://try.github.io/) and [codecademy](https://www.codecademy.com/learn/learn-git) offer free, interactive courses for beginners.
{{% /warning %}}

### Install All Themes

You can install *all* available Hugo themes by cloning the entire [Hugo Theme repository on GitHub][themesrepo] from within your working directory. Depending on your internet connection the download of all themes might take a while.

```
git clone --depth 1 --recursive https://github.com/gohugoio/hugoThemes.git themes
```

Before you use a theme, remove the .git folder in that theme's root folder. Otherwise, this will cause problem if you deploy using Git.

### Install a Single Theme

Change into the `themes` directory and download a theme by replacing `URL_TO_THEME` with the URL of the theme repository:

```
cd themes
git clone URL_TO_THEME
```

The following example shows how to use the "Hyde" theme, which has its source hosted at <https://github.com/spf13/hyde>:

{{< code file="clone-theme.sh" >}}
cd themes
git clone https://github.com/spf13/hyde
{{< /code >}}

Alternatively, you can download the theme as a `.zip` file, unzip the theme contents, and then move the unzipped source into your `themes` directory.

{{% note "Read the `README`" %}}
Always review the `README.md` file that is shipped with a theme. Often, these files contain further instructions required for theme setup; e.g., copying values from an example configuration file.
{{% /note %}}

## Theme Placement

Please make certain you have installed the themes you want to use in the
`/themes` directory. This is the default directory used by Hugo. Hugo comes with the ability to change the themes directory via the [`themesDir` variable in your site configuration][config], but this is not recommended.

## Use Themes

Hugo applies the decided theme first and then applies anything that is in the local directory. This allows for easier customization while retaining compatibility with the upstream version of the theme. To learn more, go to [customizing themes][customizethemes].

### Command Line

There are two different approaches to using a theme with your Hugo website: via the Hugo CLI or as part of your [site configuration file][config].

To change a theme via the Hugo CLI, you can pass the `-t` [flag][] when building your site:

```
hugo -t themename
```

Likely, you will want to add the theme when running the Hugo local server, especially if you are going to [customize the theme][customizethemes]:

```
hugo server -t themename
```

### `config` File

If you've already decided on the theme for your site and do not want to fiddle with the command line, you can add the theme directly to your [site configuration file][config]:

```
theme: themename
```

{{% note "A Note on `themename`" %}}
The `themename` in the above examples must match the name of the specific theme directory inside `/themes`; i.e., the directory name (likely lowercase and urlized) rather than the name of the theme displayed in the [Themes Showcase site](http://themes.gohugo.io).
{{% /note %}}

[customizethemes]: /themes/customizing/
[flag]: /getting-started/usage/ "See the full list of flags in Hugo's basic usage."
[install]: /getting-started/installing/
[config]: /getting-started/configuration/  "Learn how to customize your Hugo website configuration file in yaml, toml, or json."
[themesrepo]: https://github.com/gohugoio/hugoThemes
