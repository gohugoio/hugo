---
date: 2017-06-24T17:53:58-04:00
categories: ["Releases"]
description: "This release fixes some important archetype-related regressions from Hugo 0.24"
link: ""
title: "Hugo 0.24.1"
draft: false
author: bep
aliases: [/0-24-1/]
---

This release fixes some important **archetype-related regressions** from the recent Hugo 0.24-release.

## Fixes

* Fix archetype regression when no archetype file [4294dd8d](https://github.com/gohugoio/hugo/commit/4294dd8d9d22bd8107b7904d5389967da1f83f27) [@bep](https://github.com/bep) [#3626](https://github.com/gohugoio/hugo/issues/3626) 
* Preserve shortcodes in archetype templates [b63e4ee1](https://github.com/gohugoio/hugo/commit/b63e4ee198c875b73a6a9af6bb809589785ed589) [@bep](https://github.com/bep) [#3623](https://github.com/gohugoio/hugo/issues/3623) 
* Fix handling of timezones with positive UTC offset (e.g., +0800) in TOML [0744f81e](https://github.com/gohugoio/hugo/commit/0744f81ec00bb8888f59d6c8b5f57096e07e70b1) [@bep](https://github.com/bep) [#3628](https://github.com/gohugoio/hugo/issues/3628) 

## Enhancements

* Create default archetype on new site [bfa336d9](https://github.com/gohugoio/hugo/commit/bfa336d96173377b9bbe2298dbd101f6a718c174) [@bep](https://github.com/bep) [#3626](https://github.com/gohugoio/hugo/issues/3626) 