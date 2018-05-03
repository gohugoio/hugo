---
date: 2017-04-24T13:53:58-04:00
categories: ["Releases"]
description: "This is a bug-fix release with one important fix. But it also adds some harness around GoReleaser"
link: ""
title: "Hugo 0.20.3"
draft: false
author: bep
aliases: [/0-20-3/]
---

This is a bug-fix release with one important fix. But it also adds some harness around [GoReleaser](https://github.com/goreleaser/goreleaser) to automate the Hugo release process. Big thanks to [@caarlos0](https://github.com/caarlos0) for great and super-fast support fixing issues along the way.

Hugo now has:

* 16619&#43; [stars](https://github.com/gohugoio/hugo/stargazers)
* 458&#43; [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 156&#43; [themes](http://themes.gohugo.io/)

## Enhancement

* Automate the Hugo release process [550eba64](https://github.com/gohugoio/hugo/commit/550eba64705725eb54fdb1042e0fb4dbf6f29fd0) [@bep](https://github.com/bep) [#3358](https://github.com/gohugoio/hugo/issues/3358) 

## Fix

* Fix handling of zero-length files [9bf5c381](https://github.com/gohugoio/hugo/commit/9bf5c381b6b3e69d4d8dbfd7a40074ac44792bbf) [@bep](https://github.com/bep) [#3355](https://github.com/gohugoio/hugo/issues/3355) 