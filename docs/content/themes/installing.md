---
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
GitHub repository.  The [Hugo Themes Repo](https://github.com/spf13/hugoThemes)
itself at [github.com/spf13/hugoThemes](https://github.com/spf13/hugoThemes) is
really a meta repository which contains pointers to set of contributed themes.

## Installing all themes

If you would like to install all of the available Hugo themes, simply
clone the entire repository from within your working directory:

```bash
$ git clone --recursive https://github.com/spf13/hugoThemes.git themes
```

## Installing a specific theme

    $ mkdir themes
    $ cd themes
    $ git clone URL_TO_THEME
