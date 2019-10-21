---
date: 2017-06-13T17:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.22.1 fixes a couple of issues reported after the 0.22 release"
link: ""
title: "Hugo 0.22.1"
draft: false
author: bep
aliases: [/0-22-1/]
---

Hugo `0.22.1` fixes a couple of issues reported after the [0.22 release](https://github.com/gohugoio/hugo/releases/tag/v0.22) Monday. Most importantly a fix for detecting regular subfolders below the root-sections.

Also, we forgot to adapt the `permalink settings` with support for nested sections, which made that feature less useful than it could be.

With this release you can configure **permalinks with sections** like this:

**First level only:**

```
[permalinks]
blog = ":section/:title"
```

**Nested (all levels):**

```
[permalinks]
blog = ":sections/:title"
```
## Fixes

* Fix section logic for root folders with subfolders [a30023f5](https://github.com/gohugoio/hugo/commit/a30023f5cbafd06034807255181a5b7b17f3c25f) [@bep](https://github.com/bep) [#3586](https://github.com/gohugoio/hugo/issues/3586) 
* Support sub-sections in permalink settings [1f26420d](https://github.com/gohugoio/hugo/commit/1f26420d392a5ab4c7b7fe1911c0268b45d01ab8) [@bep](https://github.com/bep) [#3580](https://github.com/gohugoio/hugo/issues/3580) 
* Adjust rlimit to 64000 [ff54b6bd](https://github.com/gohugoio/hugo/commit/ff54b6bddcefab45339d8dc2b13776b92bdc04b9) [@bep](https://github.com/bep) [#3582](https://github.com/gohugoio/hugo/issues/3582) 
* Make error on setting rlimit a warning only [629e1439](https://github.com/gohugoio/hugo/commit/629e1439e819a7118ae483381d4634f16d3474dd) [@bep](https://github.com/bep) [#3582](https://github.com/gohugoio/hugo/issues/3582) 
* Revert: Remove the rlimit tweaking on macOS [26aa06a3](https://github.com/gohugoio/hugo/commit/26aa06a3db57ab7134a900d641fa2976f7971520) [@bep](https://github.com/bep) [#3582](https://github.com/gohugoio/hugo/issues/3582)