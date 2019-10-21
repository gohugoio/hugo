
---
date: 2018-01-22
title: "Hugo 0.34: Pattern matching to filter images and other resources"
description: "Hugo 0.34 adds full glob with super-asterisk support, for example `*.jpg`."
categories: ["Releases"]
---

Hugo `0.34` is a small release. It contains a few smaller bug-fixes, but more important is an overhaul of the API used to find images and other resources in your page bundles.

We have added two simple methods on the `Resources` object:

* `.Match` finds every resource matching a pattern. Examples: `.Match "images/*.jpg"` finds every JPEG image in `images` and `.Match "**.jpg"` finds every JPEG image in the bundle.
* `.GetMatch` finds the first resource matching the pattern given.

**Note: The path separators used are Unix-style forward slashes, even on Windows.**

It uses [standard wildcard syntax](http://tldp.org/LDP/GNU-Linux-Tools-Summary/html/x11655.htm) with the addition of the `**`, aka super-asterisk, which matches across path boundaries.

Thanks to [@gobwas](https://github.com/gobwas/glob) for a fast and easy-to-use Glob library.

This release represents **5 contributions by 1 contributors** to the main Hugo code base.

Many have also been busy writing and fixing the documentation in [hugoDocs](https://github.com/gohugoio/hugoDocs),
which has received **25 contributions by 16 contributors**. A special thanks to [@bep](https://github.com/bep), [@rmetzler](https://github.com/rmetzler), [@chris-rudmin](https://github.com/chris-rudmin), and [@stkevintan](https://github.com/stkevintan) for their work on the documentation site.


Hugo now has:

* 22689+ [stars](https://github.com/gohugoio/hugo/stargazers)
* 448+ [contributors](https://github.com/gohugoio/hugo/graphs/contributors)
* 197+ [themes](http://themes.gohugo.io/)

## Notes
* `Resources.GetByPrefix` and  `Resources.ByPrefix` are deprecated. They still work, but will eventually be removed.  Use `Resources.Match` (many) and `Resources.GetMatch`  (one).
* When filtering bundles pages in sub-folders, you need to include the sub-folder when matching. This was a bug introduced in `0.33` and gets it in line with images and other resources.

## Enhancements

* Add `Resources.Match` and `Resources.GetMatch` [94213801](https://github.com/gohugoio/hugo/commit/9421380168f66620cb73203e1267814b3086d805) [@bep](https://github.com/bep) [#4301](https://github.com/gohugoio/hugo/issues/4301)

## Fixes
* Add validation for `defaultContentLanguage` [4d5e4f37](https://github.com/gohugoio/hugo/commit/4d5e4f379a890a3c6cbc11ddb40d77a90f14c015) [@bep](https://github.com/bep) [#4298](https://github.com/gohugoio/hugo/issues/4298)
* Fix lookup of pages bundled in sub-folders in `ByPrefix` etc. [5d030869](https://github.com/gohugoio/hugo/commit/5d03086981b4a7d4bc450269a6a2e0fd22dbeed7) [@bep](https://github.com/bep) [#4295](https://github.com/gohugoio/hugo/issues/4295)
