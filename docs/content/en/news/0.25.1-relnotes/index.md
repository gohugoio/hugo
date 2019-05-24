---
date: 2017-07-10T17:53:58-04:00
categories: ["Releases"]
description: "This is a bug-fix release with a couple of important fixes"
link: ""
title: "Hugo 0.25.1"
draft: false
author: bep
aliases: [/0-25-1/]
---

This is a bug-fix release with a couple of important fixes.

Hugo now has:

* 18277+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 456+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 170+ [themes](http://themes.gohugo.io/)

## Fixes

* Fix union when the first slice is empty [dbbc5c48](https://github.com/gohugoio/hugo/commit/dbbc5c4810a04ac06fad7500d88cf5c3bfe0c7fd) [@bep](https://github.com/bep) [#3686](https://github.com/gohugoio/hugo/issues/3686)
* Navigate to changed on CREATE When working with content from IntelliJ IDE, like WebStorm, every file save is followed by two events: "RENAME" and then "CREATE". [7bcc1ce6](https://github.com/gohugoio/hugo/commit/7bcc1ce659710f2220b400ce3b76e50d2e48b241) [@miltador](https://github.com/miltador) 
* Final (!) fix for issue with escaped JSON front matter [7f82b41a](https://github.com/gohugoio/hugo/commit/7f82b41a24af0fd04d28fbfebf9254766a3c6e6f) [@bep](https://github.com/bep) [#3682](https://github.com/gohugoio/hugo/issues/3682)
* Fix issue with escaped JSON front matter [84db6c74](https://github.com/gohugoio/hugo/commit/84db6c74a084d2b52117b999d4ec343cd3389a68) [@bep](https://github.com/bep) [#3682](https://github.com/gohugoio/hugo/issues/3682)