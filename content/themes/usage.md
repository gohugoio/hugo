+++
title = "Using a Theme"
weight = 30
date = 2014-05-12T10:09:27Z
prev = "/themes/installing"
next = "/themes/customizing"

[menu]
  [menu.main]
    parent = "themes"
+++

Please make certain you have installed the themes you want to use in the
/themes directory.

To use a theme for a site:

    hugo -t ThemeName

The ThemeName must match the name of the directory inside /themes

Hugo will then apply the theme first, then apply anything that is in the local
directory. To learn more, goto [customizing themes](/themes/customizing)
