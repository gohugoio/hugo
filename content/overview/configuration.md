---
aliases:
- /doc/configuration/
date: 2013-07-01
linktitle: Configuration
menu:
  main:
    parent: getting started
next: /overview/source-directory
notoc: true
prev: /overview/usage
title: Configuring Hugo
weight: 40
---

The directory structure and templates provide the majority of the
configuration for a site. In fact, a config file isn't even needed for many
websites since the defaults follow commonly used patterns.

Hugo expects to find the config file in the root of the source directory and
will look there first for a `config.toml` file. If none is present, it will
then look for a `config.yaml` file, followed by a `config.json` file.

The config file is a site-wide config. The config file provides directions to
hugo on how to build the site as well as site-wide parameters and menus.

## Examples

The following is an example of a typical yaml config file:

    ---
    baseurl: "http://yoursite.example.com/"
    ...

The following is an example of a toml config file with some of the default values:

    contentdir = "content"
    layoutdir = "layouts"
    publishdir = "public"
    builddrafts = false
    baseurl = "http://yoursite.example.com/"
    canonifyurls = true

    [indexes]
       category = "categories"
       tag = "tags"

Here is a yaml configuration file which sets a few more options

    ---
    baseurl: "http://yoursite.example.com/"
    title: "Yoyodyne Widget Blogging"
    footnotereturnlinkcontents: "↩"
    permalinks:
      post: /:year/:month/:title/
    params:
      Subtitle: "Spinning the cogs in the widgets"
      AuthorName: "John Doe"
      GitHubUser: "spf13"
      ListOfFoo:
        - "foo1"
        - "foo2"
      SidebarRecentLimit: 5
    ...

## Notes

Config changes do not reflect with [Live Reload](/extras/livereload).

Please restart `hugo server --watch` whenever you make a config change.
