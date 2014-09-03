---
aliases:
- /doc/example/
date: 2013-07-01
linktitle: Example
menu:
  main:
    parent: content
next: /themes/overview
notoc: true
prev: /content/ordering
title: Example Content File
weight: 70
---

Some things are better shown than explained. The following is a very basic example of a content file:

**mysite/project/nitro.md  ← http://mysite.com/project/nitro.html**

    ---
    Title:       "Nitro : A quick and simple profiler for Go"
    Description: "Nitro is a simple profiler for you go lang applications"
    Tags:        [ "Development", "Go", "profiling" ]
    date:        "2013-06-19"
    Topics:      [ "Development", "Go" ]
    Slug:        "nitro"
    project_url: "http://github.com/spf13/nitro"
    ---

    # Nitro

    Quick and easy performance analyzer library for Go.

    ## Overview

    Nitro is a quick and easy performance analyzer library for Go.
    It is useful for comparing A/B against different drafts of functions
    or different functions.

    ## Implementing Nitro

    Using Nitro is simple. First use go get to install the latest version
    of the library.

        $ go get github.com/spf13/nitro

    Next include nitro in your application.


