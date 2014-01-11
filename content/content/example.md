---
title: "Example Content File"
date: "2013-07-01"
aliases: ["/doc/example/"]
linktitle: "Example"
groups: ['content']
groups_weight: 50
---

Somethings are better shown than explained. The following is a very basic example of a content file:

**mysite/project/nitro.md  <- http://mysite.com/project/nitro.html**

{{% highlight yaml %}}
---
Title:       "Nitro : A quick and simple profiler for golang"
Description: ""
Keywords:    [ "Development", "golang", "profiling" ]
Tags:        [ "Development", "golang", "profiling" ]
date:        "2013-06-19"
Topics:      [ "Development", "GoLang" ]
Slug:        "nitro"
project_url: "http://github.com/spf13/nitro"
---

# Nitro

Quick and easy performance analyzer library for golang.

## Overview

Nitro is a quick and easy performance analyzer library for golang.
It is useful for comparing A/B against different drafts of functions
or different functions.

## Implementing Nitro

Using Nitro is simple. First use go get to install the latest version
of the library.

    $ go get github.com/spf13/nitro

Next include nitro in your application.

{{% /highlight %}}

