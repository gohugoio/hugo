---
date: 2014-05-12T10:09:34Z
menu:
  main:
    parent: themes
next: /themes/creation
prev: /themes/usage
title: Customizing a Theme
weight: 40
---

Hugo themes permit you to supplement or override any template or file
from within your working directory.

## Replacing Static files

For including a different file than what the theme ships
with. For example, if you would like to use a more recent version of jQuery
than what the theme happens to include, simply place an identically named file in the same
relative location but in your working directory. For example, if the
theme has jQuery 1.6 in `/themes/themename/static/js/jquery.min.js`, simply place your file
in the same relative path `/static/js/jquery.min.js`.

## Replace a single template file

Anytime Hugo looks for a matching template, it will first check the
working directory before looking in the theme directory. If you would
like to modify a template, simply create that template in your local
`layouts` directory. In the [template documentation](/templates/overview/)
each different template type explains the rules it uses to determine
which template to use.

This is especially helpful when the theme creator used [partial
templates](/templates/partials/). These partial templates are perfect for easy
injection into the theme with minimal maintenance to ensure future
compatibility.

**Warning: This only works for templates that Hugo knows about. If the
theme imports template files in a creatively named directory,
Hugo won’t know to look for the local `/layouts` first.**

## Replace an archetype

If the archetype that ships with the theme for a given content type (or
all content types) doesn’t fit with how you are using the theme, feel
free to copy it to your `/archetypes` directory and make modifications as
you see fit.

## Beware of the default

**Default** is a very powerful force in Hugo... Especially as it pertains to
overwriting theme files. If a default is located in the local archetype
directory or `/layouts/_default/` directory, it will be used instead of
any of the similar files in the theme. 

It is usually better to override specific files rather than using the
default in your working directory.
