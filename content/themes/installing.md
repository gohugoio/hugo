---
lastmod: 2015-10-10
date: 2014-05-12T10:09:49Z
menu:
  main:
    parent: themes
next: /themes/usage
prev: /themes/overview
title: Installing Themes
weight: 20
---

Community-contributed [Hugo themes](http://themes.gohugo.io/), showcased
at [themes.gohugo.io](//themes.gohugo.io/), are hosted in a centralized
GitHub repository.  The [Hugo Themes Repo](https://github.com/gohugoio/hugoThemes)
itself at [github.com/gohugoio/hugoThemes](https://github.com/gohugoio/hugoThemes) is
really a meta repository which contains pointers to set of contributed themes.

## Installing all themes

If you would like to install all of the available Hugo themes, simply
clone the entire repository from within your working directory. Depending 
on your internet connection the download of all themes might take a while.

> **NOTE:** Make sure you've installed [Git](https://git-scm.com/) on your computer.
Otherwise you will not be able to clone the theme repositories.

```bash
$ git clone --depth 1 --recursive https://github.com/gohugoio/hugoThemes.git themes
```
Before you use a theme, remove the .git folder in that theme's root folder. Otherwise, this will cause problem if you deploy using Git.

## Installing a specific theme

Switch into the `themes` directory and download a theme by replacing `URL_TO_THEME` 
with the URL of the theme repository, e.g. `https://github.com/spf13/hyde`:

    $ cd themes
    $ git clone URL_TO_THEME
    
Alternatively, you can download the theme as `.zip` file and extract it in the 
`themes` directory.

**NOTE:** Please have a look at the `README.md` file that is shipped with all themes.
It might contain further instructions that are required to setup the theme, e.g. copying
an example configuration file.
