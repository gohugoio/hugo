---
title: Code Toggle
description: Code Toggle tryout and showcase.
date: 2018-03-16
categories: [getting started,fundamentals]
keywords: [configuration,toml,yaml,json]
weight: 60
sections_weight: 60
draft: false
toc: true
---

## The Config Toggler!

This is an example for the Config Toggle shortcode.
Its purpose is to let users choose a Config language by clicking on its corresponding tab. Upon doing so, every Code toggler on the page will be switched to the target language. Also, target language will be saved in user's `localStorage` so when they go to a different pages, Code Toggler display their last "toggled" config language.

{{% note %}}
The `code-toggler` shortcode is not an internal Hugo shortcode. This page's purpose is to test out a custom feature that we use throughout this site. See: https://github.com/gohugoio/gohugoioTheme/blob/master/layouts/shortcodes/code-toggle.html
{{% /note %}}

## That Config Toggler

{{< code-toggle file="config">}}

baseURL: "https://yoursite.example.com/"
title: "My Hugo Site"
footnoteReturnLinkContents: "↩"
permalinks:
  posts: /:year/:month/:title/
params:
  Subtitle: "Hugo is Absurdly Fast!"
  AuthorName: "Jon Doe"
  GitHubUser: "spf13"
  ListOfFoo:
    - "foo1"
    - "foo2"
  SidebarRecentLimit: 5
{{< /code-toggle >}}

## Another Config Toggler!

{{< code-toggle file="theme">}}

# theme.toml template for a Hugo theme

name = "Hugo Theme"
license = "MIT"
licenselink = "https://github.com/budparr/gohugo.io/blob/master/LICENSE.md"
description = ""
homepage = "https://github.com/budparr/gohugo.io"
tags = ["website"]
features = ["", ""]
min_version = 0.18

[author]
  name = "Bud Parr"
  homepage = "https://github.com/budparr"

{{< /code-toggle >}}

## Two regular code blocks

{{< code file="bf-config.toml" >}}
[blackfriday]
  angledQuotes = true
  fractions = false
  plainIDAnchors = true
  extensions = ["hardLineBreak"]
{{< /code >}}

{{< code file="bf-config.yml" >}}
blackfriday:
  angledQuotes: true
  fractions: false
  plainIDAnchors: true
  extensions:
    - hardLineBreak
{{< /code >}}
