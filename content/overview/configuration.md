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

The following is an example of a toml config file with some of the default values. Values under `[params]` will populate the `.Site.Params` variable for use in templates:

    contentdir = "content"
    layoutdir = "layouts"
    publishdir = "public"
    builddrafts = false
    baseurl = "http://yoursite.example.com/"
    canonifyurls = true

    [indexes]
      category = "categories"
      tag = "tags"
       
    [params]
      description = "Tesla's Awesome Hugo Site"
      author = "Nikola Tesla"

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

## Configure Blackfriday rendering

[Blackfriday](https://github.com/russross/blackfriday) is the [Markdown](http://daringfireball.net/projects/markdown/) rendering engine used in Hugo. The Blackfriday configuration in Hugo is mostly a set of sane defaults that should fit most use cases.

But Hugo does expose some options -- in the table below matched with the corresponding flag in the  [Blackfriday source](https://github.com/russross/blackfriday/blob/master/html.go):


Flag | Default | Blackfriday flag | Purpose
--- | --- | --- | ---
angledQuotes | false | HTML_SMARTYPANTS_ANGLED_QUOTES |  Enable angled double quotes (`« »`)
documentIDAnchor | true | FootnoteAnchorPrefix and HeaderIDSuffix |  Enable the prepending / appending of the unique document ID to the footnote and header anchor IDs

**Note** that these flags must be grouped under the `blackfriday` key and can be set on **both site and page level**. If set on page, it will override the site setting.

```
blackfriday:
  angledQuotes = true
  documentIDAnchor = false
```

## Notes

Config changes do not reflect with [Live Reload](/extras/livereload).

Please restart `hugo server --watch` whenever you make a config change.
