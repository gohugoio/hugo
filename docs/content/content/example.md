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
prev: /content/summaries
title: Example Content File
weight: 70
---

Some things are better shown than explained. The following is a very basic example of a content file written in [Markdown](https://help.github.com/articles/github-flavored-markdown/):

**mysite/content/project/nitro.md → http://mysite.com/project/nitro.html**

With TOML front matter:

```markdown
+++
date        = "2013-06-21T11:27:27-04:00"
title       = "Nitro: A quick and simple profiler for Go"
description = "Nitro is a simple profiler for your Golang applications"
tags        = [ "Development", "Go", "profiling" ]
topics      = [ "Development", "Go" ]
slug        = "nitro"
project_url = "https://github.com/spf13/nitro"
+++

# Nitro

Quick and easy performance analyzer library for [Go](http://golang.org/).

## Overview

Nitro is a quick and easy performance analyzer library for Go.
It is useful for comparing A/B against different drafts of functions
or different functions.

## Implementing Nitro

Using Nitro is simple. First, use `go get` to install the latest version
of the library.

    $ go get github.com/spf13/nitro

Next, include nitro in your application.
```

You may also use the equivalent YAML front matter:

```markdown
---
date:        "2013-06-21T11:27:27-04:00"
title:       "Nitro: A quick and simple profiler for Go"
description: "Nitro is a simple profiler for your Go lang applications"
tags:        [ "Development", "Go", "profiling" ]
topics:      [ "Development", "Go" ]
slug:        "nitro"
project_url: "https://github.com/spf13/nitro"
---
```

`nitro.md` would be rendered as follows:

> # Nitro
>
> Quick and easy performance analyzer library for [Go](http://golang.org/).
>
> ## Overview
>
> Nitro is a quick and easy performance analyzer library for Go.
> It is useful for comparing A/B against different drafts of functions
> or different functions.
>
> ## Implementing Nitro
>
> Using Nitro is simple. First, use `go get` to install the latest version
> of the library.
>
>     $ go get github.com/spf13/nitro
>
> Next, include nitro in your application.

The source `nitro.md` file is converted to HTML by the excellent
[Blackfriday](https://github.com/russross/blackfriday) Markdown processor,
which supports extended features found in the popular
[GitHub Flavored Markdown](https://help.github.com/articles/github-flavored-markdown/).
