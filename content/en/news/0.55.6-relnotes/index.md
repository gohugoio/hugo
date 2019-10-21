
---
date: 2019-05-18
title: "Hugo 0.55.6: One Bug Fix!"
description: "Fixes some reported paginator crashes in server mode."
categories: ["Releases"]
images:
- images/blog/hugo-bug-poster.png

---

This is a bug-fix release with one important fix. There have been reports about infrequent paginator crashes when running the Hugo server since 0.55.0. The reason have been narrowed down to that of parallel rebuilds. This isn't a new thing, but the changes in 0.55.0 made it extra important to serialize the page initialization. This release fixes that by protecting the `Build` method with a lock when running in server mode. [95ce2a40](https://github.com/gohugoio/hugo/commit/95ce2a40e734bb82b69f9a64270faf3ed69c92cc) [@bep](https://github.com/bep) [#5885](https://github.com/gohugoio/hugo/issues/5885)[#5968](https://github.com/gohugoio/hugo/issues/5968)

