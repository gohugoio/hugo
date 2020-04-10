---
title: hugo
linktitle: hugo
description: The `hugo` function provides easy access to Hugo-related data.
godocref:
date: 2019-01-31
publishdate: 2019-01-31
lastmod: 2019-01-31
keywords: []
categories: [functions]
menu:
  docs:
    parent: "functions"
toc:
signature: ["hugo"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---
  
`hugo` returns an instance that contains the following functions:

hugo.Generator
: `<meta>` tag for the version of Hugo that generated the site. `hugo.Generator` outputs a *complete* HTML tag; e.g. `<meta name="generator" content="Hugo 0.63.2" />`

hugo.Version
: the current version of the Hugo binary you are using e.g. `0.63.2`

  
`hugo` returns an instance that contains the following functions:

hugo.Environment
: the current running environment as defined through the `--environment` cli tag

hugo.CommitHash
: the git commit hash of the current Hugo binary e.g. `0e8bed9ccffba0df554728b46c5bbf6d78ae5247`

hugo.BuildDate
: the compile date of the current Hugo binary formatted with RFC 3339 e.g. `2002-10-02T10:00:00-05:00`

hugo.IsProduction
: returns true if `hugo.Environment` is set to the production environment

{{% note "Use the Hugo Generator Tag" %}}
We highly recommend using `hugo.Generator` in your website's `<head>`. `hugo.Generator` is included by default in all themes hosted on [themes.gohugo.io](https://themes.gohugo.io). The generator tag allows the Hugo team to track the usage and popularity of Hugo.
{{% /note %}}

