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

Hugo themes are located in a centralized github repository. [Hugo Themes
Repo](http://github.com/spf13/hugoThemes) itself is really a meta
repository which contains pointers to set of contributed themes.

## Installing all themes

If you would like to install all of the available hugo themes, simply
clone the entire repository from within your working directory.

    git clone --recursive https://github.com/spf13/hugoThemes.git themes


## Installing a specific theme

    mkdir themes
    cd themes
    git clone URL_TO_THEME
