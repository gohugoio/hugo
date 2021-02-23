
---
date: 2019-08-17
title: "Hugo 0.57.2: A couple of Bug Fixes"
description: "This version fixes a couple of bugs introduced in 0.57.0."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

Hugo 0.57.0 had some well-intended breaking changes. And while they made a lot of sense, one of them made a little too much noise.

This release reverts the behavior for `.Pages` on the home page to how it behaved in 0.56, but adds a `WARNING` telling you what to do to prepare for Hugo 0.58.

In short, `.Page` home will from 0.58  only return its immediate children (sections and regular pages).

In this release it returns `.Site.RegularPages`. So to prepare for Hugo 0.58 you can either use `.Site.RegularPages` in your home template, or if you have a general `list.html` or RSS template, you can do something like this:

```go-html-template
{{- $pctx := . -}}
{{- if .IsHome -}}{{ $pctx = .Site }}{{- end -}}
{{- $pages := $pctx.RegularPages -}}
```

* tpl: Use RegularPages for RSS template [88d69936](https://github.com/gohugoio/hugo/commit/88d69936122f82fffc02850516bdb37be3d0892b) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)
* hugolib: Don't use the global warning logger [ea681603](https://github.com/gohugoio/hugo/commit/ea6816030081b2cffa6c0ae9ca5429a2c6fe2fa5) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)
* tpl: Avoid "home page warning" in RSS template [564cf1bb](https://github.com/gohugoio/hugo/commit/564cf1bb11e100891992e9131b271a79ea7fc528) [@bep](https://github.com/bep) [#6238](https://github.com/gohugoio/hugo/issues/6238)
* hugolib: Allow index.md inside bundles [4b4bdcfe](https://github.com/gohugoio/hugo/commit/4b4bdcfe740d988e4cfb4fee53eced6985576abd) [@bep](https://github.com/bep) [#6208](https://github.com/gohugoio/hugo/issues/6208)
* Adjust the default paginator for sections [18836a71](https://github.com/gohugoio/hugo/commit/18836a71ce7b671fa71dd1318b99fc661755e94d) [@bep](https://github.com/bep) [#6231](https://github.com/gohugoio/hugo/issues/6231)
* hugolib: Add a site benchmark [416493b5](https://github.com/gohugoio/hugo/commit/416493b548a9bbaa27758fba9bab50a22b680e9d) [@bep](https://github.com/bep) 
* Update to Go 1.11.13 and 1.12.9 [f28efd35](https://github.com/gohugoio/hugo/commit/f28efd35820dc4909832c14dfd8ea6812ecead31) [@bep](https://github.com/bep) [#6228](https://github.com/gohugoio/hugo/issues/6228)



