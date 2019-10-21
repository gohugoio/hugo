---
date: 2017-04-27T17:53:58-04:00
categories: ["Releases"]
description: ""
link: ""
title: "Hugo 0.20.6"
draft: false
author: bep
aliases: [/0-20-6/]
---

There have been some [shouting on discuss.gohugo.io](https://discuss.gohugo.io/t/index-md-is-generated-in-subfolder-index-index-html-hugo-0-20/6338/15) about some broken sites after the release of Hugo `0.20`. This release reintroduces the old behaviour, making  `/my-blog-post/index.md` work as expected.

Hugo now has:

* 16675&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 456&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 156&#43; [themes](http://themes.gohugo.io/)

## Fixes

* Avoid index.md in /index/index.html [#3396](https://github.com/gohugoio/hugo/issues/3396) 
* Make missing GitInfo a WARNING [b30ca4be](https://github.com/gohugoio/hugo/commit/b30ca4bec811dbc17e9fd05925544db2b75e0e49) [@bep](https://github.com/bep) [#3376](https://github.com/gohugoio/hugo/issues/3376) 
* Fix some of the fpm fields for deb [3bd1d057](https://github.com/gohugoio/hugo/commit/3bd1d0571d5f2f6bf0dc8f90a8adf2dbfcb2fdfd) [@anthonyfok](https://github.com/anthonyfok) 